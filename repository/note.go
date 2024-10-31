package repository

import (
	"context"

	"table-app/domain"
	"table-app/internal/db"

	"github.com/pkg/errors"
)

type Note struct {
	db db.DB
}

func NewNote(db db.DB) Note {
	return Note{
		db: db,
	}
}

func (r Note) UpsertAll(ctx context.Context, notes []domain.Note) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.WithMessage(err, "begin upsert transaction")
	}

	for _, note := range notes {
		if note.Deleted {
			err = deleteNote(ctx, tx.ExecContext, note.Id)
			if err != nil {
				rollBackErr := tx.Rollback()
				if rollBackErr != nil {
					return errors.WithMessage(err, "rollback delete transaction")
				}

				return errors.WithMessage(err, "delete note transaction")
			}

			continue
		}

		err = upsertNote(ctx, tx.ExecContext, note)
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return errors.WithMessage(err, "rollback upsert transaction")
			}

			return errors.WithMessage(err, "upsert note transaction")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit upsert transaction")
	}

	return nil
}

func upsertNote(ctx context.Context, txExec TxFuncExec, note domain.Note) error {
	q := `
	INSERT INTO note
    	(id, text, updated_at)
	VALUES
    	($1, $2, $3)
	ON CONFLICT (id) DO UPDATE 
	SET text = $2, updated_at = $3;`

	_, err := txExec(ctx, q, note.Id, note.Text, note.UpdatedAt)
	if err != nil {
		return errors.WithMessage(err, "upsert note")
	}

	return nil
}

func deleteNote(ctx context.Context, txExec TxFuncExec, noteId string) error {
	q := `
	DELETE FROM note
	WHERE id = $1;`

	_, err := txExec(ctx, q, noteId)
	if err != nil {
		return errors.WithMessage(err, "delete note")
	}

	return nil
}

func (r Note) GetAll(ctx context.Context) ([]domain.Note, error) {
	q := `
	SELECT id, text, updated_at FROM note;`

	var notes []domain.Note
	rows, err := r.db.Select(ctx, q)
	if err != nil {
		return nil, errors.WithMessage(err, "get notes")
	}

	for rows.Next() {
		var note domain.Note
		err = rows.Scan(&note.Id, &note.Text, &note.UpdatedAt)
		if err != nil {
			return nil, errors.WithMessage(err, "scan row")
		}

		notes = append(notes, note)
	}

	err = rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "close rows")
	}

	return notes, nil
}
