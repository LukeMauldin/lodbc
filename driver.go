package lodbc

import (
	"database/sql/driver"
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
	"runtime"
)


// Implements type database/sql/driver Driver interface
type lodbcDriver struct {
}

// Returns a new connection to the databas
func (d *lodbcDriver) Open(name string) (driver.Conn, error) {
	// Allocate the connection handle
	var connHandle syscall.Handle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_DBC, envHandle, &connHandle)
	if isError(ret) {
		return nil, errorEnvironment(envHandle)
	}

	// Establish the connection with the database
	ret = odbc.SQLDriverConnect(connHandle, 0, syscall.StringToUTF16Ptr(name), odbc.SQL_NTS, nil, 0, nil, odbc.SQL_DRIVER_NOPROMPT)
	if isError(ret) {
		return nil, errorConnection(connHandle)
	}

	// Create new connection
	var conn = &connection{handle: connHandle, isTransactionActive: false, statements: make(map[driver.Stmt]bool, 0)}
	
	//Add a finalizer
	runtime.SetFinalizer(conn, (*connection).Close)

	return conn, nil
}
