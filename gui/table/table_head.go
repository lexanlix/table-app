package table

import (
	"context"

	"table-app/domain"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

func (t Table) drawTableHead() *core.Frame {
	ctx := context.Background()
	headFrame := core.NewFrame(t.tableFrame)
	headFrame.SetName("headFrame")
	headFrame.Styler(t.styler.HeadFrameStyle())

	// проходим по главным категориям
	for i, categories := range t.data.Categories {
		mainCategName := categories[0].MainCategory

		mainFrame := core.NewFrame(headFrame)
		mainFrame.SetName(mainCategName + "_frame")
		mainFrame.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.CenterAll()
			if i != len(t.data.Categories)-1 {
				s.Border.Width.Right.Dp(2)
			}
			s.Gap.Zero()
		})

		topFrame := core.NewFrame(mainFrame)
		topFrame.SetName(mainCategName + "_topFrame")
		topFrame.Styler(func(s *styles.Style) {
			s.Min.Y.Dp(25)
			s.Padding.Top.Dp(10)
			s.Padding.Bottom.Dp(10)
			s.CenterAll()
		})
		core.NewText(topFrame).SetText(mainCategName)

		bottomFrame := core.NewFrame(mainFrame)
		bottomFrame.SetName(mainCategName + "_bottomFrame")
		bottomFrame.Styler(func(s *styles.Style) {
			s.Gap.Zero()
		})

		bottomFrame.Maker(func(p *tree.Plan) {
			// проходим по категориям главной категории, добавляем ячейки
			for j, category := range t.data.Categories[i] {
				nameLen := len([]rune(t.data.Categories[i][j].Name))

				tree.AddAt(p, "cat_"+t.data.Categories[i][j].Name, func(frame *core.Frame) {
					frame.Styler(func(s *styles.Style) {
						s.Gap.Zero()
						s.Max.X.Dp(t.getCellSizeDpX(nameLen))
						s.Max.Y.Dp(t.settings.Gui.CellSizeDpY)
						s.Border.Width.SetAll(units.Dp(1))
						s.CenterAll()
					})

					tField := core.NewTextField(frame)
					tField.Type = core.TextFieldOutlined
					tField.Styler(func(s *styles.Style) {
						s.Border.Radius.Zero()
						s.Border.Width.Zero()
						s.Border.Offset.Zero()
					})
					core.Bind(&t.data.Categories[i][j].Name, tField.SetText(t.data.Categories[i][j].Name))

					if nameLen > 15 {
						tField.SetTooltip(t.data.Categories[i][j].Name)
					}

					tField.OnChange(func(e events.Event) {
						oldCategory := category
						newCategory := domain.Category{
							Name:         tField.Text(),
							MainCategory: category.MainCategory,
						}

						err := t.controller.UpdateCategoryName(context.Background(), oldCategory, newCategory)
						if err != nil {
							t.logger.Error(ctx, "update category name error: ", log.Any("err", err.Error()))
							core.MessageSnackbar(t.tableFrame, "Ошибка обновления названия категории")
						}

						t.tableFrame.Update()
					})
				})
			}
		})
	}

	consumptionFrame := core.NewFrame(headFrame)
	consumptionFrame.SetName("consumptionFrame")
	consumptionFrame.Styler(t.styler.TableHeadTextFrame())
	core.NewText(consumptionFrame).SetText("Расход в месяц")

	balanceFrame := core.NewFrame(headFrame)
	balanceFrame.SetName("balanceFrame")
	balanceFrame.Styler(t.styler.TableHeadTextFrame())
	core.NewText(balanceFrame).SetText("Остаток")

	return headFrame
}
