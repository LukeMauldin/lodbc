package lodbc

import (
	"github.com/LukeMauldin/lodbc"
	"syscall"
	 "time"
	"unsafe"
	"fmt"
	"reflect"
	"strings"
)

type Statement struct {
	//Statement handle
	handle syscall.Handle
	
	//Statement descriptor handle
	stmtDescHandle syscall.Handle
	
	//Active query for the statement
	rows IRows
	
	//Is closed -- allows Close() to be called multiple times without error
	isClosed bool
	
	//Current executing sql statement
	sqlStmt string
	
	//Store bind parameter values in a map to be sure they stay in scope
	//bindValues map[int]interface{}
	
	//Array to store bind parameter values to be sure they stay in scope
	bindValues []interface{}
}

func (stmt *Statement) bindInt(index int, value int, direction ParameterDirection) error {
	stmt.bindValues[index] = &value	
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_LONG, odbc.SQL_INTEGER, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*int))), 0, nil)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
} 


func (stmt *Statement) bindInt64(index int, value int64, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_LONG, odbc.SQL_BIGINT, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*int64))), 0, nil)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *Statement) bindBool(index int, value bool, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_BIT, odbc.SQL_BIT, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*bool))), 0, nil)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *Statement) bindNumeric(index int, value float64, precision int, scale int, direction ParameterDirection) error {
	stmt.bindValues[index] = &value
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_DOUBLE, odbc.SQL_DOUBLE, 0, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*float64))), 0, nil)
	/* Must convert to SQL_NUMERIC_STRUCT for decimal to work - http://support.microsoft.com/kb/181254
	 ret := odbc.SQLBindParameter(stmt.handle, uint16(index), direction.SQLBindParameterType(), odbc.SQL_C_NUMERIC, odbc.SQL_DECIMAL, uint64(precision), int16(scale), uintptr(unsafe.Pointer(&bindVal)), 0, nil)
	odbc.SQLSetDescField(stmt.stmtDescHandle, odbc.SQLSMALLINT(index), odbc.SQL_DESC_TYPE, odbc.SQL_NUMERIC, 0)
	odbc.SQLSetDescField(stmt.stmtDescHandle, odbc.SQLSMALLINT(index), odbc.SQL_DESC_PRECISION, int32(precision), 0)
	odbc.SQLSetDescField(stmt.stmtDescHandle, odbc.SQLSMALLINT(index), odbc.SQL_DESC_SCALE, int32(scale), 0) */
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *Statement) bindDate(index int, value time.Time, direction ParameterDirection) error {
	var bindVal odbc.SQL_DATE_STRUCT
	bindVal.Year = odbc.SQLSMALLINT(value.Year())
	bindVal.Month = odbc.SQLUSMALLINT(value.Month())
	bindVal.Day = odbc.SQLUSMALLINT(value.Day())
	
	stmt.bindValues[index] = &bindVal
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_DATE, odbc.SQL_DATE, 10, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*odbc.SQL_DATE_STRUCT))), 6, nil)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, bindVal))
	}
	return nil
}

func (stmt *Statement) bindDateTime(index int, value time.Time, direction ParameterDirection) error {
	var bindVal odbc.SQL_TIMESTAMP_STRUCT
	bindVal.Year = odbc.SQLSMALLINT(value.Year())
	bindVal.Month = odbc.SQLUSMALLINT(value.Month())
	bindVal.Day = odbc.SQLUSMALLINT(value.Day())
	bindVal.Hour = odbc.SQLUSMALLINT(value.Hour())
	bindVal.Minute = odbc.SQLUSMALLINT(value.Minute())
	bindVal.Second = odbc.SQLUSMALLINT(value.Second())

	stmt.bindValues[index] = &bindVal
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_TIMESTAMP, odbc.SQL_TIMESTAMP, 23, 0, odbc.SQLPOINTER(unsafe.Pointer(stmt.bindValues[index].(*odbc.SQL_TIMESTAMP_STRUCT))), 16, nil)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, bindVal))
	}
	return nil
}

func (stmt *Statement) bindString(index int, value string, length int, direction ParameterDirection) error {
	if length == 0 {
		length = len(value)
	}
	stmt.bindValues[index] = syscall.StringToUTF16(value)
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_WCHAR, odbc.SQL_VARCHAR, odbc.SQLULEN(length), 0, odbc.SQLPOINTER(unsafe.Pointer(&stmt.bindValues[index].([]uint16)[0])), 0, nil)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: %v", index, value))
	}
	return nil
}

func (stmt *Statement) bindNull(index int, direction ParameterDirection) error {
	return stmt.bindNullParam(index, odbc.SQL_WCHAR, direction)
}

