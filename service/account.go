package service

import (
	"context"

	"table-app/domain"
	"table-app/internal/log"
	"table-app/repository"

	"github.com/pkg/errors"
)

type AccountRepository interface {
	UpsertAll(ctx context.Context, accounts []domain.Account) error
	GetAll(ctx context.Context) ([]domain.Account, error)
}

type Account struct {
	logger log.Logger
	repo   AccountRepository
	cache  *repository.AccountCache
}

func NewAccount(logger log.Logger, repo AccountRepository, cache *repository.AccountCache) *Account {
	return &Account{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

func (s *Account) SaveAll(ctx context.Context) error {
	s.cache.Lock()
	all := s.cache.ReadAll()
	s.cache.Unlock()

	err := s.repo.UpsertAll(ctx, all)
	if err != nil {
		return errors.WithMessage(err, "upsert all")
	}

	return nil
}

func (s *Account) GetAll(ctx context.Context) ([]domain.Account, error) {
	s.cache.Lock()
	defer s.cache.Unlock()

	all := s.cache.ReadAll()
	return all, nil
}

func (s *Account) Insert(ctx context.Context, newAcc domain.Account) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	err := s.cache.Insert(newAcc)
	if err != nil {
		return errors.WithMessage(err, "insert account")
	}

	return nil
}

func (s *Account) Update(ctx context.Context, account domain.Account) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	err := s.cache.UpdateAccount(account)
	if err != nil {
		return errors.WithMessage(err, "update account")
	}

	return nil
}

func (s *Account) UpdateList(ctx context.Context, list []domain.Account) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	for _, acc := range list {
		if len(acc.Id) == 0 { // новый счет
			if acc.Deleted { // уже удален
				continue
			}

			err := s.cache.Insert(acc)
			if err != nil {
				return errors.WithMessage(err, "insert account")
			}
			continue
		}

		if acc.Deleted { // счет удален
			s.cache.Delete(acc)
			continue
		}

		err := s.cache.UpdateAccount(acc)
		if err != nil {
			return errors.WithMessage(err, "update account")
		}
	}

	return nil
}
