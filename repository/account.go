package repository

import (
	"context"

	"table-app/domain"
	"table-app/internal/db"

	"github.com/pkg/errors"
)

type Account struct {
	db db.DB
}

func NewAccount(db db.DB) Account {
	return Account{
		db: db,
	}
}

func (r Account) UpsertAll(ctx context.Context, accounts []domain.Account) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.WithMessage(err, "begin upsert transaction")
	}

	for _, account := range accounts {
		if account.Deleted {
			err = deleteAccount(ctx, tx.ExecContext, account.Id)
			if err != nil {
				rollBackErr := tx.Rollback()
				if rollBackErr != nil {
					return errors.WithMessage(err, "rollback upsert transaction")
				}

				return errors.WithMessage(err, "upsert account transaction")
			}

			continue
		}

		err = upsertAccount(ctx, tx.ExecContext, account)
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return errors.WithMessage(err, "rollback upsert transaction")
			}

			return errors.WithMessage(err, "upsert account transaction")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit upsert transaction")
	}

	return nil
}

func upsertAccount(ctx context.Context, txExec TxFuncExec, acc domain.Account) error {
	q := `
	INSERT INTO account
    	(id, name, sum, note, is_in_sum, updated_at)
	VALUES
    	($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE 
	SET name = $2, sum = $3, note = $4, is_in_sum = $5, updated_at = $6;`

	_, err := txExec(ctx, q, acc.Id, acc.Name, acc.Sum, acc.Note, acc.IsInSum, acc.UpdatedAt)
	if err != nil {
		return errors.WithMessage(err, "upsert account")
	}

	return nil
}

func deleteAccount(ctx context.Context, txExec TxFuncExec, accId string) error {
	q := `
	DELETE FROM account
	WHERE id = $1;`

	_, err := txExec(ctx, q, accId)
	if err != nil {
		return errors.WithMessage(err, "delete account")
	}

	return nil
}

func (r Account) GetAll(ctx context.Context) ([]domain.Account, error) {
	q := `
	SELECT id, name, sum, note, is_in_sum, updated_at FROM account;`

	var accounts []domain.Account
	rows, err := r.db.Select(ctx, q)
	if err != nil {
		return nil, errors.WithMessage(err, "get accounts")
	}

	for rows.Next() {
		var acc domain.Account
		err = rows.Scan(&acc.Id, &acc.Name, &acc.Sum, &acc.Note, &acc.IsInSum, &acc.UpdatedAt)
		if err != nil {
			return nil, errors.WithMessage(err, "scan row")
		}

		accounts = append(accounts, acc)
	}

	err = rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "close rows")
	}

	return accounts, nil
}
