package lodbc

import (
	"time"
)

type IConnection interface {
	Close() (error)
	Begin() (error)
	Commit() (error)
	Rollback() (error)
	IsTransactionActive() (bool)
	NewStatement() (IStatement, error)
}

type IStatement interface {
	BindInt(index int, value *int, direction ParameterDirection) (error)
	BindInt64(index int, value *int64, direction ParameterDirection) (error)
	BindBool(index int, value *bool, direction ParameterDirection) (error)
	BindNumeric(index int, value *float64, precision int, scale int, direction ParameterDirection) (error)
	BindDate(index int, value *time.Time, direction ParameterDirection) (error)
	BindDateTime(index int, value *time.Time, direction ParameterDirection) (error)
	BindString(index int, value *string, length int, direction ParameterDirection) (error)
	BindNull(index int, direction ParameterDirection) error
	Query(query string) (IRows, error) 
	QueryWithParams(query string, parameters ...BindParameter) (IRows, error)
	Exec(query string) (error)
	ExecWithParams(query string, parameters ...BindParameter) (error)
	Close() (error)
}

type IRows interface {
	Next() (bool)
	Close() (error)
	Scan(dest ...interface{}) (error)
	Err() (error)
}

type BindParameter struct {
	Value interface{}
	Length int
	Precision int
	Scale int
	DateOnly bool
	Direction ParameterDirection
}