package gui

import (
	"context"

	"table-app/domain"
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
				err := a.controller.SaveAll(ctx)
				if err != nil {
					core.MessageSnackbar(a.appBody, "Ошибка сохранения данных: "+err.Error())
					a.logger.Error(context.Background(), "save all data error", log.Any("err", err.Error()))
					return
				}

				core.MessageSnackbar(a.appBody, "-Данные сохранены-")
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Новая категория")
			w.OnClick(func(e events.Event) {
				categoryWindow := NewCategoryWindow(a.logger, a.appBody, a.controller, categories)
				categoryWindow.Run()
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
