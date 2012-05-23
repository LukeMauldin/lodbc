package lodbc

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"github.com/LukeMauldin/lodbc/odbc"
	"time"
)

// Registers the time.Time struct with the gob package
func init() {
	tp := new(time.Time)
	gob.Register(tp)
}

// Struct to hold additional metadata for bind parameters for use in 
// stmt.Query and stmt.Exec.  ODBC driver benefits and in some cases requires
// additional fields in order to correctly bind the parameters.
type BindParameter struct {
	
	// Contains the bind parameter value
	Data      driver.Value 
	
	// Valid for strings only.  Specifies the maximum length of the string.
	// If 0, defaults to the length of the string in Data
	Length    int 
	
	// Valid for float64 only
	Precision int
	
	// Valid for float64 only
	Scale     int
	
	// Valid for time.Time only
	DateOnly  bool
	
	// Specifies the direction of the ODBC parameter.  Defaults to InputParameter
	Direction ParameterDirection
}

/*
 * Implement the Valuer interface to convert a BindParameter to a driver.Value
 * Uses GOB encoding to encode as a []byte to bypass the restriction on driver.Value types
 */
func (bp *BindParameter) Value() (driver.Value, error) {
	//Return nil if bp.Data is nil
	if isNil(bp.Data) {
		return nil, nil
	}
	//GOB encode the bind parameter
	encodedBuffer := new(bytes.Buffer)
	enc := gob.NewEncoder(encodedBuffer)
	err := enc.Encode(bp)
	if err != nil {
		return nil, err
	}
	return encodedBuffer.Bytes(), nil
}

// Indicates direction of ODBC parameter. Maps to an ODBC parameter direction.
// Currently the only supported value is input parameter
type ParameterDirection int
const (
	InputParameter       ParameterDirection = 1 << iota
	// OutputParameter - not suported
	//InputOutputParameter - not suported
)

/*
 * Converts ParameterDirection to an ODBC parameter direction
 * Currently the only supported value is SQL_PARAM_INPUT
 */
func (p ParameterDirection) SQLBindParameterType() odbc.SQLBindParameterType {
	return odbc.SQL_PARAM_INPUT
	/*
	switch p {
	case InputParameter:
		return odbc.SQL_PARAM_INPUT
	case OutputParameter:
		return odbc.SQL_PARAM_OUTPUT
	case InputOutputParameter:
		return odbc.SQL_PARAM_INPUT_OUTPUT
	}
	panic("Parameter direction: " + strconv.Itoa(int(p)))
	*/
}

// Create a new bind parameter for an int
func NewParameterInt(data driver.Value) *BindParameter {
	return &BindParameter{Data: data}
}

// Create a new bind parameter for an int64
func NewParameterInt64(data driver.Value) *BindParameter {
	return &BindParameter{Data: data}
}

// Create a new bind parameter for a float
func NewParameterFloat(data driver.Value, precision int, scale int) *BindParameter {
	return &BindParameter{Data: data, Precision: precision, Scale: scale}
}

// Create a new bind parameter for a date
func NewParameterDate(data driver.Value) *BindParameter {
	return &BindParameter{Data: data, DateOnly: true}
}

// Create a new bind parameter for a date time
func NewParameterDateTime(data driver.Value) *BindParameter {
	return &BindParameter{Data: data, DateOnly: false}
}

// Create a new bind parameter for a string
func NewParameterString(data driver.Value, length int) *BindParameter {
	return &BindParameter{Data: data, Length: length}
}