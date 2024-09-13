package gui

import (
	"context"
	"strconv"

	"table-app/domain"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
	"cogentcore.org/core/styles"
)

type SumWindow struct {
	logger    log.Logger
	mainFrame *core.Frame
	sumDialog *core.Body
	sumFrame  *core.Frame
	tFields   []*core.TextField

	controller TableController
	cell       domain.Cell
	sum        int
	updateChan chan domain.Cell
}

func NewSumWindow(logger log.Logger, mainFrame *core.Frame, cell domain.Cell, controller TableController,
	updateChan chan domain.Cell) *SumWindow {
	sumBody := core.NewBody("Sum").SetTitle(cell.Category)
	sumBody.Styler(func(s *styles.Style) {
		s.Align.Self = styles.Center
		s.CenterAll()
	})

	mainSumFrame := core.NewFrame(sumBody)
	mainSumFrame.SetName("mainSumFrame")
	mainSumFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Row
		s.CenterAll()
	})

	sumFrame := core.NewFrame(mainSumFrame)
	sumFrame.SetName("sumFrame")
	sumFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})
	firstTField := core.NewTextField(sumFrame)
	secondTField := core.NewTextField(sumFrame)

	tFields := []*core.TextField{firstTField, secondTField}

	rightFrame := core.NewFrame(mainSumFrame)
	rightFrame.SetName("rightFrame")
	rightFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	textSumFrame := core.NewFrame(rightFrame)
	textSumFrame.SetName("textSumFrame")

	_ = core.NewSeparator(rightFrame)

	buttonsFrame := core.NewFrame(rightFrame)
	buttonsFrame.SetName("buttonsFrame")

	sumWindow := &SumWindow{
		logger:     logger,
		mainFrame:  mainFrame,
		sumDialog:  sumBody,
		sumFrame:   sumFrame,
		controller: controller,
		cell:       cell,
		tFields:    tFields,
		updateChan: updateChan,
		sum:        0,
	}
	sumWindow.addButtons(buttonsFrame)
	textSum := sumWindow.addTextSum(textSumFrame)

	sumWindow.addInputFunction(textSum, tFields...)

	return sumWindow
}

func (s *SumWindow) addInputFunction(textSum *core.Text, tFields ...*core.TextField) {
	for i := range tFields {
		tFields[i].Styler(func(s *styles.Style) {
			s.Min.X.Dp(120)
		})

		tFields[i].OnChange(func(e events.Event) {
			val, err := strconv.Atoi(tFields[i].Text())
			if err != nil {
				core.MessageSnackbar(s.sumDialog, "Неверный формат данных: "+err.Error())
				return
			}

			s.cell.Value += val
			s.sum += val
			textSum.Update()
			tFields[i].SetText(FormatInt(val))

			newTField := core.NewTextField(s.sumFrame)
			s.addInputFunction(textSum, newTField)
			s.sumFrame.Update()
		})
	}
}

func (s *SumWindow) addButtons(buttonsFrame *core.Frame) {
	buttonsFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	saveButton := core.NewButton(buttonsFrame).SetType(core.ButtonFilled).SetText("Сохранить")
	saveButton.OnClick(func(e events.Event) {
		err := s.controller.UpsertValue(context.Background(), s.cell)
		if err != nil {
			core.MessageSnackbar(s.sumDialog, "Ошибка сохранения данных: "+err.Error())
			s.logger.Error(context.Background(), "save all data error", log.Any("err", err.Error()))
			return
		}

		if s.updateChan != nil {
			s.updateChan <- s.cell
		}

		s.close()

		core.MessageSnackbar(s.mainFrame, "Введено: "+strconv.Itoa(s.sum))
	})

	cancelButton := core.NewButton(buttonsFrame).SetType(core.ButtonElevated).SetText("Отмена")
	cancelButton.OnClick(func(e events.Event) {
		s.close()
	})
}

func (s *SumWindow) addTextSum(textSumFrame *core.Frame) *core.Text {
	textSumFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	textSummaFrame := core.NewFrame(textSumFrame)
	textSummaFrame.SetName("textSummaFrame")
	textSummaFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(150)
		s.CenterAll()
	})
	_ = core.NewText(textSummaFrame).
		SetType(core.TextHeadlineSmall).
		SetText("Сумма:")

	textValueFrame := core.NewFrame(textSumFrame)
	textValue := core.NewText(textValueFrame).
		SetType(core.TextHeadlineSmall).
		SetText(FormatInt(s.sum))
	core.Bind(&s.sum, textValue)
	return textValue
}

func (s *SumWindow) Run(ctx core.Widget) {
	stage := s.sumDialog.NewDialog(ctx)
	stage.Run()
}

func (s *SumWindow) close() {
	s.sumDialog.Close()
}

func (s *SumWindow) isCommandEnterChord(chord key.Chord) bool {
	if !chord.PlatformChord().IsMulti() {
		return false
	}

	_, codes, _, err := chord.Decode()
	if err != nil {
		s.logger.Error(context.Background(), "chord.Decode() error", log.Any("err", err.Error()))
		return false
	}

	keys := codes.Values()

	if len(keys) != 2 {
		return false
	}

	if keys[0] == key.CodeLeftMeta || keys[0] == key.CodeRightMeta {
		if keys[1] == key.CodeReturnEnter {
			return true
		}
	}

	return false
}
