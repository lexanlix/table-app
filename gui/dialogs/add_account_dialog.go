package dialogs

import (
	"context"
	"strconv"

	"table-app/domain"
	"table-app/gui/iface"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
)

type AddAccountDialog struct {
	logger        log.Logger
	appBody       *core.Body
	dialogBody    *core.Body
	accountCtrl   iface.AccountController
	newAccount    domain.Account
	updateSumChan chan struct{}
}

func NewAddAccountDialog(
	logger log.Logger,
	appBody *core.Body,
	accountCtrl iface.AccountController,
	updateSumChan chan struct{},
) *AddAccountDialog {
	dialogBody := core.NewBody("AddAccount").
		SetTitle("Добавление счета")
	dialogBody.Styler(func(s *styles.Style) {
		s.Align.Self = styles.Center
		s.Min.X.Dp(710)
		s.Min.Y.Dp(300)
		s.CenterAll()
	})

	mainDialogFrame := core.NewFrame(dialogBody)
	mainDialogFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	titleFrame := core.NewFrame(mainDialogFrame)
	titleFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(300)
		s.CenterAll()
	})
	core.NewText(titleFrame).
		SetType(core.TextHeadlineSmall).
		SetText("Новый счет")

	core.NewSpace(mainDialogFrame).Styler(func(s *styles.Style) {
		s.Min.X.Dp(20)
	})

	dialog := &AddAccountDialog{
		logger:      logger,
		appBody:     appBody,
		dialogBody:  dialogBody,
		accountCtrl: accountCtrl,
		newAccount: domain.Account{
			IsInSum: true,
		},
		updateSumChan: updateSumChan,
	}

	dialog.addInputFields(mainDialogFrame)
	return dialog
}

func (s *AddAccountDialog) addInputFields(mainAccFrame *core.Frame) {
	inputMainFrame := core.NewFrame(mainAccFrame)
	inputMainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Gap.Zero()
	})

	// Поле ввода данных нового счета
	s.drawAddAccountFields(inputMainFrame)

	core.NewSpace(inputMainFrame).Styler(func(s *styles.Style) {
		s.Min.X.Dp(30)
	})

	// Область кнопок "Назад" и "Добавить счет"
	buttonsFrame := core.NewFrame(inputMainFrame)
	buttonsFrame.SetName("buttonsFrame")

	cancelButton := core.NewButton(buttonsFrame).
		SetType(core.ButtonElevated).
		SetText("Назад")
	cancelButton.OnClick(func(e events.Event) {
		s.close()
	})

	core.NewStretch(buttonsFrame)

	addButton := core.NewButton(buttonsFrame).
		SetType(core.ButtonFilled).
		SetText("Добавить счет")

	addButton.OnClick(func(e events.Event) {
		err := s.newAccount.Validate()
		if err != nil {
			s.logger.Error(context.Background(), "validate new account: "+err.Error())
			core.MessageSnackbar(s.dialogBody, "Необходимо ввести название счета")
			return
		}

		err = s.accountCtrl.AddAccount(context.Background(), s.newAccount)
		if err != nil {
			s.logger.Error(context.Background(), "add new account error: "+err.Error())
			core.ErrorSnackbar(s.dialogBody, err, "Ошибка добавления счета")
			return
		}

		s.appBody.Update()
		s.updateSumChan <- struct{}{}
		s.close()
	})
}

func (s *AddAccountDialog) drawAddAccountFields(inputMainFrame *core.Frame) {
	addAccountFrame := core.NewFrame(inputMainFrame)
	addAccountFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Min.X.Dp(700)
	})

	// Название
	addNameFrame := core.NewFrame(addAccountFrame)
	addNameFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})
	nameTextFrame := core.NewFrame(addNameFrame)
	nameTextFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(160)
		s.CenterAll()
	})
	core.NewText(nameTextFrame).SetText("Название")

	nameTextField := core.NewTextField(addNameFrame)
	nameTextField.Styler(func(s *styles.Style) {
		s.Min.X.Dp(500)
	})
	nameTextField.OnChange(func(e events.Event) {
		s.newAccount.Name = nameTextField.Text()
	})

	// Сумма
	addSumFrame := core.NewFrame(addAccountFrame)
	addSumFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})
	sumTextFrame := core.NewFrame(addSumFrame)
	sumTextFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(160)
		s.CenterAll()
	})
	core.NewText(sumTextFrame).SetText("Сумма")

	sumTextField := core.NewTextField(addSumFrame).SetPlaceholder("0")
	sumTextField.Styler(func(s *styles.Style) {
		s.Min.X.Dp(500)
	})
	sumTextField.OnChange(func(e events.Event) {
		sum, err := strconv.Atoi(sumTextField.Text())
		if err != nil {
			s.logger.Error(context.Background(), "convert account sum to int error: "+err.Error())
			core.MessageSnackbar(s.dialogBody, "Неверный формат суммы: "+err.Error())
			return
		}

		s.newAccount.Sum = sum
	})

	// Комментарий
	addNoteFrame := core.NewFrame(addAccountFrame)
	addNoteFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})
	noteTextFrame := core.NewFrame(addNoteFrame)
	noteTextFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(160)
		s.CenterAll()
	})
	core.NewText(noteTextFrame).SetText("Комментарий")

	noteTextField := core.NewTextField(addNoteFrame)
	noteTextField.Styler(func(s *styles.Style) {
		s.Min.X.Dp(500)
	})
	noteTextField.OnChange(func(e events.Event) {
		s.newAccount.Note = noteTextField.Text()
	})

	// Учитывать ли в сумме
	addIsInSumFrame := core.NewFrame(addAccountFrame)
	addIsInSumFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})
	isInSumTextFrame := core.NewFrame(addIsInSumFrame)
	isInSumTextFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(160)
		s.CenterAll()
	})
	core.NewText(isInSumTextFrame).SetText("Учитывать в сумме")

	isInSumSwitch := core.NewSwitch(addIsInSumFrame).
		SetType(core.SwitchCheckbox).
		SetChecked(true)
	isInSumSwitch.OnChange(func(e events.Event) {
		s.newAccount.IsInSum = isInSumSwitch.StateIs(states.Checked)
	})
	isInSumSwitch.SetTooltip(isInSumTooltip)
}

func (s *AddAccountDialog) Run() {
	stage := s.dialogBody.NewDialog(s.appBody)

	stage.Pos.X = int(s.appBody.Geom.Size.Actual.Total.X/2) - 700
	firstTableFrame := s.appBody.Child(1).AsTree().Child(0).AsTree().This.(*core.Frame)
	stage.Pos.Y = int(firstTableFrame.Geom.Size.Actual.Total.Y/2) + 100
	stage.Run()
}

func (s *AddAccountDialog) close() {
	s.dialogBody.Close()
}
