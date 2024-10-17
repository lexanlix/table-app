package gui

import (
	"context"

	"table-app/domain"
	"table-app/gui/dialogs"
	"table-app/gui/styles/format"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/tree"
)

func (a *App) createToolbar(categories [][]domain.Category) {
	ctx := context.Background()

	tbar := core.NewToolbar(a.appBody)
	tbar.Maker(func(p *tree.Plan) {
		tree.Add(p, func(w *core.Button) {
			w.SetText("Сохранить")
			w.OnClick(func(e events.Event) {
				err := a.tableController.SaveAll(ctx)
				if err != nil {
					core.ErrorSnackbar(a.appBody, err, "Ошибка сохранения данных")
					a.logger.Error(context.Background(), "save all data error", log.Any("err", err.Error()))
					return
				}

				err = a.accountController.SaveAll(ctx)
				if err != nil {
					core.ErrorSnackbar(a.appBody, err, "Ошибка сохранения данных")
					a.logger.Error(context.Background(), "save all data error", log.Any("err", err.Error()))
					return
				}

				core.MessageSnackbar(a.appBody, "-Данные сохранены-")
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Новая категория")
			w.OnClick(func(e events.Event) {
				categoryDialog := dialogs.NewCategoryDialog(a.logger, a.appBody, a.tableController, categories)
				categoryDialog.Run()
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Счета")
			w.OnClick(func(e events.Event) {
				accountsDialog, err := dialogs.NewAccountsDialog(
					a.logger, a.appBody, a.accountController, a.sumUpdater.GetUpdateAccountsChan(),
				)
				if err != nil {
					core.ErrorSnackbar(a.appBody, err, "Ошибка сохранения данных")
					a.logger.Error(context.Background(), "save all data error", log.Any("err", err.Error()))
					return
				}
				accountsDialog.Run()
			})
		})
		tree.Add(p, func(w *core.Stretch) {
			w.SetName("stretch")
		})
		tree.Add(p, func(w *core.Text) {
			w.SetText("Начальная сумма: " + format.FormatInt(a.settings.StartMoney))

			w.OnClick(func(e events.Event) {
				core.MessageSnackbar(a.appBody, "Начальная сумма задается в настройках")
			})
		})
	})

	a.toolBar = tbar
}
