package dialogs

import (
	"context"
	"strconv"

	"table-app/domain"
	"table-app/gui/iface"
	custom "table-app/gui/styles/colors"
	"table-app/gui/styles/format"
	"table-app/internal/log"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
	"github.com/pkg/errors"
)

const isInSumTooltip = "Добавлять сумму данного счета к общей сумме по всем счетам?"

type AccountsDialog struct {
	logger        log.Logger
	appBody       *core.Body
	dialogBody    *core.Body
	listFrame     *core.Frame
	accountCtrl   iface.AccountController
	newAccount    domain.Account
	accountList   []domain.Account
	updateSumChan chan struct{}
}

func NewAccountsDialog(
	logger log.Logger,
	appBody *core.Body,
	accountCtrl iface.AccountController,
	updateSumChan chan struct{},
) (*AccountsDialog, error) {
	accountList, err := accountCtrl.GetAll(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "get account list")
	}

	dialogBody := core.NewBody("Accounts").SetTitle("Редактирование счетов")
	dialogBody.Styler(func(s *styles.Style) {
		s.Align.Self = styles.Center
		s.Min.X.Dp(800)
		s.Min.Y.Dp(600)
		s.CenterAll()
	})

	mainAccFrame := core.NewFrame(dialogBody)
	mainAccFrame.SetName("mainAccFrame")
	mainAccFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	titleFrame := core.NewFrame(mainAccFrame)
	titleFrame.SetName("titleFrame")
	titleFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(400)
		s.CenterAll()
	})
	core.NewText(titleFrame).SetType(core.TextHeadlineMedium).SetText("Ваши счета")

	accDialog := &AccountsDialog{
		logger:      logger,
		appBody:     appBody,
		dialogBody:  dialogBody,
		accountCtrl: accountCtrl,
		newAccount: domain.Account{
			IsInSum: true,
		},
		accountList:   accountList,
		updateSumChan: updateSumChan,
	}

	accDialog.addAccountsList(mainAccFrame)
	accDialog.addNewAccFields(mainAccFrame)

	return accDialog, nil
}

func (s *AccountsDialog) addAccountsList(mainAccFrame *core.Frame) {
	s.listFrame = core.NewFrame(mainAccFrame)
	s.listFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Min.X.Dp(800)
		s.Border.Width.SetAll(units.Dp(1))
		s.Gap.Zero()
	})

	s.listFrame.Maker(func(p *tree.Plan) {
		tree.AddAt(p, "head", func(accountFrame *core.Frame) {
			s.drawHeadRow(accountFrame)
		})

		for i := range s.accountList {
			s.drawAccountRow(p, i)
		}
	})
}

func (s *AccountsDialog) drawHeadRow(accountFrame *core.Frame) {
	accountFrame.Styler(func(s *styles.Style) {
		s.Align.Content = styles.Center
		s.Border.Width.SetAll(units.Dp(1))
		s.Min.X.Dp(1300)
		s.Gap.Zero()
		s.Gap.X.Dp(2)
	})

	frameName := core.NewFrame(accountFrame)
	frameName.Styler(func(s *styles.Style) {
		s.Min.X.Dp(300)
		s.Min.Y.Dp(30)
		s.CenterAll()
	})
	core.NewText(frameName).SetText("Название")

	frameSum := core.NewFrame(accountFrame)
	frameSum.Styler(func(s *styles.Style) {
		s.Min.X.Dp(150)
		s.Min.Y.Dp(30)
		s.CenterAll()
	})
	core.NewText(frameSum).SetText("Сумма")

	frameComment := core.NewFrame(accountFrame)
	frameComment.Styler(func(s *styles.Style) {
		s.Min.X.Dp(510)
		s.Min.Y.Dp(30)
		s.CenterAll()
	})
	core.NewText(frameComment).SetText("Комментарий")

	frameIsInSum := core.NewFrame(accountFrame)
	frameIsInSum.Styler(func(s *styles.Style) {
		s.Min.X.Dp(180)
		s.Min.Y.Dp(30)
		s.CenterAll()
	})
	core.NewText(frameIsInSum).SetText("Учитывать в сумме")
}

