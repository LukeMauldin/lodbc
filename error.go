package lodbc

import (
	"fmt"
	"github.com/LukeMauldin/lodbc/odbc"
	"strings"
	"syscall"
	"unsafe"
)

// Error constants
const (
	sqlStateLength        = 5
	errorMaxMessageLength = 8000
)

// Contains all possible information about an ODBC error
// Exported so client programs can obtain additional error information
type StatusRecord struct {
	State       string
	NativeError int
	Message     string
	DriverInfo  string
}

// Convert StatusRecord to string
func (sr *StatusRecord) toString() string {
	if sr.DriverInfo != "" {
		return fmt.Sprintf("{%s} \n%s \n%s", sr.State, sr.DriverInfo, sr.Message)
	}
	return fmt.Sprintf("{%s} %s", sr.State, sr.Message)
}

// Contains a slice of the error(s) returned by ODBCError
// Exported so client programs can obtain additional error information
// Implements Error() interface
type ODBCError struct {
	StatusRecords []StatusRecord
}

// Implements Error() interface
func (e *ODBCError) Error() string {
	statusStrings := make([]string, len(e.StatusRecords))
	for i, sr := range e.StatusRecords {
		statusStrings[i] = sr.toString()
	}

	return strings.Join(statusStrings, "\n")
}

// Checks for SQL error
func isError(ret odbc.SQLReturn) bool {
	return !(ret == odbc.SQL_SUCCESS || ret == odbc.SQL_SUCCESS_WITH_INFO || ret == odbc.SQL_NO_DATA)
}

func errorEnvironment(handle odbc.SQLHandle) error {
	return handleError(odbc.SQL_HANDLE_ENV, handle, "")
}

func errorConnection(handle odbc.SQLHandle) error {
	return handleError(odbc.SQL_HANDLE_DBC, handle, "")
}

func errorStatement(handle odbc.SQLHandle, driverInfo string) error {
	return handleError(odbc.SQL_HANDLE_STMT, handle, driverInfo)
}

func handleError(handleType odbc.SQLSMALLINT, handle odbc.SQLHandle, driverInfo string) error {
	statusRecords := make([]StatusRecord, 0)
	if handle != 0 {
		for recNum := 1; ; recNum++ {
			sqlState := make([]uint16, sqlStateLength+1)
			var nativeError odbc.SQLINTEGER
			message := make([]uint16, errorMaxMessageLength+1)
			ret := odbc.SQLGetDiagRec(handleType, handle, odbc.SQLSMALLINT(recNum), uintptr(unsafe.Pointer(&sqlState[0])), &nativeError, uintptr(unsafe.Pointer(&message[0])), errorMaxMessageLength, nil)
			if ret == odbc.SQL_NO_DATA {
				break
			} else if !isError(ret) {
				sr := StatusRecord{State: syscall.UTF16ToString(sqlState), NativeError: int(nativeError), Message: syscall.UTF16ToString(message), DriverInfo: driverInfo}
				statusRecords = append(statusRecords, sr)
			} else {
				break
			}
		}
	}

	return &ODBCError{StatusRecords: statusRecords}
}
