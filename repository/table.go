package repository

import (
	"context"

	"table-app/domain"
	"table-app/internal/db"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const batchSize = 1000

type TxFuncExec func(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)

type Table struct {
	db db.DB
}

func NewTable(db db.DB) Table {
	return Table{
		db: db,
	}
}

func (r Table) UpsertAll(ctx context.Context, cells []domain.Cell) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.WithMessage(err, "begin upsert transaction")
	}

	for _, cell := range cells {
		err = upsertCell(ctx, tx.Exec, cell)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return errors.WithMessage(err, "rollback upsert transaction")
			}

			return errors.WithMessage(err, "upsert cell")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.WithMessage(err, "commit upsert transaction")
	}

	return nil
}

func upsertCell(ctx context.Context, txExec TxFuncExec, cell domain.Cell) error {
	q := `
	INSERT INTO table_app.finances
    	(id, main_category, category, value, month, year)
	VALUES
    	($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE 
	    SET value = $4;`

	_, err := txExec(ctx, q, cell.Id, cell.MainCategory, cell.Category, cell.Value, cell.Month, cell.Year)
	if err != nil {
		return errors.WithMessage(err, "upsert cell")
	}

	return nil
}

func (r Table) GetAll(ctx context.Context) ([]domain.Cell, error) {
	q := `
	SELECT id, main_category, category, value, month, year 
	FROM table_app.finances;`

	var cells []domain.Cell
	rows, err := r.db.Select(ctx, q)
	if err != nil {
		return nil, errors.WithMessage(err, "get cells")
	}

	defer rows.Close()
	for rows.Next() {
		var cell domain.Cell
		err = rows.Scan(&cell.Id, &cell.MainCategory, &cell.Category, &cell.Value, &cell.Month, &cell.Year)
		if err != nil {
			return nil, errors.WithMessage(err, "scan row")
		}
		cells = append(cells, cell)
	}

	return cells, nil
}
