package lodbc

import (
	"database/sql/driver"
	"fmt"
	"github.com/LukeMauldin/lodbc/odbc"
	"io"
	"syscall"
	"time"
	"unsafe"
)

// Implements type database/sql/driver Rows interface
type rows struct {
	// Statement handle
	handle syscall.Handle

	// Descriptor handle
	descHandle syscall.Handle

	// Bool indicating if any rows have been read
	isBeforeFirst bool

	// Store last error
	lastError error

	// Is closed -- allows Close() to be called multiple times without error
	isClosed bool

	// SQL statement used to generate rows -- used in error reporting
	sqlStmt string

	// Result column defintions
	ResultColumnDefs []ResultColumnDef
	
	// Result column names
	resultColumnNames []string
}

// Returns the names of the columns
func (rows *rows) Columns() []string {
	return rows.resultColumnNames	
}

// Next is called to populate the next row of data into the provided slice
func (rows *rows) Next(dest []driver.Value) error {
	//If this is the first time rows has been read, setup necessary field level information
	if rows.isBeforeFirst {
		for index, resultColumnDef := range rows.ResultColumnDefs {
			//Set precision and scale for numeric fields
			if resultColumnDef.DataType == odbc.SQL_NUMERIC || resultColumnDef.DataType == odbc.SQL_DECIMAL {
				colIndex := odbc.SQLSMALLINT(index + 1)
				odbc.SQLSetDescField(rows.descHandle, colIndex, odbc.SQL_DESC_TYPE, odbc.SQL_C_NUMERIC, 0)
				odbc.SQLSetDescField(rows.descHandle, colIndex, odbc.SQL_DESC_PRECISION, resultColumnDef.Precision, 0)
				odbc.SQLSetDescField(rows.descHandle, colIndex, odbc.SQL_DESC_SCALE, resultColumnDef.Scale, 0)
			}
		}

		//Update isBeforeFirst
		rows.isBeforeFirst = false
	}

	//Fetch a row of data
	ret := odbc.SQLFetch(rows.handle)
	if ret == odbc.SQL_NO_DATA {
		//No more data to read
		return io.EOF
	} else if IsError(ret) {
		return errorStatement(rows.handle, rows.sqlStmt)
	}

	//Get a row of data
	err := rows.getRow(dest)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the rows iterator
func (rows *rows) Close() error {
	//Verify that rows has not already been closed
	if rows.isClosed {
		return nil
	}

	//Close the cursor
	ret := odbc.SQLCloseCursor(rows.handle)
	if IsError(ret) {
		return errorStatement(rows.handle, rows.sqlStmt)
	}

	//Mark the rows as closed
	rows.isClosed = true

	return nil
}

//Get a single row of data by calling getField for each column
func (rows *rows) getRow(dest []driver.Value) error {
	for index, _ := range rows.ResultColumnDefs {
		fieldValue, ret := rows.getField(index + 1)
		if IsError(ret) {
			return errorStatement(rows.handle, rows.sqlStmt)
		}
		dest[index] = fieldValue
	}
	return nil
}

// Return a single column of data
func (rows *rows) getField(index int) (v interface{}, ret odbc.SQLReturn) {
	columnDef := rows.ResultColumnDefs[index-1]
	var fieldInd odbc.SQLValueIndicator
	switch columnDef.DataType {
	case odbc.SQL_BIT:
		var value bool
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_BIT, valuePtr, 0, &fieldInd)
		return formatGetFieldReturn(value, fieldInd, ret)
	case odbc.SQL_INTEGER, odbc.SQL_SMALLINT, odbc.SQL_TINYINT:
		var value int
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_LONG, valuePtr, 0, &fieldInd)
		return formatGetFieldReturn(value, fieldInd, ret)
	case odbc.SQL_BIGINT:
		var value int64
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_LONG, valuePtr, 0, &fieldInd)
		return formatGetFieldReturn(value, fieldInd, ret)
	case odbc.SQL_FLOAT:
		var value float64
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_FLOAT, valuePtr, 0, &fieldInd)
		return formatGetFieldReturn(value, fieldInd, ret)
	case odbc.SQL_DOUBLE, odbc.SQL_REAL:
		var value float64
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_DOUBLE, valuePtr, 0, &fieldInd)
		return formatGetFieldReturn(value, fieldInd, ret)
	case odbc.SQL_NUMERIC, odbc.SQL_DECIMAL:
		var value odbc.SQL_NUMERIC_STRUCT
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_ARD_TYPE, valuePtr, 0, &fieldInd)
		return formatGetFieldReturn(numericToFloat(value), fieldInd, ret)
	case odbc.SQL_CHAR, odbc.SQL_VARCHAR, odbc.SQL_LONGVARCHAR, odbc.SQL_WCHAR, odbc.SQL_WVARCHAR:
		value := make([]uint16, columnDef.Length*2+2)
		valuePtr := uintptr(unsafe.Pointer(&value[0]))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_WCHAR, valuePtr, odbc.SQLLEN(columnDef.Length*2+2), &fieldInd)
		return formatGetFieldReturn(syscall.UTF16ToString(value), fieldInd, ret)
	case odbc.SQL_TYPE_DATE:
		var value odbc.SQL_DATE_STRUCT
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_DATE, valuePtr, 0, &fieldInd)
		time := time.Date(int(value.Year), time.Month(value.Month), int(value.Day), 0, 0, 0, 0, time.UTC)
		return formatGetFieldReturn(time, fieldInd, ret)
	case odbc.SQL_TYPE_TIMESTAMP:
		var value odbc.SQL_TIMESTAMP_STRUCT
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_TIMESTAMP, valuePtr, 0, &fieldInd)
		time := time.Date(int(value.Year), time.Month(value.Month), int(value.Day), int(value.Hour), int(value.Minute), int(value.Second), int(value.Faction), time.UTC)
		return formatGetFieldReturn(time, fieldInd, ret)
	default:
		panic(fmt.Sprintf("ODBC type not supported: {%v}. Column name: %v", columnDef.DataType, columnDef.Name))
	}

	return nil, odbc.SQL_SUCCESS

}

// Utility function to format return value
func formatGetFieldReturn(value interface{}, fieldInd odbc.SQLValueIndicator, getDataRet odbc.SQLReturn) (interface{}, odbc.SQLReturn) {
	if IsError(getDataRet) {
		return nil, getDataRet
	} else if fieldInd == odbc.SQL_NULL_DATA {
		return nil, odbc.SQL_SUCCESS
	} else {
		return value, odbc.SQL_SUCCESS
	}
	return nil, odbc.SQL_SUCCESS
}