func (s *AccountsDialog) drawAccountRow(p *tree.Plan, idx int) {
	if s.accountList[idx].Deleted {
		return
	}

	tree.AddAt(p, "account+"+strconv.Itoa(idx), func(accountFrame *core.Frame) {
		accountFrame.Styler(func(s *styles.Style) {
			s.Align.Content = styles.Center
			s.Min.X.Dp(1300)
			s.Gap.Zero()
			s.Gap.X.Dp(2)
			s.CenterAll()
			s.Border.Width.SetAll(units.Dp(1))
		})

		// Поле названия счета
		fieldName := core.NewTextField(accountFrame)
		fieldName.SetText(s.accountList[idx].Name)
		fieldName.Styler(func(s *styles.Style) {
			s.Min.X.Dp(200)
			s.Max.Y.Dp(15)
		})
		fieldName.OnChange(func(e events.Event) {
			name := fieldName.Text()
			if len(name) == 0 {
				core.MessageSnackbar(s.dialogBody, "Название счета не может быть пустым")
				return
			}
			s.accountList[idx].Name = name
		})

		// Поле суммы счета
		fieldSum := core.NewTextField(accountFrame)
		fieldSum.Styler(func(s *styles.Style) {
			s.Min.X.Dp(70)
			s.Max.Y.Dp(15)
		})
		fieldSum.SetText(format.FormatInt(s.accountList[idx].Sum))
		fieldSum.OnChange(func(e events.Event) {
			sum, err := strconv.Atoi(fieldSum.Text())
			if err != nil {
				s.logger.Error(context.Background(), "convert account sum to int error: "+err.Error())
				core.MessageSnackbar(s.dialogBody, "Неверный формат суммы: "+err.Error())
				return
			}

			s.accountList[idx].Sum = sum
			fieldSum.SetText(format.FormatInt(sum))
		})

		// Поле комментария счета
		fieldComment := core.NewTextField(accountFrame)
		fieldComment.Styler(func(s *styles.Style) {
			s.Min.X.Dp(500)
			s.Max.Y.Dp(15)
		})
		fieldComment.SetText(s.accountList[idx].Note)
		fieldComment.OnChange(func(e events.Event) {
			s.accountList[idx].Note = fieldComment.Text()
		})

		core.NewSpace(accountFrame).Styler(func(s *styles.Style) {
			s.Min.X.Dp(90)
		})

		// Поле учитывать ли в общей сумме
		isInSumSwitch := core.NewSwitch(accountFrame).
			SetType(core.SwitchCheckbox).
			SetChecked(s.accountList[idx].IsInSum)
		isInSumSwitch.OnChange(func(e events.Event) {
			s.accountList[idx].IsInSum = isInSumSwitch.StateIs(states.Checked)
		})
		isInSumSwitch.SetTooltip(isInSumTooltip)

		core.NewStretch(accountFrame)

		// Поле удаления счета
		deleteButton := core.NewButton(accountFrame).SetType(core.ButtonOutlined).SetText("Удалить")
		deleteButton.Styler(func(s *styles.Style) {
			s.Background = colors.Uniform(colors.White)
			s.Color = custom.ColorDeleteRed
		})
		deleteButton.OnClick(func(e events.Event) {
			s.accountList[idx].Deleted = true
			core.MessageSnackbar(s.dialogBody, "Счет удален")
			s.listFrame.Update()
		})
	})
}

func (s *AccountsDialog) addNewAccFields(mainAccFrame *core.Frame) {
	bottomFrame := core.NewFrame(mainAccFrame)
	bottomFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})

	// Область с кнопками управления - "Назад" и "Сохранить и выйти"
	controlButtonsFrame := core.NewFrame(bottomFrame)
	controlButtonsFrame.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 1)
		s.Direction = styles.Column
		s.CenterAll()
	})

	cancelButton := core.NewButton(controlButtonsFrame).SetType(core.ButtonElevated).SetText("Назад")
	cancelButton.OnClick(func(e events.Event) {
		s.close()
	})

	core.NewSpace(controlButtonsFrame)

	saveButton := core.NewButton(controlButtonsFrame).SetType(core.ButtonElevated).SetText("Сохранить и выйти")
	saveButton.Styler(func(s *styles.Style) {
		s.Background = custom.ColorGreen
		s.Color = colors.Uniform(colors.White)
	})
	saveButton.OnClick(func(e events.Event) {
		err := s.accountCtrl.UpdateList(context.Background(), s.accountList)
		if err != nil {
			s.logger.Error(context.Background(), "update accounts list: "+err.Error())
			core.ErrorSnackbar(s.dialogBody, err, "Ошибка сохранения")
			return
		}

		s.appBody.Update()
		s.updateSumChan <- struct{}{}
		s.close()
	})

	// Область добавления нового счета
	addAccountMainFrame := core.NewFrame(bottomFrame)
	addAccountMainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Gap.Zero()
	})

	// Область ввода данных счета
	s.drawAddAccountFields(addAccountMainFrame)

	// Область кнопки "Добавить счет"
	addButtonFrame := core.NewFrame(addAccountMainFrame)
	addButtonFrame.SetName("buttonsFrame")

	addButton := core.NewButton(addButtonFrame).
		SetType(core.ButtonElevated).
		SetText("Добавить счет")
	addButton.Styler(func(s *styles.Style) {
		s.Color = colors.Uniform(colors.Black)
	})

	addButton.OnClick(func(e events.Event) {
		err := s.newAccount.Validate()
		if err != nil {
			s.logger.Error(context.Background(), "validate new account: "+err.Error())
			core.MessageSnackbar(s.dialogBody, "Необходимо ввести название счета")
			return
		}

		s.accountList = append(s.accountList, s.newAccount)
		s.listFrame.Update()
		s.newAccount = domain.Account{IsInSum: true}
	})
}

func (s *AccountsDialog) drawAddAccountFields(addAccountMainFrame *core.Frame) {
	addAccountFrame := core.NewFrame(addAccountMainFrame)
	addAccountFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Min.X.Dp(700)
		s.Border.Width.SetAll(units.Dp(1))
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

func (s *AccountsDialog) Run() {
	stage := s.dialogBody.NewDialog(s.appBody)

	stage.Pos.X = int(s.appBody.Geom.Size.Actual.Total.X/2) - 1050
	firstTableFrame := s.appBody.Child(1).AsTree().Child(0).AsTree().This.(*core.Frame)
	stage.Pos.Y = int(firstTableFrame.Geom.Size.Actual.Total.Y/2) - 200

	stage.Run()
}

func (s *AccountsDialog) close() {
	s.dialogBody.Close()
}
