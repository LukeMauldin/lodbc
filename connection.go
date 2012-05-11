package lodbc

import (
	"github.com/LukeMauldin/lodbc"
	"syscall"
	"unsafe"
)

type Connection struct {
	//Connection handle
	handle syscall.Handle

	//Statements created by this connection
	statements []IStatement
	
	//Is transaction active
	isTransactionActive bool
	
	//Is closed -- allows Close() to be called multiple times without error
	isClosed bool
}

func (c *Connection) Close() (error) {
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
	
	//If the transaction is active, roll it back
	ret := odbc.SQLEndTran(odbc.SQL_HANDLE_DBC, c.handle, odbc.SQL_ROLLBACK)
	if IsError(ret) {
		err = ErrorConnection(c.handle)
		isError = true
	}
		
	//Close all of the owned statements
	for _, statement := range c.statements {
		err = statement.Close()
		isError = true
	}
	
	//Disconnect connection
	ret = odbc.SQLDisconnect(c.handle)
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


func (c *Connection) Begin() (error) {
	ret := odbc.SQLSetConnectAttr(c.handle, odbc.SQL_ATTR_AUTOCOMMIT, odbc.SQL_AUTOCOMMIT_OFF, 0, nil)
	if IsError(ret) {
		return ErrorConnection(c.handle)
	}
	c.isTransactionActive = true
	return nil
}

func (c *Connection) Commit() (error) {
	return c.completeTransaction(odbc.SQL_COMMIT)
}

func (c *Connection) Rollback() (error) {
	return c.completeTransaction(odbc.SQL_ROLLBACK)
}

func (c *Connection) completeTransaction(completeType odbc.SQLTransactionOption) (error) {
	//Complete transaction by either committing or rolling back
	ret := odbc.SQLEndTran(odbc.SQL_HANDLE_DBC, c.handle, completeType)
	if IsError(ret) {
		return ErrorConnection(c.handle)
	}
	
	//Make transaction as finished and turn auto commit back on
	c.isTransactionActive = false
	ret = odbc.SQLSetConnectAttr(c.handle, odbc.SQL_ATTR_AUTOCOMMIT, odbc.SQL_AUTOCOMMIT_ON, 0, nil)
	if IsError(ret) {
		return ErrorConnection(c.handle)
	}
	return nil
}

func (c *Connection) NewStatement() (IStatement, error) {
	//Allocate the statement handle
	var stmtHandle syscall.Handle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_STMT, c.handle, &stmtHandle)
	if IsError(ret) {
		return nil, ErrorConnection(c.handle)
	}
	
	//Set the query timeout
	ret = odbc.SQLSetStmtAttr(stmtHandle, odbc.SQL_ATTR_QUERY_TIMEOUT, int32(QUERY_TIMEOUT), odbc.SQL_IS_INTEGER)
	if IsError(ret) {
		return nil, ErrorStatement(stmtHandle, "")
	}
	
	//Get the statement descriptor table
	var stmtDescHandle syscall.Handle
	ret = odbc.SQLGetStmtAttr(stmtHandle, odbc.SQL_ATTR_APP_PARAM_DESC, uintptr(unsafe.Pointer(&stmtDescHandle)), 0, nil) 
	if IsError(ret) {
		return nil, ErrorConnection(c.handle)
	}
	
	//Create new statement
	stmt := &Statement {handle: stmtHandle, stmtDescHandle: stmtDescHandle}
	
	return stmt, nil
}

func (c *Connection) IsTransactionActive() (bool) {
	return c.isTransactionActive
}