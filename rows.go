package lodbc

import (
	"fmt"
	"database/lodbc/odbc"
	"syscall"
	"time"
	"unsafe"
	"reflect"
	"errors"
	"strconv"
)

type Rows struct {
	//Statement handle
	handle syscall.Handle

	//Descriptor handle
	descHandle syscall.Handle

	//Bool indicating if any rows have been read
	isBeforeFirst bool
	
	//Store last error
	lastError error
	
	//Is closed -- allows Close() to be called multiple times without error
	isClosed bool
	
	//SQL statement used to generate rows -- used in error reporting
	sqlStmt string

	//Result column defintions
	ResultColumnDefs []ResultColumnDef
}

// RawBytes is a byte slice that holds a reference to memory owned by
// the database itself. After a Scan into a RawBytes, the slice is only
// valid until the next call to Next, Scan, or Close.
type RawBytes []byte


func (rows *Rows) Next() (bool) {
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
		return false
	} else if IsError(ret) {
		rows.lastError = ErrorStatement(rows.handle, rows.sqlStmt) 
		return false
	} else {
		return true
	}

	return false
}

func (rows *Rows) Close() error {
	//Verify that rows has not already been closed
	if rows.isClosed {
		return nil
	}
	
	//Close the cursor
	ret := odbc.SQLCloseCursor(rows.handle)
	if IsError(ret) {
		return ErrorStatement(rows.handle, rows.sqlStmt)
	}
	
	//Mark the rows as closed
	rows.isClosed = true
	
	return nil
}

func (rows *Rows) Scan(dest ...interface{}) error {
	if rows.isBeforeFirst {
		return fmt.Errorf("sql: Scan called without calling Next")
	}
	if len(dest) != len(rows.ResultColumnDefs) {
		return fmt.Errorf("sql: expected %d destination arguments in Scan, not %d", len(rows.ResultColumnDefs), len(dest))
	}
	
	dbValues, err := rows.getRow()
	if err != nil {
		return err
	}
	for i, sv := range dbValues {
		err := convertAssign(dest[i], sv)
		if err != nil {
			return fmt.Errorf("sql: Scan error on column index %d: %v", i, err)
		}
	}
	
	for _, dp := range dest {
		b, ok := dp.(*[]byte)
		if !ok {
			continue
		}
		if *b == nil {
			// If the []byte is now nil (for a NULL value),
			// don't fall through to below which would
			// turn it into a non-nil 0-length byte slice
			continue
		}
		if _, ok = dp.(*RawBytes); ok {
			continue
		}
		clone := make([]byte, len(*b))
		copy(clone, *b)
		*b = clone
	}
	return nil
}

func (rows *Rows) Err() (error) {
	return rows.lastError
}

func (rows *Rows) getRow() ([]interface{}, error) {
	dest := make([]interface{}, len(rows.ResultColumnDefs))
	for index, _ := range rows.ResultColumnDefs {
		fieldValue, ret := rows.getField(index + 1)
		if IsError(ret) {
			return nil, ErrorStatement(rows.handle, rows.sqlStmt)
		}
		dest[index] = fieldValue
	}
	return dest, nil
}

func (rows *Rows) getField(index int) (v interface{}, ret odbc.SQLReturn) {
	columnDef := rows.ResultColumnDefs[index - 1]
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
		value := make([]uint16, columnDef.Length * 2 + 2)
		valuePtr := uintptr(unsafe.Pointer(&value[0]))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_WCHAR, valuePtr, odbc.SQLLEN(columnDef.Length * 2 + 2), &fieldInd)	
		return formatGetFieldReturn(syscall.UTF16ToString(value), fieldInd, ret)
	case odbc.SQL_DATE:
		var value odbc.SQL_DATE_STRUCT
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_DATE, valuePtr, 0, &fieldInd)	
		time := time.Date(int(value.Year), time.Month(value.Month), int(value.Day), 0, 0, 0, 0, time.UTC)
		return formatGetFieldReturn(time, fieldInd, ret)
	case odbc.SQL_DATETIME:
		var value odbc.SQL_TIMESTAMP_STRUCT
		valuePtr := uintptr(unsafe.Pointer(&value))
		ret = odbc.SQLGetData(rows.handle, uint16(index), odbc.SQL_C_TIMESTAMP, valuePtr, 0, &fieldInd)
		time := time.Date(int(value.Year), time.Month(value.Month), int(value.Day), int(value.Hour), int(value.Minute), int(value.Second), int(value.Faction), time.UTC)
		return formatGetFieldReturn(time, fieldInd, ret)
	default:
		return nil, odbc.SQL_SUCCESS
	}

	return nil, odbc.SQL_SUCCESS

}

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

// Original source for function is Go source tree
// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest, src interface{}) error {
	// Common cases, without reflect.  Fall through.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			*d = s
			return nil
		case *[]byte:
			*d = []byte(s)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			*d = string(s)
			return nil
		case *interface{}:
			bcopy := make([]byte, len(s))
			copy(bcopy, s)
			*d = bcopy
			return nil
		case *[]byte:
			*d = s
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *[]byte:
			*d = nil
			return nil
		}
	}

	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = fmt.Sprintf("%v", src)
			return nil
		}
	case *interface{}:
		*d = src
		return nil
	}


	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if dv.Kind() == sv.Kind() {
		dv.Set(sv)
		return nil
	}

	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		} else {
			dv.Set(reflect.New(dv.Type().Elem()))
			return convertAssign(dv.Interface(), src)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	}

	return fmt.Errorf("unsupported driver -> Scan pair: %T -> %T", src, dest)
}

// Original source for function is Go source tree
func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	return fmt.Sprintf("%v", src)
}