package storage

import (
	"context"
	"database/sql"
	"log"
)

// abstracts db engine (*sql.DB already implements this interface)
// to allow mocking in unit tests without relying on real db
type TransactionMgr interface {
	BeginTransaction(ctx context.Context, readOnly bool) (SQLExecutor, func(error), error)
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
) (SQLExecutor, func(opErr error), error) {

	transactionExec, err := transactionMgr.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: readOnly})
	if err != nil {
		return nil, nil, err
	}

	// deferred function to either commit or rollback transaction
	finish := func(opErr error) {
		if opErr != nil {
			log.Printf("Rolling back transaction: %v", opErr)
			if err := transactionExec.Rollback(); err != nil {
				log.Printf("transaction rollback failed: %v", err)
			}
			return
		}
		if err := transactionExec.Commit(); err != nil {
			log.Printf("transaction commit failed: %v", err)
		}
	}

	return transactionExec, finish, nil
}
