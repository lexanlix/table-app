package repository

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"table-app/conf"
	"table-app/domain"
	"table-app/internal/db"

	"github.com/pkg/errors"
)

type Category struct {
	db       db.DB
	filePath string
}

func NewCategory(db db.DB, storage conf.Storage) Category {
	var filePath string

	if storage.Files != nil {
		filePath = storage.Files.CategoryFilePath
	}

	return Category{
		db:       db,
		filePath: filePath,
	}
}

func (r Category) UpsertAll(ctx context.Context, categories []domain.Category) error {
	if len(r.filePath) != 0 {
		return r.writeToFile(categories)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.WithMessage(err, "begin upsert category transaction")
	}

	for _, category := range categories {
		err = upsertCategory(ctx, tx.ExecContext, category)
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return errors.WithMessage(err, "rollback upsert category transaction")
			}

			return errors.WithMessage(err, "upsert category transaction")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit upsert category transaction")
	}

	return nil
}

func upsertCategory(ctx context.Context, txExec TxFuncExec, category domain.Category) error {
	q := `
	INSERT INTO category
    	(id, name, main_category, priority)
	VALUES
    	($1, $2, $3, $4)
	ON CONFLICT (main_category, priority) 
	DO UPDATE SET name = $2;`

	_, err := txExec(ctx, q, category.Id, category.Name, category.MainCategory, category.Priority)
	if err != nil {
		return errors.WithMessage(err, "upsert category")
	}

	return nil
}

func (r Category) GetAll(ctx context.Context) ([]domain.Category, error) {
	if len(r.filePath) != 0 {
		return r.readFromFile()
	}

	q := `
	SELECT id, name, main_category, priority FROM category;`

	var list []domain.Category
	rows, err := r.db.Select(ctx, q)
	if err != nil {
		return nil, errors.WithMessage(err, "get categories")
	}

	for rows.Next() {
		var cat domain.Category
		err = rows.Scan(&cat.Id, &cat.Name, &cat.MainCategory, &cat.Priority)
		if err != nil {
			return nil, errors.WithMessage(err, "scan row")
		}
		list = append(list, cat)
	}

	err = rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "close rows")
	}

	return list, nil
}

func (r Category) readFromFile() ([]domain.Category, error) {
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

	result := make([]domain.Category, 0)
	for _, record := range records {
		category := domain.Category{}
		category.Id = record[0]
		category.Name = record[1]
		category.MainCategory = record[2]

		priority, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, errors.WithMessage(err, "convert priority value")
		}
		category.Priority = priority

		result = append(result, category)
	}

	return result, nil
}

func (r Category) writeToFile(data []domain.Category) error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY, 0664)
	if err != nil {
		return errors.WithMessage(err, "open file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, category := range data {
		err := writer.Write([]string{
			category.Id,
			category.Name,
			category.MainCategory,
			strconv.Itoa(category.Priority),
		})
		if err != nil {
			return errors.WithMessage(err, "write to file")
		}
	}
	writer.Flush()

	return nil
}
