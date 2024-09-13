package service

import (
	"context"

	"table-app/conf"
	"table-app/domain"
	"table-app/internal/log"
	"table-app/repository"

	"github.com/pkg/errors"
)

type CategoryRepository interface {
	UpsertAll(ctx context.Context, categories []domain.Category) error
}

type Category struct {
	logger log.Logger
	cache  *repository.CategoryCache
	repo   CategoryRepository
	cfg    conf.Setting
}

func NewCategory(logger log.Logger, cache *repository.CategoryCache, repo CategoryRepository) *Category {
	return &Category{
		logger: logger,
		cache:  cache,
		repo:   repo,
	}
}

func (s *Category) SaveAll(ctx context.Context) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	categories := s.cache.ReadAll()

	err := s.repo.UpsertAll(ctx, categories)
	if err != nil {
		return errors.WithMessage(err, "upsert categories")
	}

	return nil
}

func (s *Category) AddCategory(newCat domain.Category) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	err := s.cache.Insert(newCat)
	if err != nil {
		return errors.WithMessage(err, "insert new category")
	}

	return nil
}

func (s *Category) CategoryIsExist(category domain.Category) bool {
	s.cache.Lock()
	defer s.cache.Unlock()
	return s.cache.IsInCache(category)
}

func (s *Category) UpdateCategory(old, new domain.Category) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	err := s.cache.UpdateCategory(old, new)
	if err != nil {
		return errors.WithMessage(err, "update category")
	}

	return nil
}
