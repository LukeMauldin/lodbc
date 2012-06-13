package odbc

//SQL Handle
type SQLHandle uintptr

//SQL general types
type SQLCHAR uint8
type SQLSCHAR int8
type SQLSMALLINT int16
type SQLUSMALLINT uint16
type SQLINTEGER int32
type SQLPOINTER uintptr
type SQLUINTEGER uint32
type SQLLEN int64
type SQLULEN uint64

const (
	SQL_HANDLE_ENV  SQLSMALLINT = 1
	SQL_HANDLE_DBC  SQLSMALLINT = 2
	SQL_HANDLE_STMT SQLSMALLINT = 3
	SQL_HANDLE_DESC SQLSMALLINT = 4
)

//SQL Return codes
type SQLReturn SQLSMALLINT

const (
	SQL_SUCCESS              SQLReturn = 0
	SQL_SUCCESS_WITH_INFO    SQLReturn = 1
	SQL_NO_DATA              SQLReturn = 100
	SQL_PARAM_DATA_AVAILABLE SQLReturn = 101
	SQL_ERROR                SQLReturn = -1
	SQL_INVALID_HANDLE       SQLReturn = -2
	SQL_STILL_EXECUTING      SQLReturn = 2
	SQL_NEED_DATA            SQLReturn = 99
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
	SQL_UNKNOWN_TYPE   SQLDataType = 0
	SQL_CHAR           SQLDataType = 1
	SQL_NUMERIC        SQLDataType = 2
	SQL_DECIMAL        SQLDataType = 3
	SQL_INTEGER        SQLDataType = 4
	SQL_SMALLINT       SQLDataType = 5
	SQL_FLOAT          SQLDataType = 6
	SQL_REAL           SQLDataType = 7
	SQL_DOUBLE         SQLDataType = 8
	SQL_DATE           SQLDataType = 9
	SQL_TIME           SQLDataType = 10
	SQL_VARCHAR        SQLDataType = 12
	SQL_TYPE_DATE      SQLDataType = 91
	SQL_TYPE_TIME      SQLDataType = 92
	SQL_TYPE_TIMESTAMP SQLDataType = 93
	SQL_TIMESTAMP      SQLDataType = 11
	SQL_LONGVARCHAR    SQLDataType = -1
	SQL_BINARY         SQLDataType = -2
	SQL_VARBINARY      SQLDataType = -3
	SQL_LONGVARBINARY  SQLDataType = -4
	SQL_BIGINT         SQLDataType = -5
	SQL_TINYINT        SQLDataType = -6
	SQL_BIT            SQLDataType = -7
	SQL_WCHAR          SQLDataType = -8
	SQL_WVARCHAR       SQLDataType = -9
	SQL_SS_XML         SQLDataType = -152
)

//C data types
type CDataType SQLSMALLINT

const (
	SQL_C_CHAR      CDataType = CDataType(SQL_CHAR)
	SQL_C_LONG      CDataType = CDataType(SQL_INTEGER)
	SQL_C_SHORT     CDataType = CDataType(SQL_SMALLINT)
	SQL_C_FLOAT     CDataType = CDataType(SQL_REAL)
	SQL_C_DOUBLE    CDataType = CDataType(SQL_DOUBLE)
	SQL_C_NUMERIC   CDataType = CDataType(SQL_NUMERIC)
	SQL_C_DATE      CDataType = CDataType(SQL_DATE)
	SQL_C_TIME      CDataType = CDataType(SQL_TIME)
	SQL_C_TIMESTAMP CDataType = CDataType(SQL_TIMESTAMP)
	SQL_C_BINARY    CDataType = CDataType(SQL_BINARY)
	SQL_C_BIT       CDataType = CDataType(SQL_BIT)
	SQL_C_WCHAR     CDataType = CDataType(SQL_WCHAR)
	SQL_C_DEFAULT   CDataType = CDataType(99)
)

//Misc flags
const (
	SQL_NTS = -3
)

