package service

import (
	"context"

	"table-app/conf"
	"table-app/domain"
	"table-app/internal/log"
	"table-app/repository"

	"github.com/pkg/errors"
)

type TableRepository interface {
	UpsertAll(ctx context.Context, cells []domain.Cell) error
	GetAll(ctx context.Context) ([]domain.Cell, error)
}

type Table struct {
	logger log.Logger
	cache  *repository.CellsCache
	repo   TableRepository
	cfg    conf.Setting
}

func NewTable(logger log.Logger, cache *repository.CellsCache, repo TableRepository, cfg conf.Setting) *Table {
	return &Table{
		logger: logger,
		cache:  cache,
		repo:   repo,
		cfg:    cfg,
	}
}

func (s *Table) Upsert(cell domain.Cell) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	s.cache.Upsert(cell)

	return nil
}

func (s *Table) SaveAll(ctx context.Context) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	cells := s.cache.ReadAll()

	updatedCells := make([]domain.Cell, 0)
	for _, cell := range cells {
		if !cell.IsUpdated {
			continue
		}

		updatedCells = append(updatedCells, cell)
	}

	err := s.repo.UpsertAll(ctx, updatedCells)
	if err != nil {
		return errors.WithMessage(err, "upsert cells")
	}

	return nil
}

func (s *Table) UpdateCategoryName(oldCateg, newCateg domain.Category) {
	s.cache.Lock()
	defer s.cache.Unlock()

	s.cache.UpdateCategoryName(oldCateg, newCateg, s.cfg.StartMonth, s.cfg.StartYear)
}
