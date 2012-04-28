package odbc

import (
	"syscall"
)

//SQL Handle const
type SQLHandle syscall.Handle

//SQL general types
type SQLCHAR uint8
type SQLSCHAR int8
type SQLSMALLINT int16
type SQLUSMALLINT uint16
type SQLINTEGER int32
type SQLPOINTER syscall.Handle
type SQLUINTEGER uint32
type SQLLEN int64
type SQLULEN uint64

const (
	SQL_HANDLE_ENV  = 1
	SQL_HANDLE_DBC  = 2
	SQL_HANDLE_STMT = 3
	SQL_HANDLE_DESC = 4
)

//SQL Return codes
type SQLReturn int16

const (
	SQL_SUCCESS              = 0
	SQL_SUCCESS_WITH_INFO    = 1
	SQL_NO_DATA              = 100
	SQL_PARAM_DATA_AVAILABLE = 101
	SQL_ERROR                = -1
	SQL_INVALID_HANDLE       = -2
	SQL_STILL_EXECUTING      = 2
	SQL_NEED_DATA            = 99
)

//ODBC Version
const (
	SQL_OV_ODBC3    = 3
	SQL_OV_ODBC3_80 = 380
)

//Env Attributes
const (
	SQL_ATTR_ODBC_VERSION       = 200
	SQL_ATTR_CONNECTION_POOLING = 201
	SQL_ATTR_CP_MATCH           = 202
)

//Options for SQLDriverConnect 
const (
	SQL_DRIVER_NOPROMPT = 0
)

//Options for SQLFetch
const (
	SQL_FETCH_NEXT = 1
)

//SQL data types
type SQLDataType SQLSMALLINT 

const (
	SQL_UNKNOWN_TYPE  = 0
	SQL_CHAR          = 1
	SQL_NUMERIC       = 2
	SQL_DECIMAL       = 3
	SQL_INTEGER       = 4
	SQL_SMALLINT      = 5
	SQL_FLOAT         = 6
	SQL_REAL          = 7
	SQL_DOUBLE        = 8
	SQL_DATETIME      = 9
	SQL_VARCHAR       = 12
	SQL_DATE          = 91
	SQL_TIME          = 92
	SQL_TIMESTAMP     = 11
	SQL_LONGVARCHAR   = -1
	SQL_BINARY        = -2
	SQL_VARBINARY     = -3
	SQL_LONGVARBINARY = -4
	SQL_BIGINT        = -5
	SQL_TINYINT       = -6
	SQL_BIT           = -7
	SQL_WCHAR         = -8
	SQL_WVARCHAR      = -9
)

//C data types
type CDataType SQLSMALLINT

const (
	SQL_C_CHAR      = SQL_CHAR
	SQL_C_LONG      = SQL_INTEGER
	SQL_C_SHORT     = SQL_SMALLINT
	SQL_C_FLOAT     = SQL_REAL
	SQL_C_DOUBLE    = SQL_DOUBLE
	SQL_C_NUMERIC   = SQL_NUMERIC
	SQL_C_DATE      = SQL_DATE
	SQL_C_TIME      = SQL_TIME
	SQL_C_TIMESTAMP = SQL_TIMESTAMP
	SQL_C_BINARY    = SQL_BINARY
	SQL_C_BIT       = SQL_BIT
	SQL_C_WCHAR     = SQL_WCHAR
	SQL_C_DEFAULT   = 99
)

//Misc flags
const (
	SQL_NTS = -3
)

//SQL Transaction options
type SQLTransactionOption int16

const (
	SQL_COMMIT   = 0
	SQL_ROLLBACK = 1
)

//SQLBindParameter options
type SQLBindParameterType SQLSMALLINT

const (
	SQL_PARAM_TYPE_UNKNOWN = 0
	SQL_PARAM_INPUT        = 1
	SQL_PARAM_INPUT_OUTPUT = 2
	SQL_RESULT_COL         = 3
	SQL_PARAM_OUTPUT       = 4
	SQL_RETURN_VALUE       = 5
)

//SQL descriptors
type SQLDescriptor uint16

const (
	SQL_DESC_COUNT                  = 1001
	SQL_DESC_TYPE                   = 1002
	SQL_DESC_LENGTH                 = 1003
	SQL_DESC_OCTET_LENGTH_PTR       = 1004
	SQL_DESC_PRECISION              = 1005
	SQL_DESC_SCALE                  = 1006
	SQL_DESC_DATETIME_INTERVAL_CODE = 1007
	SQL_DESC_NULLABLE               = 1008
	SQL_DESC_INDICATOR_PTR          = 1009
	SQL_DESC_DATA_PTR               = 1010
	SQL_DESC_NAME                   = 1011
	SQL_DESC_UNNAMED                = 1012
	SQL_DESC_OCTET_LENGTH           = 1013
	SQL_DESC_ALLOC_TYPE             = 1099
)

//Special length/indicator values
type SQLValueIndicator SQLLEN

const (
	SQL_NULL_DATA    = -1
	SQL_DATA_AT_EXEC = -2
)

type SQL_NUMERIC_STRUCT struct {
	Precision SQLCHAR
	Scale     SQLCHAR
	Sign      SQLSCHAR
	Val       [16]SQLCHAR
}

type SQL_DATE_STRUCT struct {
	Year  SQLSMALLINT
	Month SQLUSMALLINT
	Day   SQLUSMALLINT
}

type SQL_TIME_STRUCT struct {
	Hour   SQLUSMALLINT
	Minute SQLUSMALLINT
	Second SQLUSMALLINT
}

type SQL_TIMESTAMP_STRUCT struct {
	Year    SQLSMALLINT
	Month   SQLUSMALLINT
	Day     SQLUSMALLINT
	Hour    SQLUSMALLINT
	Minute  SQLUSMALLINT
	Second  SQLUSMALLINT
	Faction SQLUINTEGER
}

//Connection attributes
type SQLConnectionAttribute SQLINTEGER

const (
	SQL_ATTR_AUTOCOMMIT    = 102
	SQL_AUTOCOMMIT_OFF     = 0
	SQL_AUTOCOMMIT_ON      = 1
	SQL_AUTOCOMMIT_DEFAULT = SQL_AUTOCOMMIT_ON
)

//Statement attributes
type SQLStatementAttribute SQLINTEGER

const (
	SQL_QUERY_TIMEOUT           = 0
	SQL_MAX_ROWS                = 1
	SQL_NOSCAN                  = 2
	SQL_ATTR_QUERY_TIMEOUT      = SQL_QUERY_TIMEOUT
	SQL_ATTR_APP_ROW_DESC       = 10010
	SQL_ATTR_APP_PARAM_DESC     = 10011
	SQL_ATTR_IMP_ROW_DESC       = 10012
	SQL_ATTR_IMP_PARAM_DESC     = 10013
	SQL_ATTR_CURSOR_SCROLLABLE  = -1
	SQL_ATTR_CURSOR_SENSITIVITY = -2
)

//Code indicating that the application row descriptor specifies the data type
const (
	SQL_ARD_TYPE = -99
)

//Values for SQL_ATTR_CONNECTION_POOLING
const (
	SQL_CP_OFF            = 0
	SQL_CP_ONE_PER_DRIVER = 1
	SQL_CP_ONE_PER_HENV   = 2
	SQL_CP_DEFAULT        = SQL_CP_OFF
)

//Whether an attribute is a pointer or not
const (
	SQL_IS_POINTER   = -4
	SQL_IS_UINTEGER  = -5
	SQL_IS_INTEGER   = -6
	SQL_IS_USMALLINT = -7
	SQL_IS_SMALLINT  = -8
)
