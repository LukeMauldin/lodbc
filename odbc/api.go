package odbc

import ()

////Generate file command: C:\GoBuild\go\src\pkg\syscall\mksyscall_windows.pl api.go > apisys.go
////Console mode: mode con:cols=150

//sys   SQLAllocHandle(handleType SQLSMALLINT, inputHandle SQLHandle, outputHandle *SQLHandle) (ret SQLReturn) = odbc32.SQLAllocHandle
//sys   SQLSetEnvAttr(environmentHandle SQLHandle, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetEnvAttr
//sys   SQLDriverConnect(connectionHandle SQLHandle, windowHandle int, inConnString *SQLCHAR, inConnStringLength SQLSMALLINT, outConnString *SQLCHAR, outConnStringLength SQLSMALLINT, outConnStringPtr *SQLSMALLINT, driverCompletion SQLUSMALLINT) (ret SQLReturn) = odbc32.SQLDriverConnectW
//sys   SQLFreeHandle(handleType SQLSMALLINT, handle SQLHandle) (ret SQLReturn) = odbc32.SQLFreeHandle
//sys   SQLDisconnect(handle SQLHandle) (ret SQLReturn) = odbc32.SQLDisconnect
//sys   SQLCancel(statementHandle SQLHandle) (ret SQLReturn) = odbc32.SQLCancel
//sys   SQLExecDirect(statementHandle SQLHandle, statementText *SQLCHAR, textLength SQLINTEGER) (ret SQLReturn) = odbc32.SQLExecDirectW
//sys   SQLCloseCursor(statementHandle SQLHandle) (ret SQLReturn) = odbc32.SQLCloseCursor
//sys   SQLFetch(statementHandle SQLHandle) (ret SQLReturn) = odbc32.SQLFetch
//sys   SQLFetchScroll(statementHandle SQLHandle, fetchOrientation SQLSMALLINT, fetchOffset  SQLLEN) (ret SQLReturn)  = odbc32.SQLFetchScroll
//sys   SQLSetStmtAttr(statementHandle SQLHandle, attribute SQLINTEGER, valuePtr SQLPOINTER, stringLength SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetStmtAttr
//sys   SQLBindCol(statementHandle SQLHandle, columnNumber SQLUSMALLINT, targetType SQLSMALLINT, targetValuePtr SQLPOINTER, bufferLength SQLLEN, ind *SQLLEN) (ret SQLReturn) = odbc32.SQLBindCol
//sys   SQLSetConnectAttr(connectionHandle SQLHandle, attribute SQLINTEGER, valuePtr SQLPOINTER, bufferLength SQLINTEGER, stringLengthPtr *SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetConnectAttrW
//sys   SQLEndTran(handleType SQLSMALLINT, handle SQLHandle, completionType SQLSMALLINT) (ret SQLReturn) = odbc32.SQLEndTran
//sys   SQLBindParameter(statementHandle SQLHandle, parameterNumber SQLUSMALLINT, inputOutputType SQLSMALLINT, valueType CDataType, parameterType SQLDataType, columnSize SQLULEN, decimalDigits SQLSMALLINT, parameterValue SQLPOINTER, bufferLength SQLLEN, ind *SQLLEN) (ret SQLReturn) = odbc32.SQLBindParameter
//sys   SQLMoreResults(statementHandle SQLHandle) (ret SQLReturn) = odbc32.SQLMoreResults
//sys   SQLGetDescField(descriptorHandle SQLHandle, recNumber SQLSMALLINT, fieldIdentifier SQLSMALLINT, valuePtr SQLPOINTER, bufferLength SQLINTEGER, lengthPtr *SQLINTEGER) (ret SQLReturn) = odbc32.SQLGetDescField
//sys   SQLGetDescRec(descriptorHandle SQLHandle, recNumber SQLSMALLINT, name *SQLCHAR, bufferLength SQLSMALLINT, stringLengthPtr *SQLSMALLINT, typePtr *SQLSMALLINT, subTypePtr *SQLSMALLINT, lengthPtr *SQLLEN, precisionPtr *SQLSMALLINT, scalePtr *SQLSMALLINT, nullablePtr *SQLSMALLINT) (ret SQLReturn) = odbc32.SQLGetDescRecW
//sys   SQLGetDiagRec(handleType SQLSMALLINT, inputHandle SQLHandle, recNumber SQLSMALLINT, sqlState uintptr, nativeErrorPtr *SQLINTEGER, messageText uintptr, bufferLength SQLSMALLINT, textLengthPtr *SQLSMALLINT) (ret SQLReturn) = odbc32.SQLGetDiagRecW
//sys   SQLColAttribute(statementHandle SQLHandle, columnNumber SQLUSMALLINT, fieldIdentifier SQLColAttributeType, characterAttribute uintptr, bufferLength SQLSMALLINT, stringLengthPtr *SQLSMALLINT, numericAttributePtr *SQLLEN) (ret SQLReturn) = odbc32.SQLColAttributeW
//sys   SQLNumResultCols(statementHandle SQLHandle, columnCount *SQLSMALLINT)  (ret SQLReturn) = odbc32.SQLNumResultCols
//sys   SQLGetData(statementHandle SQLHandle, colNum SQLUSMALLINT, targetType CDataType, targetValuePtr uintptr, bufferLength SQLLEN, ind *SQLLEN) (ret SQLReturn) = odbc32.SQLGetData
//sys   SQLGetStmtAttr(statementHandle SQLHandle, attribute SQLINTEGER, valuePtr uintptr, bufferLength SQLINTEGER, stringLengthPtr *SQLINTEGER) (ret SQLReturn) = odbc32.SQLGetStmtAttr
//sys   SQLSetDescField(descriptorHandle SQLHandle, recNum SQLSMALLINT, fieldIdentifier SQLSMALLINT, valuePtr uintptr, bufferLength SQLINTEGER) (ret SQLReturn) = odbc32.SQLSetDescFieldW
