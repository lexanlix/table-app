package db

import (
	"context"
	"database/sql"
)

type DB interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Select(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	SelectRow(ctx context.Context, query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, txOptions *sql.TxOptions) (*sql.Tx, error)
}

type Transactional interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, txOptions *sql.TxOptions) (*sql.Tx, error)
}