func (stmt *Statement) bindNullParam(index int, paramType odbc.SQLDataType, direction ParameterDirection) error {
	var nullDataInd odbc.SQLValueIndicator
	nullDataInd = odbc.SQL_NULL_DATA
	ret := odbc.SQLBindParameter(stmt.handle, odbc.SQLUSMALLINT(index), direction.SQLBindParameterType(), odbc.SQL_C_DEFAULT, paramType, 1, 0, 0, 0, &nullDataInd)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("Bind index: %v, Value: nil", index))
	}
	return nil
}

func (stmt *Statement) Close() error {

	//Verify that stmtHandle is valid
	if stmt.handle == 0 {
		return nil
	}
	
	//Verify that statement has not already been closed
	if stmt.isClosed {
		return nil
	}

	var err error
	isError := false

	//Close any open rows
	if stmt.rows != nil {
		err = stmt.rows.Close()
		isError = true
	}
	
	//Clear any bind values
	stmt.bindValues = nil

	//Free the statement handle
	ret := odbc.SQLFreeHandle(odbc.SQL_HANDLE_STMT, stmt.handle)
	if IsError(ret) {
		err = ErrorStatement(stmt.handle, "")
		isError = true
	}

	//Return any error
	if isError {
		return err
	}
	
	//Mark the rows as closed
	stmt.isClosed = true
	
	return nil
}

func (stmt *Statement) Query(query string) (IRows, error) {
	//If rows is not nil, close rows and set to nil
	if stmt.rows != nil {
		stmt.rows.Close()
		stmt.rows = nil
	}
	
	//Store the SQL statement being executed
	stmt.sqlStmt = query

	//Execute SQL statement
	ret := odbc.SQLExecDirect(stmt.handle, syscall.StringToUTF16Ptr(query), odbc.SQL_NTS)
	if IsError(ret) {		
		return nil, ErrorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\nBind Values: %v", query, stmt.formatBindValues()))
	}

	//Get row descriptor handle
	var descRowHandle syscall.Handle
	ret = odbc.SQLGetStmtAttr(stmt.handle, odbc.SQL_ATTR_APP_ROW_DESC, uintptr(unsafe.Pointer(&descRowHandle)), 0, nil)
	if IsError(ret) {
		return nil, ErrorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\nBind Values: %v", query, stmt.formatBindValues()))
	}

	//Get definition of result columns
	resultColumnDefs, ret := stmt.getResultColumnDefintion()
	if IsError(ret) {
		return nil, ErrorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\nBind Values: %v", query, stmt.formatBindValues()))
	}

	//Create rows
	stmt.rows = &Rows{handle: stmt.handle, descHandle: descRowHandle, isBeforeFirst: true, ResultColumnDefs: resultColumnDefs, sqlStmt: query}
	
	return stmt.rows, nil
}

func (stmt *Statement) QueryWithParams(query string, parameters ...BindParameter) (IRows, error) {
	//Clear any existing bind values
	stmt.bindValues = make([]interface{}, len(parameters) + 1)

	//Bind the parameters
	stmt.bindParameters(parameters...)
	
	//Execute the query
	return stmt.Query(query)
}

func (stmt *Statement) Exec(query string) error {
	//If rows is not nil, close rows and set to nil
	if stmt.rows != nil {
		stmt.rows.Close()
		stmt.rows = nil
	}
	
	//Store the SQL statement being executed
	stmt.sqlStmt = query

	//Execute SQL statement
	ret := odbc.SQLExecDirect(stmt.handle, syscall.StringToUTF16Ptr(query), odbc.SQL_NTS)
	if IsError(ret) {
		return ErrorStatement(stmt.handle, fmt.Sprintf("SQL Stmt: %v\n Bind Values: %v", query, stmt.formatBindValues()))
	}
	
	return nil
}

func (stmt *Statement) ExecWithParams(query string, parameters ...BindParameter) (error) {
	//Clear any existing bind values
	stmt.bindValues = make([]interface{}, len(parameters) + 1)

	//Bind the parameters
	stmt.bindParameters(parameters...)
	
	//Execute the statement
	return stmt.Exec(query)
}

