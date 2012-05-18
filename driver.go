package lodbc

import (
	"database/sql/driver"
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
)

type lodbcDriver struct {
}

func (d *lodbcDriver) Open(name string) (driver.Conn, error) {
	//Allocate the connection handle
	var connHandle syscall.Handle
	ret := odbc.SQLAllocHandle(odbc.SQL_HANDLE_DBC, envHandle, &connHandle)
	if IsError(ret) {
		return nil, ErrorEnvironment(envHandle)
	}

	//Establish the connection with the database
	ret = odbc.SQLDriverConnect(connHandle, 0, syscall.StringToUTF16Ptr(name), odbc.SQL_NTS, nil, 0, nil, odbc.SQL_DRIVER_NOPROMPT)
	if IsError(ret) {
		return nil, ErrorConnection(connHandle)
	}

	//Create new connection
	var conn = &connection{handle: connHandle, isTransactionActive: false}

	return conn, nil
}
