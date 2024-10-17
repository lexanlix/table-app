package controller

import (
	"context"
	"time"

	"table-app/internal/log"
)

type UpdatingService interface {
	UpsertAllData(ctx context.Context) error
	SetLastUpdated(updatedTime time.Time) error
	GetLastUpdated() *string
	SetLastRecord(lastRecord string) error
	GetLastRecord() string
}

type Updating struct {
	logger  log.Logger
	service UpdatingService
}

func NewUpdating(logger log.Logger, service UpdatingService) Updating {
	return Updating{
		logger:  logger,
		service: service,
	}
}

// GetLastUpdated
// Получить дату последнего обновления таблицы
func (c Updating) GetLastUpdated() *string {
	return c.service.GetLastUpdated()
}

// GetLastRecord
// Получить данные последней записи
func (c Updating) GetLastRecord() string {
	return c.service.GetLastRecord()
}

// SetLastRecord
// Сохранить данные последней записи
func (c Updating) SetLastRecord(lastRecord string) error {
	return c.service.SetLastRecord(lastRecord)
}
