package repository

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"table-app/conf"
	"table-app/domain"
	"table-app/internal/db"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const batchSize = 1000

type TxFuncExec func(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)

type Table struct {
	db       db.DB
	filePath string
}

func NewTable(db db.DB, storage conf.Storage) Table {
	var filePath string

	if storage.Files != nil {
		filePath = storage.Files.TableFilePath
	}

	return Table{
		db:       db,
		filePath: filePath,
	}
}

func (r Table) UpsertAll(ctx context.Context, cells []domain.Cell) error {
	if len(r.filePath) != 0 {
		return r.writeToFile(cells)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.WithMessage(err, "begin upsert transaction")
	}

	for _, cell := range cells {
		err = upsertCell(ctx, tx.Exec, cell)
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				return errors.WithMessage(err, "rollback upsert transaction")
			}

			return errors.WithMessage(err, "upsert cell transaction")
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
	if len(r.filePath) != 0 {
		return r.readFromFile()
	}

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

func (r Table) readFromFile() ([]domain.Cell, error) {
	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return nil, errors.WithMessage(err, "open file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.WithMessage(err, "read all file")
	}

	result := make([]domain.Cell, 0)
	for _, record := range records {
		cell := domain.Cell{}
		cell.Id = record[0]
		cell.MainCategory = record[1]
		cell.Category = record[2]

		value, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, errors.WithMessage(err, "convert cell value")
		}
		cell.Value = value

		month, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, errors.WithMessage(err, "convert month value")
		}
		cell.Month = time.Month(month)

		year, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, errors.WithMessage(err, "convert year value")
		}
		cell.Year = year

		result = append(result, cell)
	}

	return result, nil
}

func (r Table) writeToFile(data []domain.Cell) error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY, 0664)
	if err != nil {
		return errors.WithMessage(err, "open file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, cell := range data {
		err := writer.Write([]string{
			cell.Id,
			cell.MainCategory,
			cell.Category,
			strconv.Itoa(cell.Value),
			strconv.Itoa(int(cell.Month)),
			strconv.Itoa(cell.Year),
		})
		if err != nil {
			return errors.WithMessage(err, "write to file")
		}
	}
	writer.Flush()

	return nil
}