func (stmt *Statement) bindParameters(parameters ...BindParameter) error {
//Call bind statements based on the type of the parameter
	for index, parameter := range parameters {
		if isNil(parameter.Value) {
			err := stmt.bindNull(index + 1, parameter.Direction)
			if err != nil {
				return err
			}
			continue
		}
		switch value := parameter.Value.(type) {	    
			case nil:
				err := stmt.bindNull(index + 1, parameter.Direction)
				if err != nil {
					return err
				}
			case bool:
				err := stmt.bindBool(index + 1, value, parameter.Direction)
				if err != nil {
					return err
				}
			case *bool:
				err := stmt.bindBool(index + 1, *value, parameter.Direction)
				if err != nil {
					return err
				}				
			case int:
				err := stmt.bindInt(index + 1, value, parameter.Direction)		
				if err != nil {
					return err
				}	
			case *int:
				err := stmt.bindInt(index + 1, *value, parameter.Direction)		
				if err != nil {
					return err
				}	
			case int64:
				err := stmt.bindInt64(index + 1, value, parameter.Direction)
				if err != nil {
					return err
				}
			case *int64:
				err := stmt.bindInt64(index + 1, *value, parameter.Direction)
				if err != nil {
					return err
				}
			case float64:
				err := stmt.bindNumeric(index + 1, value, parameter.Precision, parameter.Scale, parameter.Direction)
				if err != nil {
					return err
				}
			case *float64:
				err := stmt.bindNumeric(index + 1, *value, parameter.Precision, parameter.Scale, parameter.Direction)
				if err != nil {
					return err
				}
			case string:
				err := stmt.bindString(index + 1, value, parameter.Length, parameter.Direction)
				if err != nil {
					return err
				}
			case *string:
				err := stmt.bindString(index + 1, *value, parameter.Length, parameter.Direction)
				if err != nil {
					return err
				}
			case time.Time:
				if parameter.DateOnly {
					err := stmt.bindDate(index + 1, value, parameter.Direction)
					if err != nil {
						return err
					}
				} else {
					err := stmt.bindDateTime(index + 1, value, parameter.Direction)
						if err != nil {
							return err
						}
				}				
			case *time.Time:
				if parameter.DateOnly {
					err := stmt.bindDate(index + 1, *value, parameter.Direction)
					if err != nil {
						return err
					}
				} else {
					err := stmt.bindDateTime(index + 1, *value, parameter.Direction)
						if err != nil {
							return err
						}
				} 
			default:
				return fmt.Errorf("Error binding parameter number: %v.  Parameter type not supported: %T", index + 1, parameter.Value)  				
		}
	}
	
	return nil
}

func (stmt *Statement) getResultColumnDefintion() ([]ResultColumnDef, odbc.SQLReturn) {
	//Get number of result columns
	var numColumns int16
	ret := odbc.SQLNumResultCols(stmt.handle, &numColumns)
	if IsError(ret) {
		ErrorStatement(stmt.handle, stmt.sqlStmt)
	}

	resultColumnDefs := make([]ResultColumnDef, 0, numColumns)
	for colNum, lNumColumns := uint16(1), uint16(numColumns); colNum <= lNumColumns; colNum++ {
		//Get odbc.SQL type
		var sqlType int32
		ret := odbc.SQLColAttribute(stmt.handle, colNum, odbc.SQL_DESC_TYPE, 0, 0, nil, &sqlType)
		if IsError(ret) {
			ErrorStatement(stmt.handle, stmt.sqlStmt)
		}

		//Get length
		var length int32
		ret = odbc.SQLColAttribute(stmt.handle, colNum, odbc.SQL_DESC_LENGTH, 0, 0, nil, &length)
		if IsError(ret) {
			ErrorStatement(stmt.handle, stmt.sqlStmt)
		}
		
		//If the type is a CHAR or VARCHAR, add 4 to the length
		if sqlType == odbc.SQL_CHAR || sqlType ==  odbc.SQL_VARCHAR || sqlType == odbc.SQL_WCHAR || sqlType == odbc.SQL_WVARCHAR {
			length = length + 4
		}

		//For numeric and decimal types, get the precision
		var precision int32
		if sqlType == odbc.SQL_NUMERIC || sqlType == odbc.SQL_DECIMAL {
			ret = odbc.SQLColAttribute(stmt.handle, colNum, odbc.SQL_DESC_PRECISION, 0, 0, nil, &precision)
			if IsError(ret) {
				ErrorStatement(stmt.handle, stmt.sqlStmt)
			}
		}

		//For numeric and decimal types, get the scale
		var scale int32
		if sqlType == odbc.SQL_NUMERIC || sqlType == odbc.SQL_DECIMAL {
			ret = odbc.SQLColAttribute(stmt.handle, colNum, odbc.SQL_DESC_SCALE, 0, 0, nil, &scale)
			if IsError(ret) {
				ErrorStatement(stmt.handle, stmt.sqlStmt)
			}
		}

		resultColumnDef := ResultColumnDef{RecNum: colNum, DataType: odbc.SQLDataType(sqlType), Length: length, Precision: precision, Scale: scale}
		resultColumnDefs = append(resultColumnDefs, resultColumnDef)
	}

	return resultColumnDefs, odbc.SQL_SUCCESS
}

func (stmt *Statement) formatBindValues() string {
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
				default:
					strValues = append(strValues, fmt.Sprintf("%v: Unknown type: <%t>", index, val))	
			}
			
		}
	} 
	
	return strings.Join(strValues, ", ")
} 

//Checks the type v for nil
func isNil(v interface{}) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return val.IsNil()
	}
	return false
}
