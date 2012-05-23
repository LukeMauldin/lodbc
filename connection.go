package lodbc

import (
	"database/sql/driver"
	"fmt"
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
	"unsafe"
	"runtime"
)

// Implements type database/sql/driver Conn interface
type connection struct {
	
	// Connection handle
	handle syscall.Handle

	// Is transaction active
	isTransactionActive bool
	
	// Statements owned by the connection
	statements map[driver.Stmt]bool

	// Is closed -- allows Close() to be called multiple times without error
	isClosed bool
}

// Prepare returns a prepared statement, bound to this connection
func (c *connection) Prepare(query string) (driver.Stmt, error) {
	
	// Allocate the statement handle
	var stmtHandle syscall.Handle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_STMT, c.handle, &stmtHandle)
	if isError(ret) {
		return nil, errorConnection(c.handle)
	}

	// Set the query timeout
	ret = odbc.SQLSetStmtAttr(stmtHandle, odbc.SQL_ATTR_QUERY_TIMEOUT, int32(queryTimeout.Seconds()), odbc.SQL_IS_INTEGER)
	if isError(ret) {
		return nil, errorStatement(stmtHandle, query)
	}

	// Get the statement descriptor table
	var stmtDescHandle syscall.Handle
	ret = odbc.SQLGetStmtAttr(stmtHandle, odbc.SQL_ATTR_APP_PARAM_DESC, uintptr(unsafe.Pointer(&stmtDescHandle)), 0, nil)
	if isError(ret) {
		return nil, errorConnection(c.handle)
	}

	// Create new statement
	stmt := &statement{handle: stmtHandle, stmtDescHandle: stmtDescHandle, sqlStmt: query, conn: c}
	
	// Add to map of statements owned by the connection
	c.statements[stmt] = true
	
	//Add a finalizer
	runtime.SetFinalizer(stmt, (*statement).Close)

	return stmt, nil
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
func (c *connection) Close() error {
	
	// Verify that connHandle is valid
	if c.handle == 0 {
		return nil
	}

	// Verify that connection has not already been closed
	if c.isClosed {
		return nil
	}

	var err error
	
	// Close all of the statements owned by the connection
	for key, _ := range c.statements {
		// Skip the statement if it is already nil
		if isNil(key) {
			continue
		}
		key.Close()
	}
	c.statements = nil

	// If the transaction is active, roll it back
	if c.isTransactionActive {
		ret := odbc.SQLEndTran(odbc.SQL_HANDLE_DBC, c.handle, odbc.SQL_ROLLBACK)
		if isError(ret) {
			err = errorConnection(c.handle)
		}
		
		//Turn AutoCommit back on
		ret = odbc.SQLSetConnectAttr(c.handle, odbc.SQL_ATTR_AUTOCOMMIT, odbc.SQL_AUTOCOMMIT_ON, 0, nil)
		if isError(ret) {
			err = errorConnection(c.handle)
		}
	}

	// Disconnect connection
	ret := odbc.SQLDisconnect(c.handle)
	if isError(ret) {
		err = errorConnection(c.handle)
	}

	// Deallocate connection 
	ret = odbc.SQLFreeHandle(odbc.SQL_HANDLE_DBC, c.handle)
	if isError(ret) {
		err = errorConnection(c.handle)
	}
	
	// Clear the handle
	c.handle = 0
	
	// Set connection to closed
	c.isClosed = true
	
	//Clear the finalizer
	runtime.SetFinalizer(c, nil)

	// Return any error
	if err != nil {
		return err
	}

	return nil
}

// Begin starts and returns a new transaction
// Only one transaction is supported at a time for a connection
func (c *connection) Begin() (driver.Tx, error) {
	// Do not allow a  new transaction if one already exists
	if c.isTransactionActive {
		return nil, fmt.Errorf("Transaction already active for connection")
	}

	ret := odbc.SQLSetConnectAttr(c.handle, odbc.SQL_ATTR_AUTOCOMMIT, odbc.SQL_AUTOCOMMIT_OFF, 0, nil)
	if isError(ret) {
		return nil, errorConnection(c.handle)
	}
	c.isTransactionActive = true

	tx := &transaction{conn: c}
	return tx, nil
}

// To be called by the statements owned by the connection when the statement is closed
// Removed the statement from the connection's list of statements'
func (c *connection) closeStatement(stmt driver.Stmt) {
	delete(c.statements, stmt)
}