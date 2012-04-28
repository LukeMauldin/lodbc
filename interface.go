package lodbc

import (
	_ "time"
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