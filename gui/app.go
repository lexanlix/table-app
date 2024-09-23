package gui

import (
	"context"
	"strconv"
	"strings"
	"time"

	"table-app/conf"
	"table-app/domain"
	"table-app/entity"
	"table-app/internal/log"
	"table-app/utils"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

type App struct {
	logger     log.Logger
	appBody    *core.Body
	toolBar    *core.Toolbar
	controller TableController
	settings   conf.Setting
	updater    *Updater
	sumUpdater *SumUpdater

	frames    []*core.Frame
	txtFields []*core.TextField
	texts     []*core.Text
}

func NewApp(logger log.Logger, cfg Config, controller TableController, settings conf.Setting, shutdownFunc func()) *App {
	body := core.NewBody(cfg.Title)
	body.Styler(func(s *styles.Style) {
		s.Min.Set(units.Dp(cfg.SizeDp))
		s.Pos.Y.Dp(0)
		s.Pos.X.Dp(0)
	})

	updater := NewUpdater(logger)
	sumUpdater := NewSumUpdater(logger, controller)

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
		frames:     []*core.Frame{},
		txtFields:  []*core.TextField{},
		texts:      []*core.Text{},
	}
}

func (a *App) Upgrade(data *domain.GuiTableData) {
	a.createToolbar(data.Categories)

	mainFrame := a.withFrame(a.appBody)
	mainFrame.SetName("app_mainFrame")
	mainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	for year := a.settings.StartYear; year <= time.Now().Year(); year++ {
		a.DrawYearTable(year, mainFrame, data)
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

// DrawYearTable отрисовывает табличку каждого года
func (a *App) DrawYearTable(year int, frame *core.Frame, data *domain.GuiTableData) {
	tableFrame := a.withFrame(frame)
	tableFrame.SetName("table_" + strconv.Itoa(year) + "_Frame")
	tableFrame.Styler(func(s *styles.Style) {
		s.Display = styles.Grid
		s.Columns = 2
		s.Border.Width.Set(units.Dp(1))
		s.Border.Radius.Zero()
		s.Gap.Zero()
		s.CenterAll()
	})

	yearFrame := a.withFrame(tableFrame)
	yearFrame.SetName("yearFrame")
	yearFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(83)
		s.Border.Width.Right.Dp(1)
		s.Border.Width.Bottom.Dp(2)
		s.CenterAll()
		s.Gap.Zero()
	})
	_ = a.withText(yearFrame, strconv.Itoa(year)+" год")

	_ = a.getTableHead(tableFrame, data)
	_ = a.getMonthsColumn(year, tableFrame)
	_ = a.getValuesFrame(year, tableFrame, data)
}

func (a *App) getTableHead(frame *core.Frame, data *domain.GuiTableData) *core.Frame {
	ctx := context.Background()
	headFrame := a.withFrame(frame)
	headFrame.SetName("headFrame")
	headFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
		s.Border.Width.Left.Dp(1)
		s.Border.Width.Bottom.Dp(1)
		s.Gap.Zero()
	})

	// проходим по главным категориям
	for i, categories := range data.Categories {
		mainCategName := categories[0].MainCategory

		mainFrame := a.withFrame(headFrame)
		mainFrame.SetName(mainCategName + "_frame")
		mainFrame.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.CenterAll()
			if i != len(data.Categories)-1 {
				s.Border.Width.Right.Dp(2)
			}
			s.Gap.Zero()
		})

		topFrame := a.withFrame(mainFrame)
		topFrame.SetName(mainCategName + "_topFrame")
		topFrame.Styler(func(s *styles.Style) {
			s.Min.Y.Dp(25)
			s.Padding.Top.Dp(10)
			s.Padding.Bottom.Dp(10)
			s.CenterAll()
		})
		_ = a.withText(topFrame, mainCategName)

		bottomFrame := a.withFrame(mainFrame)
		bottomFrame.SetName(mainCategName + "_bottomFrame")
		bottomFrame.Styler(func(s *styles.Style) {
			s.Gap.Zero()
		})

		bottomFrame.Maker(func(p *tree.Plan) {
			// проходим по категориям главной категории, добавляем ячейки
			for j, category := range data.Categories[i] {
				tree.AddAt(p, "cat_"+data.Categories[i][j].Name, func(frame *core.Frame) {
					frame.Styler(func(s *styles.Style) {
						s.Gap.Zero()
						s.Max.X.Dp(a.settings.Gui.CellSizeDpX)
						s.Max.Y.Dp(a.settings.Gui.CellSizeDpY)
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
					core.Bind(&data.Categories[i][j].Name, tField.SetText(data.Categories[i][j].Name))

					tField.OnChange(func(e events.Event) {
						oldCategory := category
						newCategory := domain.Category{
							Name:         tField.Text(),
							MainCategory: category.MainCategory,
						}

						err := a.controller.UpdateCategoryName(context.Background(), oldCategory, newCategory)
						if err != nil {
							a.logger.Error(ctx, "update category name error: ", log.Any("err", err.Error()))
							core.MessageSnackbar(a.appBody, "Ошибка обновления названия категории")
						}

						a.appBody.Update()
					})
				})
			}
		})
	}

	consumptionFrame := a.withFrame(headFrame)
	consumptionFrame.SetName("consumptionFrame")
	consumptionFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(81)
		s.Border.Width.SetAll(units.Dp(1))
		s.CenterAll()
	})
	a.withText(consumptionFrame, "Расход в месяц")

	balanceFrame := a.withFrame(headFrame)
	balanceFrame.SetName("balanceFrame")
	balanceFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(81)
		s.Border.Width.SetAll(units.Dp(1))
		s.CenterAll()
	})
	a.withText(balanceFrame, "Остаток")

	return headFrame
}

