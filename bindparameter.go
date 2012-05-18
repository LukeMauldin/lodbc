package lodbc

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"github.com/LukeMauldin/lodbc/odbc"
	"strconv"
	"time"
)

type BindParameter struct {
	Data      driver.Value
	Length    int
	Precision int
	Scale     int
	DateOnly  bool
	Direction ParameterDirection
}

func init() {
	/* var bp BindParameter
	gob.Register(bp)
	var t time.Time
	gob.Register(t) */
	//var t time.Time
	//gob.Register(t)
	tp := new(time.Time)
	gob.Register(tp)
}

func (bp BindParameter) Value() (driver.Value, error) {
	//No conversion necessary if bp.Data is nil
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

type ParameterDirection int16

const (
	InputParameter       ParameterDirection = 0
	OutputParameter      ParameterDirection = 1
	InputOutputParameter ParameterDirection = 2
)

func (p ParameterDirection) SQLBindParameterType() odbc.SQLBindParameterType {
	switch p {
	case InputParameter:
		return odbc.SQL_PARAM_INPUT
	case OutputParameter:
		return odbc.SQL_PARAM_OUTPUT
	case InputOutputParameter:
		return odbc.SQL_PARAM_INPUT_OUTPUT
	}
	panic("Parameter direction: " + strconv.Itoa(int(p)))
}
