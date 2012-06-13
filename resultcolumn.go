package lodbc

import (
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
	"unsafe"
)

// Metadata about each result column
type resultColumnDef struct {
	RecNum    odbc.SQLSMALLINT
	DataType  odbc.SQLDataType
	Precision odbc.SQLLEN
	Scale     odbc.SQLLEN
	Name      string
}

// Build metadata for each result column
func buildResultColumnDefinitions(stmtHandle odbc.SQLHandle, sqlStmt string) ([]resultColumnDef, odbc.SQLReturn) {

	//Get number of result columns
	var numColumns odbc.SQLSMALLINT
	ret := odbc.SQLNumResultCols(stmtHandle, &numColumns)
	if isError(ret) {
		errorStatement(stmtHandle, sqlStmt)
	}

	resultColumnDefs := make([]resultColumnDef, 0, numColumns)
	for colNum, lNumColumns := odbc.SQLSMALLINT(1), numColumns; colNum <= lNumColumns; colNum++ {
		//Get odbc.SQL type
		var sqlType odbc.SQLLEN
		ret := odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_TYPE, 0, 0, nil, &sqlType)
		if isError(ret) {
			errorStatement(stmtHandle, sqlStmt)
		}

		/* Disabled because it is no longer needed
		//Get length
		var length odbc.SQLLEN
		ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_LENGTH, 0, 0, nil, &length)
		if isError(ret) {
			errorStatement(stmtHandle, sqlStmt)
		}

		//If the type is a CHAR or VARCHAR, add 4 to the length
		if odbc.SQLDataType(sqlType) == odbc.SQL_CHAR || odbc.SQLDataType(sqlType)  == odbc.SQL_VARCHAR || odbc.SQLDataType(sqlType)  == odbc.SQL_WCHAR || odbc.SQLDataType(sqlType)  == odbc.SQL_WVARCHAR {
			length = length + 4
		} */

		//Get name
		const namelength = 1000
		nameArr := make([]uint16, namelength)
		ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_DESC_LABEL, uintptr(unsafe.Pointer(&nameArr[0])), namelength, nil, nil)
		if isError(ret) {
			errorStatement(stmtHandle, sqlStmt)
		}
		name := syscall.UTF16ToString(nameArr)

		//For numeric and decimal types, get the precision
		var precision odbc.SQLLEN
		if odbc.SQLDataType(sqlType) == odbc.SQL_NUMERIC || odbc.SQLDataType(sqlType) == odbc.SQL_DECIMAL {
			ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_PRECISION, 0, 0, nil, &precision)
			if isError(ret) {
				errorStatement(stmtHandle, sqlStmt)
			}
		}

		//For numeric and decimal types, get the scale
		var scale odbc.SQLLEN
		if odbc.SQLDataType(sqlType) == odbc.SQL_NUMERIC || odbc.SQLDataType(sqlType) == odbc.SQL_DECIMAL {
			ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_SCALE, 0, 0, nil, &scale)
			if isError(ret) {
				errorStatement(stmtHandle, sqlStmt)
			}
		}

		col := resultColumnDef{RecNum: colNum, DataType: odbc.SQLDataType(sqlType), Name: name, Precision: precision, Scale: scale}
		resultColumnDefs = append(resultColumnDefs, col)
	}

	return resultColumnDefs, odbc.SQL_SUCCESS
}
