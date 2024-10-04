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

	cellsData, err := tableRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get cells")
	}

	categoryList, err := categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get categories")
	}

	if len(categoryList) == 0 {
		categoryList = domain.GetStartingCategories(cfg.Settings.MainCategoryOrder)
	}

	cellsCache := repository.NewCellsCache()
	cellsCache.InitCache(cellsData)
	cellsList := cellsCache.GetList()

	categoryCache := repository.NewCategoryCache(cfg.Settings.MainCategoryOrder)
	categoryCache.InitCache(categoryList)
	categoryArray := categoryCache.GetCategoryArray()

	calculationCache := repository.NewCalculationCache(cfg.Settings)
	err = calculationCache.InitCache(cellsList, categoryArray)
	if err != nil {
		return nil, errors.WithMessage(err, "init calculation cache")
	}

	isFileStorage := false
	if cfg.Storage.Files != nil {
		isFileStorage = true
	}

	tableService := service.NewTable(l.logger, cellsCache, tableRepo, cfg.Settings, isFileStorage)
	categoryService := service.NewCategory(l.logger, categoryCache, categoryRepo)
	calculationService := service.NewCalculation(calculationCache, cellsCache, categoryCache, cfg.Settings)
	updatingService, err := service.NewUpdating(updatingRepo)
	if err != nil {
		return nil, errors.WithMessage(err, "create updating service")
	}

	tableCtrl := controller.NewTable(l.logger, tableService, categoryService, calculationService, updatingService)

	guiApp := gui.NewApp(l.logger, gui.NewAppConfig(), tableCtrl, cfg.Settings, shutdownFunc)

	guiApp.Upgrade(&domain.GuiTableData{
		Categories:        categoryArray,
		ValuesList:        cellsList,
		MainCategoryOrder: cfg.Settings.MainCategoryOrder,
	})

	return guiApp, nil
}
