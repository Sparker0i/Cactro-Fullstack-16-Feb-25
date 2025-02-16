package postgres

import (
	"context"
	"fmt"

	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type transactionManager struct {
	db *pgxpool.Pool
}

func NewTransactionManager(db *pgxpool.Pool) repository.TransactionManager {
	return &transactionManager{db: db}
}

type transaction struct {
	tx pgx.Tx
}

func (tm *transactionManager) Begin(ctx context.Context) (repository.Transaction, error) {
	tx, err := tm.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &transaction{tx: tx}, nil
}

func (t *transaction) Commit() error {
	if err := t.tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (t *transaction) Rollback() error {
	if err := t.tx.Rollback(context.Background()); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
