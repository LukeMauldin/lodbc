package lodbc

import (
	"database/sql/driver"
	"fmt"
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
	"unsafe"
)

type connection struct {
	//Connection handle
	handle syscall.Handle

	//Is transaction active
	isTransactionActive bool
	
	//Statements owned by the connection
	statements map[driver.Stmt]bool

	//Is closed -- allows Close() to be called multiple times without error
	isClosed bool
}

func (c *connection) Prepare(query string) (driver.Stmt, error) {
	//Allocate the statement handle
	var stmtHandle syscall.Handle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_STMT, c.handle, &stmtHandle)
	if IsError(ret) {
		return nil, ErrorConnection(c.handle)
	}

	//Set the query timeout
	ret = odbc.SQLSetStmtAttr(stmtHandle, odbc.SQL_ATTR_QUERY_TIMEOUT, int32(QUERY_TIMEOUT), odbc.SQL_IS_INTEGER)
	if IsError(ret) {
		return nil, ErrorStatement(stmtHandle, query)
	}

	//Get the statement descriptor table
	var stmtDescHandle syscall.Handle
	ret = odbc.SQLGetStmtAttr(stmtHandle, odbc.SQL_ATTR_APP_PARAM_DESC, uintptr(unsafe.Pointer(&stmtDescHandle)), 0, nil)
	if IsError(ret) {
		return nil, ErrorConnection(c.handle)
	}

	//Create new statement
	stmt := &statement{handle: stmtHandle, stmtDescHandle: stmtDescHandle, sqlStmt: query, conn: c}
	
	//Add to map of statements owned by the connection
	c.statements[stmt] = true

	return stmt, nil
}

func (c *connection) Close() error {
	//Verify that connHandle is valid
	if c.handle == 0 {
		return nil
	}

	//Verify that connection has not already been closed
	if c.isClosed {
		return nil
	}

	var err error
	isError := false
	
	//Close all of the statements owned by the connection
	for key, _ := range c.statements {
		//Skip the statement if it is already nil
		if isNil(key) {
			continue
		}
		key.Close()
	}
	c.statements = nil

	//If the transaction is active, roll it back
	if c.isTransactionActive {
		ret := odbc.SQLEndTran(odbc.SQL_HANDLE_DBC, c.handle, odbc.SQL_ROLLBACK)
		if IsError(ret) {
			err = ErrorConnection(c.handle)
			isError = true
		}
	}

	//Disconnect connection
	ret := odbc.SQLDisconnect(c.handle)
	if IsError(ret) {
		err = ErrorConnection(c.handle)
		isError = true
	}

	//Deallocate connection 
	ret = odbc.SQLFreeHandle(odbc.SQL_HANDLE_DBC, c.handle)
	if IsError(ret) {
		err = ErrorConnection(c.handle)
		isError = true
	}

	//Return any error
	if isError {
		return err
	}

	//Set connection to closed
	c.isClosed = true

	return nil
}

func (c *connection) Begin() (driver.Tx, error) {
	//Do not allow a  new transaction if one already exists
	if c.isTransactionActive {
		return nil, fmt.Errorf("Transaction already active for connection")
	}

	ret := odbc.SQLSetConnectAttr(c.handle, odbc.SQL_ATTR_AUTOCOMMIT, odbc.SQL_AUTOCOMMIT_OFF, 0, nil)
	if IsError(ret) {
		return nil, ErrorConnection(c.handle)
	}
	c.isTransactionActive = true

	tx := &transaction{conn: c}
	return tx, nil
}

func (c *connection) IsTransactionActive() bool {
	return c.isTransactionActive
}


func (c *connection) closeStatement(stmt driver.Stmt) {
	delete(c.statements, stmt)
}