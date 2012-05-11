package lodbc

import (
	"github.com/LukeMauldin/go-lodbc/odbc"
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
}

func NewConnection(connStr string) (IConnection, error) {
	//Allocate the connection handle
	var connHandle syscall.Handle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_DBC, envHandle, &connHandle)
	if IsError(ret) {
		return nil, ErrorEnvironment(envHandle)
	}

	//Establish the connection with the database
	ret = odbc.SQLDriverConnect(connHandle, 0, syscall.StringToUTF16Ptr(connStr), odbc.SQL_NTS, nil, 0, nil, odbc.SQL_DRIVER_NOPROMPT)
	if IsError(ret) {
		return nil, ErrorConnection(connHandle)
	}
	
	//Create new connection
	var conn = &Connection {handle: connHandle, statements: make([]IStatement, 0), isTransactionActive: false}
	
	return conn, nil
}

func FreeEnvironment() (error) {
	ret := odbc.SQLFreeHandle(odbc.SQL_HANDLE_ENV, envHandle)
	if IsError(ret) {
		return ErrorEnvironment(envHandle)
	}
	return nil
}