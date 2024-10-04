package repository

import (
	"context"
	"database/sql"
	"time"

	"table-app/internal/db"

	"github.com/pkg/errors"
)

const (
	lastUpdatedTimeId   = "last_updated_time_id"
	lastUpdatedRecordId = "last_updated_record_id"
)

type Updated struct {
	db db.DB
}

func NewUpdated(db db.DB) Updated {
	return Updated{
		db: db,
	}
}

func (r Updated) UpsertLastRecord(ctx context.Context, note string) error {
	q := `
	INSERT INTO last_updated_data
    	(id, note)
	VALUES
    	($1, $2)
	ON CONFLICT (id) DO UPDATE SET note = $2;`

	_, err := r.db.Exec(ctx, q, lastUpdatedRecordId, note)
	if err != nil {
		return errors.WithMessage(err, "upsert last record")
	}

	return nil
}

func (r Updated) GetLastRecord(ctx context.Context) (string, error) {
	var lastRecord string

	q := `
	SELECT note 
	FROM last_updated_data
    WHERE id = $1;`

	row := r.db.SelectRow(ctx, q, lastUpdatedRecordId)
	err := row.Scan(&lastRecord)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return "", nil
		}

		return "", errors.WithMessage(err, "get last updated record")
	}

	return lastRecord, nil
}

func (r Updated) UpsertLastUpdated(ctx context.Context, updatedTime time.Time) error {
	q := `
	INSERT INTO last_updated_data
    	(id, updated_at)
	VALUES
    	($1, $2)
	ON CONFLICT (id) DO UPDATE SET updated_at = $2;`

	_, err := r.db.Exec(ctx, q, lastUpdatedTimeId, updatedTime)
	if err != nil {
		return errors.WithMessage(err, "upsert updated time")
	}

	return nil
}

func (r Updated) GetLastUpdated(ctx context.Context) (time.Time, error) {
	var updatedTime time.Time

	q := `
	SELECT updated_at 
	FROM last_updated_data
    WHERE id = $1;`

	row := r.db.SelectRow(ctx, q, lastUpdatedTimeId)
	err := row.Scan(&updatedTime)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return time.Now(), nil
		}

		return time.Time{}, errors.WithMessage(err, "get last updated time")
	}

	return updatedTime, nil
}
