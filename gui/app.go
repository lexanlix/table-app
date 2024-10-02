package gui

import (
	"context"
	"time"

	"table-app/conf"
	"table-app/domain"
	"table-app/gui/iface"
	"table-app/gui/table"
	"table-app/gui/updaters"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
)

type App struct {
	logger     log.Logger
	appBody    *core.Body
	toolBar    *core.Toolbar
	controller iface.TableController
	settings   conf.Setting
	updater    *updaters.Updater
	sumUpdater *updaters.SumUpdater
}

func NewApp(logger log.Logger, cfg Config, controller iface.TableController, settings conf.Setting, shutdownFunc func()) *App {
	body := core.NewBody(cfg.Title)
	body.Styler(func(s *styles.Style) {
		s.Min.Set(units.Dp(cfg.SizeDp))
		s.Pos.Y.Dp(0)
		s.Pos.X.Dp(0)
	})

	updater := updaters.NewUpdater(logger)
	sumUpdater := updaters.NewSumUpdater(logger, controller)

	body.OnClose(func(e events.Event) {
		ctx := context.Background()
		logger.Info(ctx, "starting close")

		err := controller.SaveAll(context.Background())
		if err != nil {
			logger.Error(context.Background(), "save all data error: "+err.Error())
		}

		updater.Close()
		sumUpdater.Close()
		shutdownFunc()

		logger.Info(ctx, "close completed")
	})

	return &App{
		logger:     logger,
		appBody:    body,
		controller: controller,
		settings:   settings,
		updater:    updater,
		sumUpdater: sumUpdater,
	}
}

func (a *App) Upgrade(data *domain.GuiTableData) {
	a.createToolbar(data.Categories)

	mainFrame := core.NewFrame(a.appBody)
	mainFrame.SetName("app_mainFrame")
	mainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	for year := a.settings.StartYear; year <= time.Now().Year(); year++ {
		yearTable := table.NewTable(mainFrame, year, data, a.settings, a.controller, a.updater, a.sumUpdater)
		yearTable.Draw()
	}
}

func (a *App) Run() {
	a.updater.Start()
	a.sumUpdater.Start()
	a.appBody.RunMainWindow()
}

func (a *App) Shutdown() {
	a.appBody.Scene.RenderWindow().SystemWindow.CloseReq()
	a.appBody.Close()
}