func (a *App) getMonthsColumn(year int, frame *core.Frame) *core.Frame {
	mainFrame := a.withFrame(frame)
	mainFrame.SetName("monthsFrame")
	mainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Border.Width.Right.Dp(1)
		s.Border.Width.Top.Dp(1)
		s.Gap.Zero()
	})

	var lastMonth time.Month
	if year != time.Now().Year() {
		// если год не текущий
		lastMonth = time.December
	} else {
		lastMonth = time.Now().Month()
	}

	for month := 1; month <= int(lastMonth); month++ {
		monthFrame := a.withFrame(mainFrame)
		monthFrame.SetName(time.Month(month).String() + "_frame")
		monthFrame.Styler(func(s *styles.Style) {
			s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
			s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
			s.Gap.Zero()
			s.Border.Width.Top.Dp(1)
			s.Border.Width.Bottom.Dp(1)
			s.CenterAll()

			if year == time.Now().Year() && month == int(lastMonth) {
				s.Background = ColorYellow
			}
		})

		monthName, ok := domain.RusMonths[month]
		if !ok {
			monthName = time.Month(month).String()
		}

		core.NewText(monthFrame).SetText(monthName)
	}

	if year == time.Now().Year() {
		return mainFrame
	}

	// строка итогов года
	resultFrame := a.withFrame(mainFrame)
	resultFrame.SetName("resultFrame")
	resultFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
		s.Gap.Zero()
		s.Border.Width.Top.Dp(1)
		s.Border.Width.Bottom.Dp(1)
		s.CenterAll()
		s.Background = ColorPurple
	})

	resultText := core.NewText(resultFrame).SetText(strconv.Itoa(year))
	resultText.Styler(func(s *styles.Style) {
		s.Font.Weight = styles.WeightBold
	})

	return mainFrame
}