//SQL Transaction options
const (
	SQL_COMMIT   SQLSMALLINT = 0
	SQL_ROLLBACK SQLSMALLINT = 1
)

//SQLBindParameter options
const (
	SQL_PARAM_TYPE_UNKNOWN SQLSMALLINT = 0
	SQL_PARAM_INPUT        SQLSMALLINT = 1
	SQL_PARAM_INPUT_OUTPUT SQLSMALLINT = 2
	SQL_RESULT_COL         SQLSMALLINT = 3
	SQL_PARAM_OUTPUT       SQLSMALLINT = 4
	SQL_RETURN_VALUE       SQLSMALLINT = 5
)

//SQL descriptors
type SQLDescriptor SQLSMALLINT

const (
	SQL_DESC_COUNT                  SQLSMALLINT = 1001
	SQL_DESC_TYPE                   SQLSMALLINT = 1002
	SQL_DESC_LENGTH                 SQLSMALLINT = 1003
	SQL_DESC_OCTET_LENGTH_PTR       SQLSMALLINT = 1004
	SQL_DESC_PRECISION              SQLSMALLINT = 1005
	SQL_DESC_SCALE                  SQLSMALLINT = 1006
	SQL_DESC_DATETIME_INTERVAL_CODE SQLSMALLINT = 1007
	SQL_DESC_NULLABLE               SQLSMALLINT = 1008
	SQL_DESC_INDICATOR_PTR          SQLSMALLINT = 1009
	SQL_DESC_DATA_PTR               SQLSMALLINT = 1010
	SQL_DESC_NAME                   SQLSMALLINT = 1011
	SQL_DESC_UNNAMED                SQLSMALLINT = 1012
	SQL_DESC_OCTET_LENGTH           SQLSMALLINT = 1013
	SQL_DESC_ALLOC_TYPE             SQLSMALLINT = 1099
)

//SQLColAttributes
type SQLColAttributeType SQLUSMALLINT

const (
	SQL_COLUMN_TYPE      SQLColAttributeType = 2
	SQL_COLUMN_LENGTH    SQLColAttributeType = 3
	SQL_COLUMN_PRECISION SQLColAttributeType = 4
	SQL_COLUMN_SCALE     SQLColAttributeType = 5
	SQL_COLUMN_NULLABLE  SQLColAttributeType = 7
	SQL_COLUMN_LABEL     SQLColAttributeType = 18
	SQL_DESC_LABEL       SQLColAttributeType = SQL_COLUMN_LABEL
)

//Special length/indicator values
const (
	SQL_NULL_DATA    SQLLEN = -1
	SQL_DATA_AT_EXEC SQLLEN = -2
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
const (
	SQL_ATTR_AUTOCOMMIT    SQLINTEGER = 102
	SQL_AUTOCOMMIT_OFF     SQLINTEGER = 0
	SQL_AUTOCOMMIT_ON      SQLINTEGER = 1
	SQL_AUTOCOMMIT_DEFAULT SQLINTEGER = SQL_AUTOCOMMIT_ON
)

//Statement attributes
const (
	SQL_QUERY_TIMEOUT           SQLINTEGER = 0
	SQL_MAX_ROWS                SQLINTEGER = 1
	SQL_NOSCAN                  SQLINTEGER = 2
	SQL_ATTR_QUERY_TIMEOUT      SQLINTEGER = SQL_QUERY_TIMEOUT
	SQL_ATTR_APP_ROW_DESC       SQLINTEGER = 10010
	SQL_ATTR_APP_PARAM_DESC     SQLINTEGER = 10011
	SQL_ATTR_IMP_ROW_DESC       SQLINTEGER = 10012
	SQL_ATTR_IMP_PARAM_DESC     SQLINTEGER = 10013
	SQL_ATTR_CURSOR_SCROLLABLE  SQLINTEGER = -1
	SQL_ATTR_CURSOR_SENSITIVITY SQLINTEGER = -2
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
