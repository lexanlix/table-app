package dialogs

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

type CategoryDialog struct {
	logger         log.Logger
	appBody        *core.Body
	dialogBody     *core.Body
	controller     iface.TableController
	mainCategories []string
	category       domain.Category
}

func NewCategoryDialog(
	logger log.Logger,
	appBody *core.Body,
	controller iface.TableController,
	categories [][]domain.Category,
) *CategoryDialog {
	dialogBody := core.NewBody("NewCategory").
		SetTitle("Добавление категории")
	dialogBody.Styler(func(s *styles.Style) {
		s.Align.Self = styles.Center
		s.CenterAll()
	})

	mainDialogFrame := core.NewFrame(dialogBody)
	mainDialogFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	inputFrame := core.NewFrame(mainDialogFrame)

	core.NewSpace(mainDialogFrame).Styler(func(s *styles.Style) {
		s.Min.Y.Dp(10)
	})

	buttonsFrame := core.NewFrame(mainDialogFrame)

	dialog := &CategoryDialog{
		logger:         logger,
		appBody:        appBody,
		dialogBody:     dialogBody,
		controller:     controller,
		mainCategories: getMainCategories(categories),
		category:       domain.Category{},
	}

	dialog.addCategoryInput(inputFrame)
	dialog.addControlButtons(buttonsFrame)

	return dialog
}

func (s *CategoryDialog) addCategoryInput(inputFrame *core.Frame) {
	inputFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	// Область заголовка диалога
	titleFrame := core.NewFrame(inputFrame)
	titleFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(400)
		s.CenterAll()
	})
	core.NewText(titleFrame).
		SetType(core.TextDisplaySmall).
		SetText("Добавить категорию")

	core.NewSpace(inputFrame).Styler(func(s *styles.Style) {
		s.Min.Y.Dp(10)
	})

	// Область выбора главной категории
	s.addNewCategoryFields(inputFrame)
}

func (s *CategoryDialog) addNewCategoryFields(inputFrame *core.Frame) {
	// Область выбора главной категории
	choiceFrame := core.NewFrame(inputFrame)
	choiceFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Min.X.Dp(300)
	})
	core.NewText(choiceFrame).
		SetType(core.TextBodyLarge).
		SetText("Выбрать главную категорию")

	switcher := core.NewSwitches(choiceFrame).
		SetType(core.SwitchChip).
		SetMutex(true).
		SetStrings(s.mainCategories...)

	switcher.Styler(func(s *styles.Style) {
		s.Font.Size.Set(8, units.UnitPt)
	})

	switcher.OnChange(func(e events.Event) {
		if switcher.SelectedItem() != nil {
			s.category.MainCategory = switcher.SelectedItem().Value.(string)
		}
	})

	core.NewSpace(inputFrame).Styler(func(s *styles.Style) {
		s.Min.Y.Dp(10)
	})

	// Область ввода названия категории
	addNameFrame := core.NewFrame(inputFrame)
	addNameFrame.SetName("inputFrame")
	addNameFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	core.NewText(addNameFrame).
		SetType(core.TextBodyLarge).
		SetText("Название категории")

	nameTextField := core.NewTextField(addNameFrame).
		SetType(core.TextFieldOutlined)
	nameTextField.Styler(func(s *styles.Style) {
		s.Min.X.Dp(300)
		s.Font.Size.Set(8, units.UnitPt)
	})

	nameTextField.OnInput(func(e events.Event) {
		s.category.Name = nameTextField.Text()
	})
}

func (s *CategoryDialog) addControlButtons(buttonsFrame *core.Frame) {
	buttonsFrame.Styler(func(s *styles.Style) {
		s.Max.X.Dp(400)
		s.CenterAll()
	})

	cancelButton := core.NewButton(buttonsFrame).
		SetType(core.ButtonElevated).
		SetText("Отмена")
	cancelButton.OnClick(func(e events.Event) {
		s.close()
	})

	core.NewStretch(buttonsFrame)

	addButton := core.NewButton(buttonsFrame).
		SetType(core.ButtonFilled).
		SetText("Добавить")
	addButton.OnClick(func(e events.Event) {
		if len(s.category.MainCategory) == 0 || len(s.category.Name) == 0 {
			core.MessageSnackbar(s.dialogBody, "Введите данные для добавления категории")
			return
		}

		if s.controller.CategoryIsExist(context.Background(), s.category) {
			core.MessageSnackbar(s.dialogBody, "Категория уже существует")
			return
		}

		err := s.controller.AddCategory(context.Background(), s.category)
		if err != nil {
			core.MessageSnackbar(s.dialogBody, "Ошибка добавления категории: "+err.Error())
			s.logger.Error(context.Background(), "add category error", log.Any("err", err.Error()))
			return
		}

		s.appBody.Update()
		s.close()
	})
}

func (s *CategoryDialog) Run() {
	stage := s.dialogBody.NewDialog(s.appBody)

	stage.Pos.X = int(s.appBody.Geom.Size.Actual.Total.X/2) - 400
	firstTableFrame := s.appBody.Child(1).AsTree().Child(0).AsTree().This.(*core.Frame)
	stage.Pos.Y = int(firstTableFrame.Geom.Size.Actual.Total.Y/2) - 100

	stage.Run()
}

func (s *CategoryDialog) close() {
	s.dialogBody.Close()
}

func getMainCategories(categoryArr [][]domain.Category) []string {
	mainCategories := make([]string, 0)

	for i := range categoryArr {
		mainCategories = append(mainCategories, categoryArr[i][0].MainCategory)
	}

	return mainCategories
}
