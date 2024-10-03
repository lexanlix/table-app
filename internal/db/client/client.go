package db

import (
	"context"
	"database/sql"
	"fmt"

	"table-app/internal/log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type MigrationRunner interface {
	Run(ctx context.Context, db *sql.DB, gooseOpts ...goose.ProviderOption) error
}

type Client struct {
	logger log.Logger
	cli    *sql.DB

	migrationRunner MigrationRunner
}

func NewClient(logger log.Logger, opts ...Option) *Client {
	client := &Client{
		logger: logger,
		cli:    &sql.DB{},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Upgrade(ctx context.Context, config StorageConfig) error {
	dsn, err := config.getDsn()
	if err != nil {
		return errors.WithMessage(err, "get dsn from cfg")
	}

	if config.Schema != "public" && config.Schema != "" {
		err = c.createSchema(ctx, config, dsn)
		if err != nil {
			return errors.WithMessage(err, "create schema")
		}
	}

	cli, err := c.Open(dsn)
	if err != nil {
		return errors.WithMessage(err, "open db client")
	}

	if c.migrationRunner != nil {
		err = c.migrationRunner.Run(ctx, cli)
		if err != nil {
			return errors.WithMessage(err, "run migration")
		}
	}

	c.cli = cli
	return nil
}

func (c *Client) Open(dsn string) (*sql.DB, error) {
	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "parse config")
	}

	return stdlib.OpenDB(*pgxConfig), nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.cli.ExecContext(ctx, query, args...)
}

func (c *Client) Select(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.cli.QueryContext(ctx, query, args...)
}

func (c *Client) SelectRow(ctx context.Context, query string, args ...any) *sql.Row {
	return c.cli.QueryRowContext(ctx, query, args...)
}

func (c *Client) Begin() (*sql.Tx, error) {
	return c.cli.Begin()
}

func (c *Client) BeginTx(ctx context.Context, txOptions *sql.TxOptions) (*sql.Tx, error) {
	return c.cli.BeginTx(ctx, txOptions)
}

func (c *Client) createSchema(ctx context.Context, config StorageConfig, dsn string) error {
	schema := config.Schema

	config.Schema = ""
	dbCli, err := c.Open(dsn)
	if err != nil {
		return errors.WithMessage(err, "open db")
	}

	_, err = dbCli.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	if err != nil {
		return errors.WithMessage(err, "exec query")
	}

	err = dbCli.Close()
	if err != nil {
		return errors.WithMessage(err, "close db")
	}

	return nil
}
