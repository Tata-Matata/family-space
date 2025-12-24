package storage

import (
	"context"
	"database/sql"
)

// abstracts db engine (*sql.DB already implements this interface)
// to allow mocking in unit tests without relying on real db
type TransactionMgr interface {
	BeginTransaction(ctx context.Context, readOnly bool) (SQLExecutor, func(error) error, error)
}

// concrete implementation of TransactionManager for sql.DB
// wraps *sql.DB for production use, not testing
type sqlTransactionMgr struct {
	db *sql.DB
}

// starts transactional execution and returns sql.Tx as SQLExecutor
func (transactionMgr *sqlTransactionMgr) BeginTransaction(
	ctx context.Context,
	readOnly bool,
) (SQLExecutor, func(opErr error) error, error) {

	transactionExec, err := transactionMgr.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: readOnly})
	if err != nil {
		return nil, nil, err
	}

	// deferred function to either commit or rollback transaction
	finish := func(opErr error) error {
		if opErr != nil {
			_ = transactionExec.Rollback()
			return opErr
		}
		return transactionExec.Commit()
	}

	return transactionExec, finish, nil
}
