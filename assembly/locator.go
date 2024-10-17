package assembly

import (
	"context"

	"table-app/conf"
	"table-app/controller"
	"table-app/domain"
	"table-app/gui"
	"table-app/internal/db"
	"table-app/internal/log"
	"table-app/repository"
	"table-app/service"

	"github.com/pkg/errors"
)

type DB interface {
	db.DB
}

type Locator struct {
	db     DB
	logger log.Logger
}

func NewLocator(db DB, logger log.Logger) Locator {
	return Locator{
		db:     db,
		logger: logger,
	}
}

func (l Locator) Config(ctx context.Context, cfg conf.Remote, shutdownFunc func()) (*gui.App, error) {
	tableRepo := repository.NewTable(l.db, cfg.Storage)
	categoryRepo := repository.NewCategory(l.db, cfg.Storage)
	updatingRepo := repository.NewUpdated(l.db)
	accountRepo := repository.NewAccount(l.db)

	cellsData, err := tableRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get cells")
	}

	categoryData, err := categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get categories")
	}

	if len(categoryData) == 0 {
		categoryData = domain.GetStartingCategories(cfg.Settings.MainCategoryOrder)
	}

	accountData, err := accountRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get accounts")
	}

	if len(accountData) == 0 {
		accountData = domain.GetStartingAccounts()
	}

	cellsCache := repository.NewCellsCache()
	cellsCache.InitCache(cellsData)
	cellsList := cellsCache.GetList()

	categoryCache := repository.NewCategoryCache(cfg.Settings.MainCategoryOrder)
	categoryCache.InitCache(categoryData)
	categoryList := categoryCache.GetCategoryArray()

	accountCache := repository.NewAccountCache()
	accountCache.InitCache(accountData)
	accountList := accountCache.GetListPtr()

	calculationCache := repository.NewCalculationCache(cfg.Settings)
	err = calculationCache.InitCache(cellsList, categoryList)
	if err != nil {
		return nil, errors.WithMessage(err, "init calculation cache")
	}

	isFileStorage := false
	if cfg.Storage.Files != nil {
		isFileStorage = true
	}

	tableService := service.NewTable(l.logger, cellsCache, tableRepo, cfg.Settings, isFileStorage)
	categoryService := service.NewCategory(l.logger, categoryCache, categoryRepo)
	calculationService := service.NewCalculation(calculationCache, cellsCache, categoryCache, accountCache, cfg.Settings)
	updatingService, err := service.NewUpdating(updatingRepo)
	if err != nil {
		return nil, errors.WithMessage(err, "create updating service")
	}

	accountService := service.NewAccount(l.logger, accountRepo, accountCache)

	tableCtrl := controller.NewTable(l.logger, tableService, categoryService, calculationService, updatingService)
	accountCtrl := controller.NewAccount(l.logger, accountService, calculationService)
	updatingCtrl := controller.NewUpdating(l.logger, updatingService)

	guiApp := gui.NewApp(l.logger, gui.NewAppConfig(), tableCtrl, accountCtrl, updatingCtrl, cfg.Settings, shutdownFunc)
	guiApp.Upgrade(&domain.GuiTableData{
		Categories:        categoryList,
		ValuesList:        cellsList,
		MainCategoryOrder: cfg.Settings.MainCategoryOrder,
		Accounts:          accountList,
	})

	return guiApp, nil
}
