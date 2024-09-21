package controller

import (
	"context"

	"table-app/domain"
	"table-app/internal/log"

	"github.com/pkg/errors"
)

type TableService interface {
	Upsert(cell domain.Cell) error
	UpdateCategoryName(oldCateg, newCateg domain.Category)
	SaveAll(ctx context.Context) error
}

type CategoryService interface {
	AddCategory(newCat domain.Category) error
	CategoryIsExist(category domain.Category) bool
	UpdateCategory(old, new domain.Category) error
	SaveAll(ctx context.Context) error
}

type CalculationService interface {
	ConsumptionSum(month, year int) int
	UpsertBalance(month, year int) (map[string]int, error)
	BalanceSum(month, year int) (int, error)
}

type Table struct {
	logger             log.Logger
	service            TableService
	categoryService    CategoryService
	calculationService CalculationService
}

func NewTable(logger log.Logger, service TableService, categoryService CategoryService, calculationService CalculationService) Table {
	return Table{
		logger:             logger,
		service:            service,
		categoryService:    categoryService,
		calculationService: calculationService,
	}
}

// UpsertValue
// Обновление/добавление нового значения в кеш ячеек
func (c Table) UpsertValue(ctx context.Context, cell domain.Cell) error {
	c.logger.Debug(ctx, "upsert new cell value",
		log.String("category", cell.Category),
		log.Int("value", cell.Value))

	err := cell.Validate()
	if err != nil {
		return errors.WithMessage(err, "validate cell")
	}

	return c.service.Upsert(cell)
}

// SaveAll
// Сохранение всех кешей в БД
func (c Table) SaveAll(ctx context.Context) error {
	c.logger.Debug(ctx, "save all")

	// сначала сохраняются изменения в категориях, так как в таблице ячеек обновляются
	// названия категорий, и по ним далее идет обновление значений ячеек
	err := c.categoryService.SaveAll(ctx)
	if err != nil {
		return errors.WithMessage(err, "save all categories")
	}

	return c.service.SaveAll(ctx)
}

// AddCategory
// Добавление новой категории в кеш категорий
func (c Table) AddCategory(ctx context.Context, category domain.Category) error {
	c.logger.Debug(ctx, "add category",
		log.String("mainCategory", category.MainCategory),
		log.String("category", category.Name))

	return c.categoryService.AddCategory(category)
}

// UpdateCategoryName
// Обновление названия категории в кеше ячеек и в кеше категорий
func (c Table) UpdateCategoryName(ctx context.Context, old, new domain.Category) error {
	c.logger.Debug(ctx, "update category name",
		log.String("old category", old.Name),
		log.String("new category", new.Name))

	err := c.categoryService.UpdateCategory(old, new)
	if err != nil {
		return errors.WithMessage(err, "update category")
	}

	c.service.UpdateCategoryName(old, new)
	return nil
}

// CategoryIsExist
// Поиск категории в кеше
func (c Table) CategoryIsExist(ctx context.Context, category domain.Category) bool {
	return c.categoryService.CategoryIsExist(category)
}

// GetConsumptionSum
// Получение суммы расходов по конкретному месяцу и году
func (c Table) GetConsumptionSum(month, year int) int {
	return c.calculationService.ConsumptionSum(month, year)
}

// GetBalanceSum
// Получение суммы остатка по конкретному месяцу и году
func (c Table) GetBalanceSum(month, year int) (int, error) {
	res, err := c.calculationService.BalanceSum(month, year)
	if err != nil {
		return 0, errors.WithMessage(err, "get balance sum")
	}

	return res, nil
}

// UpsertBalance
// Обновление остатка
func (c Table) UpsertBalance(month, year int) (map[string]int, error) {
	res, err := c.calculationService.UpsertBalance(month, year)
	if err != nil {
		return nil, errors.WithMessage(err, "upsert balance")
	}

	return res, nil
}
