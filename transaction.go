package lodbc

import (
	"github.com/LukeMauldin/lodbc/odbc"
)

// Implements type database/sql/driver TX interface
type transaction struct {
	conn *connection
}

// Commit transaction
func (tx *transaction) Commit() error {
	return tx.completeTransaction(odbc.SQL_COMMIT)
}

// Rollback transaction
func (tx *transaction) Rollback() error {
	return tx.completeTransaction(odbc.SQL_ROLLBACK)
}

// Commit or rollback transaction in consistent manner
func (tx *transaction) completeTransaction(completeType odbc.SQLTransactionOption) error {
	//Complete transaction by either committing or rolling back
	ret := odbc.SQLEndTran(odbc.SQL_HANDLE_DBC, tx.conn.handle, completeType)
	if isError(ret) {
		return errorConnection(tx.conn.handle)
	}

	//Make transaction as finished and turn auto commit back on
	tx.conn.isTransactionActive = false
	ret = odbc.SQLSetConnectAttr(tx.conn.handle, odbc.SQL_ATTR_AUTOCOMMIT, odbc.SQL_AUTOCOMMIT_ON, 0, nil)
	if isError(ret) {
		return errorConnection(tx.conn.handle)
	}
	return nil
}
