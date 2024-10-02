package table

import (
	"context"
	"strconv"
	"strings"
	"time"

	"table-app/domain"
	"table-app/entity"
	"table-app/gui/styles/colors"
	"table-app/gui/styles/format"
	"table-app/internal/log"
	"table-app/utils"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

func (t Table) drawValuesGrid() *core.Frame {
	mainFrame := core.NewFrame(t.tableFrame)
	mainFrame.SetName("valuesFrame")
	mainFrame.Styler(t.styler.ValuesFrameStyle())

	var lastMonth time.Month
	if t.year != time.Now().Year() {
		// если год не текущий
		lastMonth = time.December
	} else {
		lastMonth = time.Now().Month()
	}

	for month := 1; month <= int(lastMonth); month++ {
		monthFrame := core.NewFrame(mainFrame)
		monthFrame.SetName(time.Month(month).String() + "_frame")
		monthFrame.Styler(func(s *styles.Style) {
			s.CenterAll()
			s.Gap.Zero()
		})

		for i, categories := range t.data.Categories {
			mainCategFrame := core.NewFrame(monthFrame)
			mainCategFrame.SetName(categories[0].MainCategory + "_frame")
			mainCategFrame.Styler(func(s *styles.Style) {
				s.Gap.Zero()
				if i != len(t.data.Categories)-1 {
					s.Border.Width.Right.Dp(1)
					s.Background = colors.ColorSoftGreen
				} else {
					s.Background = colors.ColorSoftGrey
				}
				s.CenterAll()
			})

			mainCategFrame.Maker(func(p *tree.Plan) {
				for j, category := range t.data.Categories[i] {
					ctx := context.Background()
					compositeId := category.MainCategory + category.Name + strconv.Itoa(month) + strconv.Itoa(t.year)

					cellIsCreated := false
					cell, ok := t.controller.GetCellById(compositeId)
					if ok {
						cellIsCreated = true
					}

					tree.AddAt(p, compositeId, func(frame *core.Frame) {
						nameLen := len([]rune(t.data.Categories[i][j].Name))

						frame.Styler(func(s *styles.Style) {
							s.Gap.Zero()
							s.Max.X.Dp(t.getCellSizeDpX(nameLen))
							s.Max.Y.Dp(t.settings.Gui.CellSizeDpY)
							s.Border.Width.SetAll(units.Dp(1))
							s.CenterAll()
						})

						tField := core.NewTextField(frame)
						tField.SetName(compositeId + "_tField")
						tField.Type = core.TextFieldOutlined
						tField.Styler(func(s *styles.Style) {
							s.Border.Radius.Zero()
							s.Border.Width.Zero()
							s.Border.Offset.Zero()
						})
						t.updater.AddTextField(compositeId, tField)

						tField.OnDoubleClick(func(e events.Event) {
							// проверка через кеш
							cell, ok = t.controller.GetCellById(compositeId)
							if ok {
								cellIsCreated = true
							}
							if !cellIsCreated {
								cell = domain.Cell{
									MainCategory: category.MainCategory,
									Category:     category.Name,
									Month:        time.Month(month),
									Year:         t.year,
								}
							}

							sumWindow := NewSumWindow(t.logger, frame, cell, t.controller,
								t.updater.GetUpdateChan(), t.sumUpdater.GetUpdateChan())
							sumWindow.Run(tField)
						})

						tField.OnChange(func(e events.Event) {
							val, err := strconv.Atoi(strings.Join(strings.Fields(tField.Text()), ""))
							if err != nil {
								core.MessageSnackbar(mainFrame, "Неверный формат данных: "+err.Error())
								t.logger.Error(ctx, "convert tField to int: "+err.Error())
								return
							}

							if !cellIsCreated {
								cell = domain.Cell{
									MainCategory: category.MainCategory,
									Category:     category.Name,
									Value:        val,
									Month:        time.Month(month),
									Year:         t.year,
								}
							} else {
								cell.Value = val
							}

							err = t.controller.UpsertValue(ctx, cell)
							if err != nil {
								core.MessageSnackbar(mainFrame, "Ошибка сохранения данных: "+err.Error())
								t.logger.Error(ctx, "save all data")
								return
							}

							t.sumUpdater.SendToChannel(entity.MonthYear{
								Month: month,
								Year:  t.year,
							})

							tField.SetText(format.FormatInt(val))
							core.MessageSnackbar(mainFrame, "Введено: "+tField.Text())
						})

						if !cellIsCreated {
							tField.SetText("")
							return
						}

						if cell.Category != category.Name {
							t.logger.Error(context.Background(), "bad cell month value", log.String("category", cell.Category))
							return
						}

						tField.SetText(format.FormatInt(cell.Value))
					})
				}
			})
		}

		consFrame := core.NewFrame(monthFrame)
		consFrame.SetName("consumptionFrame")

		consFrame.Styler(t.styler.StandardSizeWithAllBorders())

		consText := core.NewText(consFrame)
		t.sumUpdater.AddConsumptionText(month, t.year, consText)

		balanceFrame := core.NewFrame(monthFrame)
		balanceFrame.SetName("balanceFrame")

		balanceFrame.Styler(t.styler.StandardSizeWithAllBorders())

		balanceText := core.NewText(balanceFrame)
		t.sumUpdater.AddBalanceText(month, t.year, balanceText)
	}

	if t.year == time.Now().Year() {
		return mainFrame
	}

	// строка итогов года
	resultFrame := core.NewFrame(mainFrame)
	resultFrame.SetName("resultFrame")
	resultFrame.Styler(func(s *styles.Style) {
		s.Gap.Zero()
		s.CenterAll()
	})

	resultByCategoryId := t.controller.GetAnnualResult(t.year)

	for i, categories := range t.data.Categories {
		mainCategFrame := core.NewFrame(resultFrame)
		mainCategFrame.SetName(categories[0].MainCategory + "_resultFrame")
		mainCategFrame.Styler(func(s *styles.Style) {
			s.Gap.Zero()
			if i != len(t.data.Categories)-1 {
				s.Border.Width.Right.Dp(1)
			}
			s.CenterAll()
		})

		mainCategFrame.Maker(func(p *tree.Plan) {
			for j, category := range t.data.Categories[i] {
				compositeCategory := utils.GetCompositeCategory(category.MainCategory, category.Name)

				tree.AddAt(p, compositeCategory, func(frame *core.Frame) {
					nameLen := len([]rune(t.data.Categories[i][j].Name))

					frame.Styler(func(s *styles.Style) {
						s.Gap.Zero()
						s.Min.X.Dp(t.getCellSizeDpX(nameLen))
						s.Min.Y.Dp(t.settings.Gui.CellSizeDpY)
						s.Border.Width.SetAll(units.Dp(1))
						s.Background = colors.ColorPurple
						s.CenterAll()
					})

					textResult := core.NewText(frame)
					textResult.SetName(compositeCategory + "_text")
					textResult.Styler(func(s *styles.Style) {
						s.Font.Weight = styles.WeightBold
					})

					var result int
					var ok bool

					result, ok = resultByCategoryId[compositeCategory]
					if !ok {
						result = 0
					}

					textResult.SetText(format.FormatInt(result))
				})
			}
		})
	}

	consResFrame := core.NewFrame(resultFrame)
	consResFrame.SetName("consumptionResFrame")

	consResFrame.Styler(t.styler.StandardSizeWithAllBorders(colors.ColorPurple))

	consumptionRes := resultByCategoryId[domain.ColumnConsumption]
	core.NewText(consResFrame).SetText(format.FormatInt(consumptionRes, format.AddMinus)).Styler(func(s *styles.Style) {
		s.Font.Weight = styles.WeightBold
	})

	balanceResFrame := core.NewFrame(resultFrame)
	balanceResFrame.SetName("balanceResFrame")

	balanceResFrame.Styler(t.styler.StandardSizeWithAllBorders(colors.ColorPurple))

	balanceRes := resultByCategoryId[domain.ColumnBalance]
	core.NewText(balanceResFrame).SetText(format.FormatInt(balanceRes)).Styler(func(s *styles.Style) {
		s.Font.Weight = styles.WeightBold
	})

	return mainFrame
}
