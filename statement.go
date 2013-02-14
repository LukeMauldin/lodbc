package lodbc

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"fmt"
	"github.com/LukeMauldin/lodbc/odbc"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type statement struct {
	//Statement handle
	handle odbc.SQLHandle

	//Statement descriptor handle
	stmtDescHandle odbc.SQLHandle

	//Current executing sql statement
	sqlStmt string

	//Active query for the statement
	rows driver.Rows

	//Owning connection
	conn *connection

	//Is closed -- allows Close() to be called multiple times without error
	isClosed bool

	//Array to store bind parameter values to be sure they stay in scope
	bindValues []interface{}

	//SQL statement options
	queryOptions []QueryOption
}

func (stmt *statement) bindInt(index int, value int, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_LONG, odbc.SQL_INTEGER, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*int))), 0, nil)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *statement) bindInt64(index int, value int64, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_LONG, odbc.SQL_BIGINT, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*int64))), 0, nil)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *statement) bindBool(index int, value bool, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_BIT, odbc.SQL_BIT, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*bool))), 0, nil)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *statement) bindNumeric(index int, value float64, precision int, scale int, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_DOUBLE, odbc.SQL_DOUBLE, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*float64))), 0, nil)
	/* Must convert to SQL_NUMERIC_STRUCT for decimal to work - http://support.microsoft.com/kb/181254
	 ret := odbc.SQLBindParameter(stmt.handle, uint16(index), direction.SQLBindParameterType(), odbc.SQL_C_NUMERIC, odbc.SQL_DECIMAL, uint64(precision), int16(scale), uintptr(unsafe.Pointer(&bindVal)), 0, nil)
	odbc.SQLSetDescField(stmt.stmtDescHandle, odbc.SQLSMALLINT(index), odbc.SQL_DESC_TYPE, odbc.SQL_NUMERIC, 0)
	odbc.SQLSetDescField(stmt.stmtDescHandle, odbc.SQLSMALLINT(index), odbc.SQL_DESC_PRECISION, int32(precision), 0)
	odbc.SQLSetDescField(stmt.stmtDescHandle, odbc.SQLSMALLINT(index), odbc.SQL_DESC_SCALE, int32(scale), 0) */
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *statement) bindDate(index int, value time.Time, direction ParameterDirection) error {
	var bindVal odbc.SQL_DATE_STRUCT
	bindVal.Year = odbc.SQLSMALLINT(value.Year())
	bindVal.Month = odbc.SQLUSMALLINT(value.Month())
	bindVal.Day = odbc.SQLUSMALLINT(value.Day())

	stmt.bindValues[index] = &bindVal
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_DATE, odbc.SQL_DATE, 10, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*odbc.SQL_DATE_STRUCT))), 6, nil)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, bindVal))
	}
	return nil
}

func (stmt *statement) bindDateTime(index int, value time.Time, direction ParameterDirection) error {
	var bindVal odbc.SQL_TIMESTAMP_STRUCT
	bindVal.Year = odbc.SQLSMALLINT(value.Year())
	bindVal.Month = odbc.SQLUSMALLINT(value.Month())
	bindVal.Day = odbc.SQLUSMALLINT(value.Day())
	bindVal.Hour = odbc.SQLUSMALLINT(value.Hour())
	bindVal.Minute = odbc.SQLUSMALLINT(value.Minute())
	bindVal.Second = odbc.SQLUSMALLINT(value.Second())

	stmt.bindValues[index] = &bindVal
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_TIMESTAMP, odbc.SQL_TIMESTAMP, 23, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*odbc.SQL_TIMESTAMP_STRUCT))), 16, nil)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, bindVal))
	}
	return nil
}

func (stmt *statement) bindString(index int, value string, length int, direction ParameterDirection) error {
	if length == 0 {
		length = len(value)
	}
	stmt.bindValues[index] = syscall.StringToUTF16(value)
	var sqlType odbc.SQLDataType
	if length < 4000 {
		sqlType = odbc.SQL_VARCHAR
	} else {
		sqlType = odbc.SQL_LONGVARCHAR
	}
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_WCHAR, sqlType, odbc.SQLULEN(length), 0, odbc.SQLPOINTER(unsafe.Pointer(&stmt.bindValues[index].([]uint16)[0])), 0, nil)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *statement) bindByteArray(index int, value []byte, direction ParameterDirection) error {
	// Store both value and lenght, because we need a pointer to the lenght in
	// the last parameter of SQLBindParamter. Otherwise the data is assumed to
	// be a null terminated string.
	bindVal := &struct {
		value  []byte
		length int
	}{
		value,
		len(value),
	}
	sqlType := odbc.SQL_VARBINARY
	if bindVal.length > 4000 {
		sqlType = odbc.SQL_LONGVARBINARY
	}

	// Protect against index out of range on &bindVal.value[0] when value is zero-length.
	// We can't pass NULL to SQLBindParameter so this is needed, it will still
	// write a zero length value to the database since the length parameter is
	// zero.
	if bindVal.length == 0 {
		bindVal.value = []byte{'\x00'}
	}

	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_BINARY, sqlType, odbc.SQLULEN(bindVal.length), 0, odbc.SQLPOINTER(unsafe.Pointer(&bindVal.value[0])), 0, (*odbc.SQLLEN)(unsafe.Pointer(&bindVal.length)))
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *statement) bindNull(index int, direction ParameterDirection) error {
	return stmt.bindNullParam(index, odbc.SQL_WCHAR, direction)
}

