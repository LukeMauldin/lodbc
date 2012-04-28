package odbc

import ()

////Generate file command: C:\GoBuild\go\src\pkg\syscall\mksyscall_windows.pl api.go > apisys.go

//sys   SQLAllocHandle(handleType SQLHandle, inputHandle syscall.Handle, outputHandle *syscall.Handle) (ret SQLReturn) = odbc32.SQLAllocHandle
//sys   SQLSetEnvAttr(environmentHandle syscall.Handle, attribute SQLINTEGER, valuePtr int32, stringLength SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetEnvAttr
//sys   SQLDriverConnect(connectionHandle syscall.Handle, windowHandle int, inConnString *uint16, inConnStringLength int16, outConnString *int, outConnStringLength int16, outConnStringPtr *int16, driverCompletion uint16) (ret SQLReturn) = odbc32.SQLDriverConnectW
//sys   SQLFreeHandle(handleType SQLHandle, handle syscall.Handle) (ret SQLReturn) = odbc32.SQLFreeHandle
//sys   SQLDisconnect(handle syscall.Handle) (ret SQLReturn) = odbc32.SQLDisconnect
//sys   SQLCancel(statementHandle syscall.Handle) (ret SQLReturn) = odbc32.SQLCancel
//sys   SQLExecDirect(statementHandle syscall.Handle, statementText *uint16, textLength int) (ret SQLReturn) = odbc32.SQLExecDirectW
//sys   SQLCloseCursor(statementHandle syscall.Handle) (ret SQLReturn) = odbc32.SQLCloseCursor
//sys   SQLFetch(statementHandle syscall.Handle) (ret SQLReturn) = odbc32.SQLFetch
//sys   SQLFetchScroll(statementHandle syscall.Handle, fetchOrientation int16, fetchOffset int32) (ret SQLReturn)  = odbc32.SQLFetchScroll
//sys   SQLSetStmtAttr(statementHandle syscall.Handle, attribute int, valuePtr int32, stringLength int) (ret SQLReturn) = odbc32.SQLSetStmtAttr
//sys   SQLBindCol(statementHandle syscall.Handle, columnNumber uint16, targetType CDataType, targetValuePtr uintptr, bufferLength SQLLEN, ind *SQLValueIndicator) (ret SQLReturn) = odbc32.SQLBindCol
//sys   SQLSetConnectAttr(connectionHandle syscall.Handle, attribute SQLConnectionAttribute, valuePtr int32, bufferLength SQLINTEGER, stringLengthPtr *SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetConnectAttrW
//sys   SQLEndTran(handleType SQLHandle, handle syscall.Handle, completionType SQLTransactionOption) (ret SQLReturn) = odbc32.SQLEndTran
//sys   SQLBindParameter(statementHandle syscall.Handle, parameterNumber SQLUSMALLINT, inputOutputType SQLBindParameterType, valueType CDataType, parameterType SQLDataType, columnSize SQLULEN, decimalDigits SQLSMALLINT, parameterValue SQLPOINTER, bufferLength SQLLEN, ind *SQLValueIndicator) (ret SQLReturn) = odbc32.SQLBindParameter
//sys   SQLMoreResults(statementHandle syscall.Handle) (ret SQLReturn) = odbc32.SQLMoreResults
//sys   SQLGetDescField(descriptorHandle syscall.Handle, recNumber int16, fieldIdentifier int16, valuePtr uintptr, bufferLength int, lengthPtr *int) (ret SQLReturn) = odbc32.SQLGetDescField
//sys   SQLGetDescRec(descriptorHandle syscall.Handle, recNumber int16, name *uint16, bufferLength int16, stringLengthPtr *int16, typePtr *int16, subTypePtr *int16, lengthPtr *int, precisionPtr *int16, scalePtr *int16, nullablePtr *int16) (ret SQLReturn) = odbc32.SQLGetDescRecW
//sys   SQLGetDiagRec(handleType SQLHandle, inputHandle syscall.Handle, recNumber int16, sqlState uintptr, nativeErrorPtr *int, messageText uintptr, bufferLength int16, textLengthPtr *int16) (ret SQLReturn) = odbc32.SQLGetDiagRecW
//sys   SQLColAttribute(statementHandle syscall.Handle, columnNumber uint16, fieldIdentifier SQLDescriptor, characterAttribute uintptr, bufferLength SQLSMALLINT, stringLengthPtr *int16, numericAttributePtr *int32) (ret SQLReturn) = odbc32.SQLColAttributeW
//sys   SQLNumResultCols(statementHandle syscall.Handle, columnCount *int16)  (ret SQLReturn) = odbc32.SQLNumResultCols
//sys   SQLGetData(statementHandle syscall.Handle, colNum uint16 , targetType CDataType, targetValuePtr uintptr, bufferLength SQLLEN, ind *SQLValueIndicator) (ret SQLReturn) = odbc32.SQLGetData
//sys   SQLGetStmtAttr(statementHandle syscall.Handle, attribute SQLStatementAttribute, valuePtr uintptr, bufferLength SQLINTEGER, stringLengthPtr *SQLINTEGER) (ret SQLReturn) = odbc32.SQLGetStmtAttr
//sys   SQLSetDescField(descriptorHandle syscall.Handle, recNum SQLSMALLINT, fieldIdentifier SQLSMALLINT, valuePtr int32, bufferLength SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetDescFieldW