package lodbc

import (
	"database/sql/driver"
	"github.com/LukeMauldin/lodbc/odbc"
	"runtime"
	"syscall"
	"unsafe"
)

// Implements type database/sql/driver Driver interface
type lodbcDriver struct {
}

// Returns a new connection to the database
func (d *lodbcDriver) Open(name string) (driver.Conn, error) {
	// Allocate the connection handle
	var connHandle odbc.SQLHandle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_DBC, envHandle, &connHandle)
	if isError(ret) {
		return nil, errorEnvironment(envHandle)
	}

	// Establish the connection with the database
	nameSqlPtr := (*odbc.SQLCHAR)(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	ret = odbc.SQLDriverConnect(connHandle, 0, nameSqlPtr, odbc.SQLSMALLINT(odbc.SQL_NTS), nil, 0, nil, odbc.SQL_DRIVER_NOPROMPT)
	if isError(ret) {
		return nil, errorConnection(connHandle)
	}

	// Create new connection
	var conn = &connection{handle: connHandle, isTransactionActive: false, statements: make(map[driver.Stmt]bool, 0)}

	//Add a finalizer
	runtime.SetFinalizer(conn, (*connection).Close)

	return conn, nil
}