func (a *App) getValuesFrame(year int, frame *core.Frame, data *domain.GuiTableData) *core.Frame {
	mainFrame := a.withFrame(frame)
	mainFrame.SetName("valuesFrame")
	mainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Border.Width.Left.Dp(1)
		s.Border.Width.Top.Dp(1)
		s.Border.Width.Right.Dp(0.5)
		s.Gap.Zero()
	})

	var lastMonth time.Month
	if year != time.Now().Year() {
		// если год не текущий
		lastMonth = time.December
	} else {
		lastMonth = time.Now().Month()
	}

	consumptionArr := make([]*int, 0)

	for month := 1; month <= int(lastMonth); month++ {
		monthFrame := a.withFrame(mainFrame)
		monthFrame.SetName(time.Month(month).String() + "_frame")
		monthFrame.Styler(func(s *styles.Style) {
			s.CenterAll()
			s.Gap.Zero()
		})

		for i, categories := range data.Categories {
			mainCategFrame := core.NewFrame(monthFrame)
			mainCategFrame.SetName(categories[0].MainCategory + "_frame")
			mainCategFrame.Styler(func(s *styles.Style) {
				s.Gap.Zero()
				if i != len(data.Categories)-1 {
					s.Border.Width.Right.Dp(1)
					s.Background = ColorSoftGreen
				} else {
					s.Background = ColorSoftGrey
				}
				s.CenterAll()
			})

			mainCategFrame.Maker(func(p *tree.Plan) {
				for _, category := range data.Categories[i] {
					ctx := context.Background()
					compositeId := category.MainCategory + category.Name + strconv.Itoa(month) + strconv.Itoa(year)

					cellIsCreated := false
					cell, ok := data.ValuesList[compositeId]
					if ok {
						cellIsCreated = true
					}

					tree.AddAt(p, compositeId, func(frame *core.Frame) {
						frame.Styler(func(s *styles.Style) {
							s.Gap.Zero()
							s.Max.X.Dp(a.settings.Gui.CellSizeDpX)
							s.Max.Y.Dp(a.settings.Gui.CellSizeDpY)
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
						a.updater.AddTextField(compositeId, tField)

						tField.OnDoubleClick(func(e events.Event) {
							if !cellIsCreated {
								cell = domain.Cell{
									MainCategory: category.MainCategory,
									Category:     category.Name,
									Month:        time.Month(month),
									Year:         year,
								}
							}

							sumWindow := NewSumWindow(a.logger, frame, cell, a.controller,
								a.updater.updateChan, a.sumUpdater.updateChan)
							sumWindow.Run(tField)
						})

						tField.OnChange(func(e events.Event) {
							val, err := strconv.Atoi(strings.Join(strings.Fields(tField.Text()), ""))
							if err != nil {
								core.MessageSnackbar(mainFrame, "Неверный формат данных: "+err.Error())
								a.logger.Error(ctx, "convert tField to int: "+err.Error())
								return
							}

							if !cellIsCreated {
								cell = domain.Cell{
									MainCategory: category.MainCategory,
									Category:     category.Name,
									Value:        val,
									Month:        time.Month(month),
									Year:         year,
								}
							} else {
								cell.Value = val
							}

							err = a.controller.UpsertValue(ctx, cell)
							if err != nil {
								core.MessageSnackbar(mainFrame, "Ошибка сохранения данных: "+err.Error())
								a.logger.Error(ctx, "save all data")
								return
							}

							a.sumUpdater.updateChan <- entity.MonthYear{
								Month: month,
								Year:  year,
							}

							tField.SetText(FormatInt(val))
							core.MessageSnackbar(mainFrame, "Введено: "+tField.Text())
						})

						if !cellIsCreated {
							tField.SetText("0")
							return
						}

						if cell.Category != category.Name {
							a.logger.Error(context.Background(), "bad cell month value", log.String("category", cell.Category))
							return
						}

						consumptionArr = append(consumptionArr)
						tField.SetText(FormatInt(cell.Value))
					})
				}
			})
		}

		consFrame := core.NewFrame(monthFrame)
		consFrame.SetName("consumptionFrame")

		consFrame.Styler(func(s *styles.Style) {
			s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
			s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
			s.Border.Width.SetAll(units.Dp(1))
			s.CenterAll()
		})

		consText := core.NewText(consFrame)
		a.sumUpdater.AddConsumptionText(month, year, consText)

		balanceFrame := core.NewFrame(monthFrame)
		balanceFrame.SetName("balanceFrame")

		balanceFrame.Styler(func(s *styles.Style) {
			s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
			s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
			s.Border.Width.SetAll(units.Dp(1))
			s.CenterAll()
		})

		balanceText := core.NewText(balanceFrame)
		a.sumUpdater.AddBalanceText(month, year, balanceText)
	}

	if year == time.Now().Year() {
		return mainFrame
	}

	// строка итогов года
	resultFrame := a.withFrame(mainFrame)
	resultFrame.SetName("resultFrame")
	resultFrame.Styler(func(s *styles.Style) {
		s.Gap.Zero()
		s.CenterAll()
	})

	resultByCategoryId := a.controller.GetAnnualResult(year)

	for i, categories := range data.Categories {
		mainCategFrame := core.NewFrame(resultFrame)
		mainCategFrame.SetName(categories[0].MainCategory + "_resultFrame")
		mainCategFrame.Styler(func(s *styles.Style) {
			s.Gap.Zero()
			if i != len(data.Categories)-1 {
				s.Border.Width.Right.Dp(1)
			}
			s.CenterAll()
		})

		mainCategFrame.Maker(func(p *tree.Plan) {
			for _, category := range data.Categories[i] {
				compositeCategory := utils.GetCompositeCategory(category.MainCategory, category.Name)

				tree.AddAt(p, compositeCategory, func(frame *core.Frame) {
					frame.Styler(func(s *styles.Style) {
						s.Gap.Zero()
						s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
						s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
						s.Border.Width.SetAll(units.Dp(1))
						s.Background = ColorPurple
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

					textResult.SetText(FormatInt(result))
				})
			}
		})
	}

	consResFrame := core.NewFrame(resultFrame)
	consResFrame.SetName("consumptionResFrame")

	consResFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
		s.Border.Width.SetAll(units.Dp(1))
		s.Background = ColorPurple
		s.CenterAll()
	})

	consumptionRes := resultByCategoryId[domain.ColumnConsumption]
	core.NewText(consResFrame).SetText(FormatInt(consumptionRes, addMinus)).Styler(func(s *styles.Style) {
		s.Font.Weight = styles.WeightBold
	})

	balanceResFrame := core.NewFrame(resultFrame)
	balanceResFrame.SetName("balanceResFrame")

	balanceResFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(a.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(a.settings.Gui.CellSizeDpY)
		s.Border.Width.SetAll(units.Dp(1))
		s.Background = ColorPurple
		s.CenterAll()
	})

	balanceRes := resultByCategoryId[domain.ColumnBalance]
	core.NewText(balanceResFrame).SetText(FormatInt(balanceRes)).Styler(func(s *styles.Style) {
		s.Font.Weight = styles.WeightBold
	})

	return mainFrame
}

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
			w.SetText("Начальная сумма: " + FormatInt(a.settings.StartMoney))

			w.OnClick(func(e events.Event) {
				core.MessageSnackbar(a.appBody, "Начальная сумма задается в настройках")
			})
		})
	})

	a.toolBar = tbar
}
