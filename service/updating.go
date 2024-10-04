package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

const timeFormatUpdated = "02.01.2006 15:04:05"

type UpdatingRepository interface {
	UpsertLastUpdated(ctx context.Context, updatedTime time.Time) error
	GetLastUpdated(ctx context.Context) (time.Time, error)
	UpsertLastRecord(ctx context.Context, note string) error
	GetLastRecord(ctx context.Context) (string, error)
}

type Updating struct {
	repo UpdatingRepository

	updatedTime *string
	lastRecord  string
}

func NewUpdating(repo UpdatingRepository) (*Updating, error) {
	updatedTime, err := repo.GetLastUpdated(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "get last updated time")
	}
	updatedTimeStr := updatedTime.Format(timeFormatUpdated)

	lastRecord, err := repo.GetLastRecord(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "get last record")
	}

	return &Updating{
		repo:        repo,
		updatedTime: &updatedTimeStr,
		lastRecord:  lastRecord,
	}, nil
}

func (s *Updating) UpsertAllData(ctx context.Context) error {
	err := s.repo.UpsertLastRecord(ctx, s.lastRecord)
	if err != nil {
		return errors.WithMessage(err, "upsert last record")
	}

	timeParsed, err := time.Parse(timeFormatUpdated, *s.updatedTime)
	if err != nil {
		return errors.WithMessage(err, "parse time str")
	}

	err = s.repo.UpsertLastUpdated(ctx, timeParsed)
	if err != nil {
		return errors.WithMessage(err, "upsert last updated")
	}

	return nil
}

func (s *Updating) SetLastUpdated(updatedTime time.Time) error {
	updatedTimeStr := updatedTime.Format(timeFormatUpdated)
	*s.updatedTime = updatedTimeStr

	return nil
}

func (s *Updating) GetLastUpdated() *string {
	return s.updatedTime
}

func (s *Updating) SetLastRecord(lastRecord string) error {
	s.lastRecord = lastRecord
	return nil
}

func (s *Updating) GetLastRecord() string {
	return s.lastRecord
}