func (stmt *statement) bindNullParam(index int, paramType odbc.SQLDataType, direction ParameterDirection) error {
	nullDataInd := odbc.SQL_NULL_DATA
	stmt.bindValues[index] = &nullDataInd
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_DEFAULT, paramType, 1, 0, 0, 0, &nullDataInd)
	if isError(ret) {
		return errorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: nil", index))
	}
	return nil
}

func (stmt *statement) Close() error {

	//Verify that stmtHandle is valid
	if stmt.handle == 0 {
		return nil
	}

	//Verify that statement has not already been closed
	if stmt.isClosed {
		return nil
	}

	var err error

	//Close any open rows
	if stmt.rows != nil {
		err = stmt.rows.Close()
	}

	//Clear any bind values
	stmt.bindValues = nil

	//Free the statement handle
	ret := odbc.SQLFreeHandle(odbc.SQL_HANDLE_STMT, stmt.handle)
	if isError(ret) {
		err = errorStatement(stmt.handle, stmt.sqlStmt)
	}

	//Mark the statement as closed with the connection
	stmt.conn.closeStatement(stmt)

	//Clear the handles
	stmt.handle = 0
	stmt.stmtDescHandle = 0

	//Clear the finalizer
	runtime.SetFinalizer(stmt, nil)

	//Mark the rows as closed
	stmt.isClosed = true

	//Return any error
	if err != nil {
		return err
	}

	return nil
}

func (stmt *statement) Query(args []driver.Value) (driver.Rows, error) {
	//Clear any existing bind values
	stmt.bindValues = make([]interface{}, len(args)+1)

	//Bind the parameters
	bindParameters, err := stmt.convertToBindParameters(args)
	if err != nil {
		return nil, err
	}
	stmt.bindParameters(bindParameters)

	//If rows is not nil, close rows and set to nil
	if stmt.rows != nil {
		stmt.rows.Close()
		stmt.rows = nil
	}

	//Execute SQL statement
	sqlStmtSqlPtr := (*odbc.SQLCHAR)(unsafe.Pointer(syscall.StringToUTF16Ptr(stmt.sqlStmt)))
	ret := odbc.SQLExecDirect(stmt.handle, sqlStmtSqlPtr, odbc.SQL_NTS)
	if isError(ret) {
		return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\nBind Values: %v", stmt.sqlStmt, stmt.formatBindValues()))
	}

	//Get row descriptor handle
	var descRowHandle odbc.SQLHandle
	ret = odbc.SQLGetStmtAttr(stmt.handle, odbc.SQL_ATTR_APP_ROW_DESC, uintptr(unsafe.Pointer(&descRowHandle)), 0, nil)
	if isError(ret) {
		return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\nBind Values: %v", stmt.sqlStmt, stmt.formatBindValues()))
	}

	//Check to see if the query option ResultSetNum was passed and if so, iterate through result sets
	optionValue, optionFound := getOptionValue(stmt.queryOptions, ResultSetNum)
	if optionFound {
		for counter, resultSetNum := 0, int(optionValue.(float64)); counter < resultSetNum; counter++ {
			ret := odbc.SQLMoreResults(stmt.handle)
			if isError(ret) {
				return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v", stmt.sqlStmt))
			}
		}
	} else {
		//If query option ResultSetNum was not passed, iterate through result sets until at least one column is found
		for {
			var numColumns odbc.SQLSMALLINT
			ret := odbc.SQLNumResultCols(stmt.handle, &numColumns)
			if isError(ret) {
				return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v", stmt.sqlStmt))
			}
			if numColumns > 0 {
				break
			} else {
				ret := odbc.SQLMoreResults(stmt.handle)
				if isError(ret) {
					return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v", stmt.sqlStmt))
				}
			}
		}
	}

	//Get definition of result columns
	resultColumnDefs, ret := buildResultColumnDefinitions(stmt.handle, stmt.sqlStmt)
	if isError(ret) {
		return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\nBind Values: %v", stmt.sqlStmt, stmt.formatBindValues()))
	}

	//Make a slice of the column names
	columnNames := make([]string, len(resultColumnDefs))
	for index, resultCol := range resultColumnDefs {
		columnNames[index] = fmt.Sprint(resultCol.Name)
	}

	//Create rows
	stmt.rows = &rows{handle: stmt.handle, descHandle: descRowHandle, isBeforeFirst: true, resultColumnDefs: resultColumnDefs, resultColumnNames: columnNames, sqlStmt: stmt.sqlStmt}

	//Add a finalizer
	runtime.SetFinalizer(stmt.rows, (*rows).Close)

	return stmt.rows, nil
}

