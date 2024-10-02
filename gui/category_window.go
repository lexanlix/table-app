package gui

import (
	"context"

	"table-app/domain"
	"table-app/gui/iface"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
)

type CategoryWindow struct {
	logger  log.Logger
	appBody *core.Body

	catDialog      *core.Body
	controller     iface.TableController
	mainCategories []string
	category       domain.Category
}

func NewCategoryWindow(logger log.Logger, appBody *core.Body, controller iface.TableController,
	categories [][]domain.Category) *CategoryWindow {
	catBody := core.NewBody("NewCategory").SetTitle("Добавление категории")
	catBody.Styler(func(s *styles.Style) {
		s.Align.Self = styles.Center
		s.CenterAll()
	})

	mainCatFrame := core.NewFrame(catBody)
	mainCatFrame.SetName("mainCatFrame")
	mainCatFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	catFrame := core.NewFrame(mainCatFrame)
	catFrame.SetName("catFrame")

	core.NewSpace(mainCatFrame).Styler(func(s *styles.Style) {
		s.Min.Y.Dp(10)
	})

	buttonsFrame := core.NewFrame(mainCatFrame)
	buttonsFrame.SetName("buttonsFrame")

	catWindow := &CategoryWindow{
		logger:         logger,
		appBody:        appBody,
		catDialog:      catBody,
		controller:     controller,
		mainCategories: getMainCategories(categories),
		category:       domain.Category{},
	}

	catWindow.addInput(catFrame)
	catWindow.addButtons(buttonsFrame)

	return catWindow
}

func (s *CategoryWindow) addInput(catFrame *core.Frame) {
	catFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	titleFrame := core.NewFrame(catFrame)
	titleFrame.SetName("titleFrame")
	titleFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(400)
		s.CenterAll()
	})
	core.NewText(titleFrame).SetType(core.TextDisplaySmall).SetText("Добавить категорию")

	core.NewSpace(catFrame).Styler(func(s *styles.Style) {
		s.Min.Y.Dp(10)
	})

	choiceFrame := core.NewFrame(catFrame)
	choiceFrame.SetName("choiceMainCatFrame")
	choiceFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Min.X.Dp(300)
	})
	core.NewText(choiceFrame).SetType(core.TextBodyLarge).SetText("Выбрать главную категорию")
	switcher := core.NewSwitches(choiceFrame).SetType(core.SwitchChip).SetMutex(true).SetStrings(s.mainCategories...)
	switcher.SetName("switcherMain")
	switcher.Styler(func(s *styles.Style) {
		s.Font.Size.Set(8, units.UnitPt)
	})
	switcher.OnChange(func(e events.Event) {
		if switcher.SelectedItem() != nil {
			s.category.MainCategory = switcher.SelectedItem().Value.(string)
		}
	})

	core.NewSpace(catFrame).Styler(func(s *styles.Style) {
		s.Min.Y.Dp(10)
	})

	inputFrame := core.NewFrame(catFrame)
	inputFrame.SetName("inputFrame")
	inputFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})
	core.NewText(inputFrame).SetType(core.TextBodyLarge).SetText("Название категории")
	tField := core.NewTextField(inputFrame).SetType(core.TextFieldOutlined)
	tField.Styler(func(s *styles.Style) {
		s.Min.X.Dp(300)
		s.Font.Size.Set(8, units.UnitPt)
	})
	tField.OnInput(func(e events.Event) {
		s.category.Name = tField.Text()
	})
}

func (s *CategoryWindow) addButtons(buttonsFrame *core.Frame) {
	buttonsFrame.Styler(func(s *styles.Style) {
		s.Max.X.Dp(400)
		s.CenterAll()
	})

	cancelButton := core.NewButton(buttonsFrame).SetType(core.ButtonElevated).SetText("Отмена")
	cancelButton.OnClick(func(e events.Event) {
		s.close()
	})

	core.NewStretch(buttonsFrame)

	addButton := core.NewButton(buttonsFrame).SetType(core.ButtonFilled).SetText("Добавить")
	addButton.OnClick(func(e events.Event) {
		if len(s.category.MainCategory) == 0 || len(s.category.Name) == 0 {
			core.MessageSnackbar(s.catDialog, "Введите данные для добавления категории")
			return
		}

		if s.controller.CategoryIsExist(context.Background(), s.category) {
			core.MessageSnackbar(s.catDialog, "Категория уже существует")
			return
		}

		err := s.controller.AddCategory(context.Background(), s.category)
		if err != nil {
			core.MessageSnackbar(s.catDialog, "Ошибка добавления категории: "+err.Error())
			s.logger.Error(context.Background(), "add category error", log.Any("err", err.Error()))
			return
		}

		s.appBody.Update()
		s.close()
	})
}

func (s *CategoryWindow) Run() {
	stage := s.catDialog.NewDialog(s.appBody)

	stage.Pos.X = int(s.appBody.Geom.Size.Actual.Total.X/2) - 400
	firstTableFrame := s.appBody.Child(1).AsTree().Child(0).AsTree().This.(*core.Frame)
	stage.Pos.Y = int(firstTableFrame.Geom.Size.Actual.Total.Y/2) - 100

	stage.Run()
}

func (s *CategoryWindow) close() {
	s.catDialog.Close()
}

func getMainCategories(categoryArr [][]domain.Category) []string {
	mainCategories := make([]string, 0)

	for i := range categoryArr {
		mainCategories = append(mainCategories, categoryArr[i][0].MainCategory)
	}

	return mainCategories
}
