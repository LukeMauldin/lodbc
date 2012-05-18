package lodbc

import (
	"database/sql"
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
)

//Global variables
var (
	QUERY_TIMEOUT = 240 //240 seconds
)

//Shared global environment
var envHandle syscall.Handle

//Allocates shared environment
func init() {

	//Set environment handle for connection pooling
	ret := odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_CONNECTION_POOLING, odbc.SQL_CP_ONE_PER_DRIVER, 0)
	if IsError(ret) {
		panic(ErrorEnvironment(envHandle))
	}

	// Allocate the environment handle	
	ret = odbc.SQLAllocHandle(odbc.SQL_HANDLE_ENV, 0, &envHandle)
	if IsError(ret) {
		panic(ErrorEnvironment(envHandle))
	}

	// Set the environment handle to use ODBC v3
	ret = odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_ODBC_VERSION, odbc.SQL_OV_ODBC3, 0)
	if IsError(ret) {
		panic(ErrorEnvironment(envHandle))
	}

	//Register with the SQL package
	d := &lodbcDriver{}
	sql.Register("lodbc", d)
}

func FreeEnvironment() error {
	ret := odbc.SQLFreeHandle(odbc.SQL_HANDLE_ENV, envHandle)
	if IsError(ret) {
		return ErrorEnvironment(envHandle)
	}
	return nil
}

type ODBCVersion int

const (
	ODBCVersion_3   ODBCVersion = 1
	ODBCVersion_380 ODBCVersion = 2
)

func SetODBCVersion(version ODBCVersion) {
	switch version {
	case ODBCVersion_3:
		ret := odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_ODBC_VERSION, odbc.SQL_OV_ODBC3, 0)
		if IsError(ret) {
			panic(ErrorEnvironment(envHandle))
		}
		break
	case ODBCVersion_380:
		ret := odbc.SQLSetEnvAttr(envHandle, odbc.SQL_ATTR_ODBC_VERSION, odbc.SQL_OV_ODBC3_80, 0)
		if IsError(ret) {
			panic(ErrorEnvironment(envHandle))
		}
		break
	}
}
