package controller

import (
	"context"

	"table-app/domain"
	"table-app/internal/log"
)

type AccountService interface {
	SaveAll(ctx context.Context) error
	GetAll(ctx context.Context) ([]domain.Account, error)
	Insert(ctx context.Context, newAcc domain.Account) error
	Update(ctx context.Context, account domain.Account) error
	UpdateList(ctx context.Context, list []domain.Account) error
}

type Account struct {
	logger             log.Logger
	service            AccountService
	calculationService CalculationService
}

func NewAccount(logger log.Logger, service AccountService, calculationService CalculationService) Account {
	return Account{
		logger:             logger,
		service:            service,
		calculationService: calculationService,
	}
}

func (c Account) SaveAll(ctx context.Context) error {
	c.logger.Debug(ctx, "save all accounts")

	return c.service.SaveAll(ctx)
}

func (c Account) AddAccount(ctx context.Context, account domain.Account) error {
	c.logger.Debug(ctx, "add new account", log.String("name", account.Name))

	return c.service.Insert(ctx, account)
}

func (c Account) UpdateList(ctx context.Context, list []domain.Account) error {
	c.logger.Debug(ctx, "update accounts list")

	return c.service.UpdateList(ctx, list)
}

func (c Account) UpdateAccount(ctx context.Context, account domain.Account) error {
	c.logger.Debug(ctx, "update account", log.String("name", account.Name))

	return c.service.Update(ctx, account)
}

func (c Account) GetAll(ctx context.Context) ([]domain.Account, error) {
	c.logger.Debug(ctx, "get all accounts")

	return c.service.GetAll(ctx)
}

func (c Account) GetSum(ctx context.Context) (domain.Sum, error) {
	c.logger.Debug(ctx, "get sum")

	return c.calculationService.GetSum()
}
