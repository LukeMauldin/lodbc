package lodbc

import (
	"syscall"
	"fmt"
	"unsafe"
	"github.com/LukeMauldin/go-lodbc/odbc"
	"strings"
)

const SQLStateLength = 5
const ErrorMaxMessageLength = 8000

type StatusRecord struct {
	State       string
	NativeError int
	Message     string
	DriverInfo  string
}

func (sr *StatusRecord) toString() string {
	if sr.DriverInfo != "" {
		return fmt.Sprintf("{%s} \n%s \n%s", sr.State, sr.DriverInfo, sr.Message)
	}
	return fmt.Sprintf("{%s} %s", sr.State, sr.Message)
}

type ODBCError struct {
	StatusRecords []StatusRecord
}

func (e *ODBCError) Error() string {
	statusStrings := make([]string, len(e.StatusRecords))
	for i, sr := range e.StatusRecords {
		statusStrings[i] = sr.toString()
	}

	return strings.Join(statusStrings, "\n")
}

func ErrorEnvironment(handle syscall.Handle) error {
	return handleError(odbc.SQL_HANDLE_ENV, handle, "")
}

func ErrorConnection(handle syscall.Handle) error {
	return handleError(odbc.SQL_HANDLE_DBC, handle, "")
}

func ErrorStatement(handle syscall.Handle, driverInfo string) error {
	return handleError(odbc.SQL_HANDLE_STMT, handle, driverInfo)
}

func handleError(handleType odbc.SQLHandle, handle syscall.Handle, driverInfo string) error {
	statusRecords := make([]StatusRecord, 0)
	if handle != 0 {		
		for recNum := 1; ; recNum++ {
			sqlState := make([]uint16, SQLStateLength+1)
			var nativeError int
			message := make([]uint16, ErrorMaxMessageLength + 1)
			ret := odbc.SQLGetDiagRec(handleType, handle, int16(recNum), uintptr(unsafe.Pointer(&sqlState[0])), &nativeError, uintptr(unsafe.Pointer(&message[0])), ErrorMaxMessageLength, nil)
			if ret == odbc.SQL_NO_DATA {
				break
			} else if !IsError(ret) {
				sr := StatusRecord{State: syscall.UTF16ToString(sqlState), NativeError: nativeError, Message: syscall.UTF16ToString(message), DriverInfo: driverInfo}
				statusRecords = append(statusRecords, sr)
			} else {
				break
			}
		}
	}
	
	return &ODBCError{StatusRecords: statusRecords}
}