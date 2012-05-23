package lodbc

import (
	"database/sql"
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
	"time"
)

//Global variables
var (
	queryTimeout = 240 * time.Second // Query timeout
)

// Shared global environment
var envHandle syscall.Handle

// Allocates shared environment
func init() {

	// Set environment handle for connection pooling
	ret := odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_CONNECTION_POOLING, odbc.SQL_CP_ONE_PER_DRIVER, 0)
	if isError(ret) {
		panic(errorEnvironment(envHandle))
	}

	// Allocate the environment handle	
	ret = odbc.SQLAllocHandle(odbc.SQL_HANDLE_ENV, 0, &envHandle)
	if isError(ret) {
		panic(errorEnvironment(envHandle))
	}

	// Set the environment handle to use ODBC v3
	ret = odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_ODBC_VERSION, odbc.SQL_OV_ODBC3, 0)
	if isError(ret) {
		panic(errorEnvironment(envHandle))
	}

	// Register with the SQL package
	d := &lodbcDriver{}
	sql.Register("lodbc", d)
}

// Frees environment handle -- calling this will make the lodbc package unusable because all setup is performed in init()
func FreeEnvironment() error {
	ret := odbc.SQLFreeHandle(odbc.SQL_HANDLE_ENV, envHandle)
	if isError(ret) {
		return errorEnvironment(envHandle)
	}
	return nil
}

// Enumeration for supported ODBC version
type ODBCVersion int
const (
	ODBCVersion_3   ODBCVersion = 1
	ODBCVersion_380 ODBCVersion = 2
)

// Sets the ODBC version for the environment
func SetODBCVersion(version ODBCVersion) {
	switch version {
	case ODBCVersion_3:
		ret := odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_ODBC_VERSION, odbc.SQL_OV_ODBC3, 0)
		if isError(ret) {
			panic(errorEnvironment(envHandle))
		}
		break
	case ODBCVersion_380:
		ret := odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_ODBC_VERSION, odbc.SQL_OV_ODBC3_80, 0)
		if isError(ret) {
			panic(errorEnvironment(envHandle))
		}
		break
	}
}

//Sets the global query timeout
func SetQueryTimeout(timeout time.Duration) {
	queryTimeout = timeout
}