func (stmt *statement) Exec(args []driver.Value) (driver.Result, error) {
	//Clear any existing bind values
	stmt.bindValues = make([]interface{}, len(args)+1)

	//Bind the parameters
	bindParameters, err := stmt.convertToBindParameters(args)
	if err != nil {
		return nil, err
	}
	stmt.bindParameters(bindParameters)

	//If rows is not nil, close rows and set to nil
	if stmt.rows != nil {
		stmt.rows.Close()
		stmt.rows = nil
	}

	//Execute SQL statement
	sqlStmtSqlPtr := (*odbc.SQLCHAR)(unsafe.Pointer(syscall.StringToUTF16Ptr(stmt.sqlStmt)))
	ret := odbc.SQLExecDirect(stmt.handle, sqlStmtSqlPtr, odbc.SQL_NTS)
	if isError(ret) {
		return nil, errorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\n Bind Values: %v", stmt.sqlStmt, stmt.formatBindValues()))
	}

	return driver.ResultNoRows, nil
}

func (stmt *statement) NumInput() int {
	return -1 //No checking by the driver
}

func (stmt *statement) convertToBindParameters(args []driver.Value) ([]BindParameter, error) {
	bindParameters := make([]BindParameter, len(args))
	//Check each item in args and see if it is an encoded byte array or a driver.Value
	for index, arg := range args {
		if encodedBytes, ok := arg.([]byte); ok {
			//If arg is an encoded byte array, attempt to decode into a bindParameter
			decodedBuffer := bytes.NewBuffer(encodedBytes)
			dec := gob.NewDecoder(decodedBuffer)
			var bindParameter BindParameter
			err := dec.Decode(&bindParameter)
			if err == nil {
				bindParameters[index] = bindParameter
				continue
			}
		}

		bindParameters[index] = BindParameter{Data: arg}
	}
	return bindParameters, nil
}

func (stmt *statement) bindParameters(parameters []BindParameter) error {
	//Call bind statements based on the type of the parameter
	for index, parameter := range parameters {
		//Bind a null parameter
		if isNil(parameter.Data) {
			err := stmt.bindNull(index+1, parameter.Direction)
			if err != nil {
				return err
			}
			continue
		}

		//Flatten out pointers
		elemValue := reflect.Indirect(reflect.ValueOf(parameter.Data)).Interface()

		switch value := elemValue.(type) {
		case nil:
			err := stmt.bindNull(index+1, parameter.Direction)
			if err != nil {
				return err
			}
		case bool:
			err := stmt.bindBool(index+1, value, parameter.Direction)
			if err != nil {
				return err
			}
		case int:
			err := stmt.bindInt(index+1, value, parameter.Direction)
			if err != nil {
				return err
			}
		case int64:
			err := stmt.bindInt64(index+1, value, parameter.Direction)
			if err != nil {
				return err
			}
		case float64:
			err := stmt.bindNumeric(index+1, value, parameter.Precision, parameter.Scale, parameter.Direction)
			if err != nil {
				return err
			}
		case string:
			err := stmt.bindString(index+1, value, parameter.Length, parameter.Direction)
			if err != nil {
				return err
			}
		case []byte:
			err := stmt.bindByteArray(index+1, value, parameter.Direction)
			if err != nil {
				return err
			}
		case time.Time:
			if parameter.DateOnly {
				err := stmt.bindDate(index+1, value, parameter.Direction)
				if err != nil {
					return err
				}
			} else {
				err := stmt.bindDateTime(index+1, value, parameter.Direction)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("Error binding parameter number: %v.  Parameter type not supported: %T", index+1, parameter.Data)
		}
	}

	return nil
}

func (stmt *statement) formatBindValues() string {
	strValues := make([]string, 0, len(stmt.bindValues))
	for index, bvalue := range stmt.bindValues {
		//Skip 0 index
		if index == 0 {
			continue
		}
		if bvalue == nil {
			strValues = append(strValues, fmt.Sprintf("%v: <nil>", index))
		} else {
			switch val := bvalue.(type) {
			case *int, *int64, *bool, *float64, *odbc.SQL_DATE_STRUCT, *odbc.SQL_TIMESTAMP_STRUCT:
				refValue := reflect.ValueOf(val)
				interfaceValue := reflect.Indirect(refValue).Interface()
				name := reflect.TypeOf(interfaceValue).Name()
				strValues = append(strValues, fmt.Sprintf("%v: <%v> {%v}", index, name, interfaceValue))
			case []uint16:
				str := syscall.UTF16ToString(val)
				strValues = append(strValues, fmt.Sprintf("%v: <string> {%v}", index, str))
			case *odbc.SQLLEN:
				if *val == odbc.SQL_NULL_DATA {
					strValues = append(strValues, fmt.Sprintf("%v: {NULL}", index))
				} else {
					strValues = append(strValues, fmt.Sprintf("%v: <SQLLEN> %v", index, val))
				}
			default:
				strValues = append(strValues, fmt.Sprintf("%v: Unknown type: <%t>", index, val))
			}

		}
	}

	return strings.Join(strValues, ", ")
}
