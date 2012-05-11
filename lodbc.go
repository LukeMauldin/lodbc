package lodbc

import (
	"github.com/lukemauldin/lodbc/odbc"
	"strconv"
)

type ParameterDirection int16

const (
	InputParameter ParameterDirection = 0
	OutputParameter ParameterDirection = 1
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

type ResultColumnDef struct {
	RecNum    uint16
	DataType  odbc.SQLDataType
	Length    int32
	Precision int32
	Scale     int32
}


