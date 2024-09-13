package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"table-app/internal/log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Database interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Client struct {
	logger      log.Logger
	cli         *pgxpool.Pool
	maxAttempts int
}

func NewClient(logger log.Logger) *Client {
	return &Client{
		logger:      logger,
		cli:         &pgxpool.Pool{},
		maxAttempts: 3,
	}
}

func (c *Client) Upgrade(ctx context.Context, cfg StorageConfig) error {
	dsn, err := getDsn(cfg)
	if err != nil {
		return errors.WithMessage(err, "get dsn from cfg")
	}

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		c.cli, err = pgxpool.New(ctx, dsn)
		if err != nil {
			c.logger.Error(ctx, "connect to postgres")
			return errors.WithMessage(err, "connect to postgres")
		}

		return nil
	}, c.maxAttempts, 5*time.Second)

	if err != nil {
		return errors.WithMessage(err, "can't connect to database")
	}

	return nil
}

func (c *Client) Close() error {
	c.cli.Close()
	return nil
}

func (c *Client) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return c.cli.Exec(ctx, query, args...)
}

func (c *Client) Select(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return c.cli.Query(ctx, query, args...)
}

func (c *Client) SelectRow(ctx context.Context, query string, args ...any) pgx.Row {
	return c.cli.QueryRow(ctx, query, args...)
}

func (c *Client) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.cli.Begin(ctx)
}

func (c *Client) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.cli.BeginTx(ctx, txOptions)
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}

	return
}

func getDsn(cfg StorageConfig) (string, error) {
	if len(cfg.Host) == 0 || len(cfg.Port) == 0 {
		return "", errors.New("invalid db configuration: host and port are required")
	}

	if len(cfg.Database) == 0 {
		return "", errors.New("invalid db configuration: database is required")
	}

	if len(cfg.Username) == 0 || len(cfg.Password) == 0 {
		return "", errors.New("invalid db configuration: username and password are required")
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database), nil
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}
