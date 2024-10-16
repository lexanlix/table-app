package iface

import (
	"context"

	"table-app/domain"
)

type TableController interface {
	UpsertValue(ctx context.Context, cell domain.Cell) error
	AddCategory(ctx context.Context, category domain.Category) error
	UpdateCategoryName(ctx context.Context, old, new domain.Category) error
	CategoryIsExist(ctx context.Context, category domain.Category) bool
	SaveAll(ctx context.Context) error
	GetCellById(compositeId string) (domain.Cell, bool)

	GetConsumptionSum(month, year int) int
	GetBalanceSum(month, year int) (int, error)
	UpsertBalance(month, year int) (map[string]int, error)

	GetAnnualResult(year int) map[string]int

	GetLastUpdated() *string
	GetLastRecord() string
	SetLastRecord(lastRecord string) error
}

type AccountController interface {
	SaveAll(ctx context.Context) error
	GetAll(ctx context.Context) ([]domain.Account, error)
	AddAccount(ctx context.Context, account domain.Account) error
	UpdateAccount(ctx context.Context, account domain.Account) error
	UpdateList(ctx context.Context, list []domain.Account) error
	GetSum(ctx context.Context) (domain.Sum, error)
}
