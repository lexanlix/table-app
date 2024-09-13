package repository

import (
	"context"

	"table-app/domain"
	"table-app/internal/db"

	"github.com/pkg/errors"
)

type Category struct {
	db db.DB
}

func NewCategory(db db.DB) Category {
	return Category{
		db: db,
	}
}

func (r Category) UpsertAll(ctx context.Context, categories []domain.Category) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.WithMessage(err, "begin upsert category transaction")
	}

	for _, category := range categories {
		err = upsertCategory(ctx, tx.Exec, category)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return errors.WithMessage(err, "rollback upsert category transaction")
			}

			return errors.WithMessage(err, "upsert cell")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.WithMessage(err, "commit upsert category transaction")
	}

	return nil
}

func upsertCategory(ctx context.Context, txExec TxFuncExec, category domain.Category) error {
	q := `
	INSERT INTO table_app.category
    	(name, main_category, priority)
	VALUES
    	($1, $2, $3)
	ON CONFLICT (main_category, priority) 
	DO UPDATE SET name = $1;`

	_, err := txExec(ctx, q, category.Name, category.MainCategory, category.Priority)
	if err != nil {
		return errors.WithMessage(err, "upsert category")
	}

	return nil
}

func (r Category) GetAll(ctx context.Context) ([]domain.Category, error) {
	q := `
	SELECT id, name, main_category, priority
	FROM table_app.category;`

	var list []domain.Category
	rows, err := r.db.Select(ctx, q)
	if err != nil {
		return nil, errors.WithMessage(err, "get categories")
	}

	defer rows.Close()
	for rows.Next() {
		var cat domain.Category
		err = rows.Scan(&cat.Id, &cat.Name, &cat.MainCategory, &cat.Priority)
		if err != nil {
			return nil, errors.WithMessage(err, "scan row")
		}
		list = append(list, cat)
	}

	return list, nil
}
