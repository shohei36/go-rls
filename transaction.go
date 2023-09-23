package main

import (
	"context"
	"database/sql"
	"fmt"
)

type Transaction interface {
	DoInTx(ctx context.Context, tenantID string, f func(context.Context) (any, error)) (any, error)
}

type ctxTxKey struct{}

type transaction struct {
	tdb TenantDB
}

func NewTransaction(tdb TenantDB) Transaction {
	return &transaction{
		tdb: tdb,
	}
}

func (t *transaction) DoInTx(
	ctx context.Context, tenantID string, f func(context.Context) (any, error),
) (any, error) {
	conn, err := t.tdb.Conn(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, ctxTxKey{}, tx)

	v, err := f(ctx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("rollback failed: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, err
		}
		return nil, err
	}
	return v, nil
}

func GetTx(ctx context.Context) (DBTX, bool) {
	tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx)
	return tx, ok
}

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}
