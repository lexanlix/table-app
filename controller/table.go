package controller

import (
	"context"
	"time"

	"table-app/domain"
	"table-app/internal/log"

	"github.com/pkg/errors"
)

type TableService interface {
	Upsert(cell domain.Cell) error
	UpdateCategoryName(oldCateg, newCateg domain.Category)
	SaveAll(ctx context.Context) error
	GetCellById(compositeId string) (domain.Cell, bool)
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
	GetAnnualResult(year int) map[string]int
}

type UpdatingService interface {
	UpsertAllData(ctx context.Context) error
	SetLastUpdated(updatedTime time.Time) error
	GetLastUpdated() *string
	SetLastRecord(lastRecord string) error
	GetLastRecord() string
}

type Table struct {
	logger             log.Logger
	service            TableService
	categoryService    CategoryService
	calculationService CalculationService
	updatingService    UpdatingService
}

func NewTable(
	logger log.Logger,
	service TableService,
	categoryService CategoryService,
	calculationService CalculationService,
	updatingService UpdatingService,
) Table {
	return Table{
		logger:             logger,
		service:            service,
		categoryService:    categoryService,
		calculationService: calculationService,
		updatingService:    updatingService,
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

	err = c.service.Upsert(cell)
	if err != nil {
		return errors.WithMessage(err, "upsert cell")
	}

	err = c.updatingService.SetLastUpdated(time.Now())
	if err != nil {
		return errors.WithMessage(err, "set last updated time")
	}

	return nil
}

// SaveAll
// Сохранение всех кешей в БД
func (c Table) SaveAll(ctx context.Context) error {
	c.logger.Debug(ctx, "save all tables data")

	// сначала сохраняются изменения в категориях, так как в таблице ячеек обновляются
	// названия категорий, и по ним далее идет обновление значений ячеек
	err := c.categoryService.SaveAll(ctx)
	if err != nil {
		return errors.WithMessage(err, "save all categories")
	}

	err = c.service.SaveAll(ctx)
	if err != nil {
		return errors.WithMessage(err, "save all cells")
	}

	err = c.updatingService.UpsertAllData(ctx)
	if err != nil {
		return errors.WithMessage(err, "upsert all updated data")
	}

	return nil
}

// AddCategory
// Добавление новой категории в кеш категорий
func (c Table) AddCategory(ctx context.Context, category domain.Category) error {
	c.logger.Debug(ctx, "add category",
		log.String("mainCategory", category.MainCategory),
		log.String("category", category.Name))

	err := c.categoryService.AddCategory(category)
	if err != nil {
		return errors.WithMessage(err, "add category")
	}

	err = c.updatingService.SetLastUpdated(time.Now())
	if err != nil {
		return errors.WithMessage(err, "set last updated time")
	}
	return nil
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

	err = c.updatingService.SetLastUpdated(time.Now())
	if err != nil {
		return errors.WithMessage(err, "set last updated time")
	}
	return nil
}

// GetCellById
// Получить ячейку по compositeId
func (c Table) GetCellById(compositeId string) (domain.Cell, bool) {
	return c.service.GetCellById(compositeId)
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

// GetAnnualResult
// Годовой итог
func (c Table) GetAnnualResult(year int) map[string]int {
	return c.calculationService.GetAnnualResult(year)
}

// GetLastUpdated
// Получить дату последнего обновления таблицы
func (c Table) GetLastUpdated() *string {
	return c.updatingService.GetLastUpdated()
}

// GetLastRecord
// Получить данные последней записи
func (c Table) GetLastRecord() string {
	return c.updatingService.GetLastRecord()
}

// SetLastRecord
// Сохранить данные последней записи
func (c Table) SetLastRecord(lastRecord string) error {
	return c.updatingService.SetLastRecord(lastRecord)
}
