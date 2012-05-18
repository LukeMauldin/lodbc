package lodbc

import (
	"github.com/LukeMauldin/lodbc/odbc"
)

type ResultColumnDef struct {
	RecNum    uint16
	DataType  odbc.SQLDataType
	Length    int32
	Precision int32
	Scale     int32
	Name      string
}
