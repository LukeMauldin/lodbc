package lodbc

import (
	"github.com/LukeMauldin/lodbc/odbc"
	"syscall"
	"unsafe"
)

// Metadata about each result column
type resultColumnDef struct {
	RecNum    uint16
	DataType  odbc.SQLDataType
	Length    int32
	Precision int32
	Scale     int32
	Name      string
}

// Build metadata for each result column
func buildResultColumnDefinitions(stmtHandle syscall.Handle, sqlStmt string) ([]resultColumnDef, odbc.SQLReturn) {
	
	//Get number of result columns
	var numColumns int16
	ret := odbc.SQLNumResultCols(stmtHandle, &numColumns)
	if isError(ret) {
		errorStatement(stmtHandle, sqlStmt)
	}

	resultColumnDefs := make([]resultColumnDef, 0, numColumns)
	for colNum, lNumColumns := uint16(1), uint16(numColumns); colNum <= lNumColumns; colNum++ {
		//Get odbc.SQL type
		var sqlType odbc.SQLLEN
		ret := odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_TYPE, 0, 0, nil, &sqlType)
		if isError(ret) {
			errorStatement(stmtHandle, sqlStmt)
		}

		//Get length
		var length odbc.SQLLEN
		ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_LENGTH, 0, 0, nil, &length)
		if isError(ret) {
			errorStatement(stmtHandle, sqlStmt)
		}

		//If the type is a CHAR or VARCHAR, add 4 to the length
		if sqlType == odbc.SQL_CHAR || sqlType == odbc.SQL_VARCHAR || sqlType == odbc.SQL_WCHAR || sqlType == odbc.SQL_WVARCHAR {
			length = length + 4
		}

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
		if sqlType == odbc.SQL_NUMERIC || sqlType == odbc.SQL_DECIMAL {
			ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_PRECISION, 0, 0, nil, &precision)
			if isError(ret) {
				errorStatement(stmtHandle, sqlStmt)
			}
		}

		//For numeric and decimal types, get the scale
		var scale odbc.SQLLEN
		if sqlType == odbc.SQL_NUMERIC || sqlType == odbc.SQL_DECIMAL {
			ret = odbc.SQLColAttribute(stmtHandle, odbc.SQLUSMALLINT(colNum), odbc.SQL_COLUMN_SCALE, 0, 0, nil, &scale)
			if isError(ret) {
				errorStatement(stmtHandle, sqlStmt)
			}
		}

		col := resultColumnDef{RecNum: colNum, DataType: odbc.SQLDataType(sqlType), Name: name, Length: int32(length), Precision: int32(precision), Scale: int32(scale)}
		resultColumnDefs = append(resultColumnDefs, col)
	}

	return resultColumnDefs, odbc.SQL_SUCCESS
}
