package undertable

import (
	"context"
	"strconv"

	"table-app/gui/dialogs"
	"table-app/gui/styles/format"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

// drawAccounts Отрисовка информации о счетах
// строка с ячейками, содержащими название счета и сумму на счете
func (t *UnderTable) drawAccounts() {
	mainAccountsFrame := core.NewFrame(t.underTableFrame)
	mainAccountsFrame.SetName("mainAccountsFrame")
	mainAccountsFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	textHeadAccountsFrame := core.NewFrame(mainAccountsFrame)
	textHeadAccountsFrame.SetName("textHeadAccountsFrame")
	core.NewText(textHeadAccountsFrame).SetText("Счета")

	accountListFrame := core.NewFrame(mainAccountsFrame)
	accountListFrame.SetName("accountFrame")
	accountListFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})

	accountListFrame.Maker(func(p *tree.Plan) {
		for i, account := range *t.accountList {
			if (*t.accountList)[i].Deleted {
				continue
			}

			tree.AddAt(p, "account_"+(*t.accountList)[i].Name+"_Frame", func(accountFrame *core.Frame) {
				accountFrame.Styler(func(s *styles.Style) {
					s.Direction = styles.Column
					s.Border.Width.SetAll(units.Dp(1))
					s.Gap.Zero()
				})

				accountFieldName := core.NewTextField(accountFrame)
				accountFieldName.SetName("accountFieldName")

				accountFieldName.Styler(t.textFieldStyler.AccountName(&(*t.accountList)[i].IsInSum))
				accountFieldName.SetText((*t.accountList)[i].Name)
				accountFieldName.SetTooltip((*t.accountList)[i].Note)

				accountFieldName.OnChange(func(e events.Event) {
					account.Name = accountFieldName.Text()
					err := t.accountController.UpdateAccount(context.Background(), account)
					if err != nil {
						t.logger.Error(context.Background(), "update account error: "+err.Error())
						core.MessageSnackbar(t.underTableFrame, "Ошибка обновления данных счета: "+err.Error())
					}

					t.updateSumChan <- struct{}{}
				})

				accountFieldSum := core.NewTextField(accountFrame)
				accountFieldSum.SetName("accountFieldSum")

				accountFieldSum.Styler(t.textFieldStyler.AccountSum())
				accountFieldSum.SetText(format.FormatInt((*t.accountList)[i].Sum))

				accountFieldSum.OnChange(func(e events.Event) {
					sum, err := strconv.Atoi(accountFieldSum.Text())
					if err != nil {
						core.MessageSnackbar(t.underTableFrame,
							"Неверный формат введенной суммы на счете "+(*t.accountList)[i].Name)
						return
					}

					account.Sum = sum

					err = t.accountController.UpdateAccount(context.Background(), account)
					if err != nil {
						t.logger.Error(context.Background(), "update account error: "+err.Error())
						core.MessageSnackbar(t.underTableFrame, "Ошибка обновления данных счета: "+err.Error())
					}

					accountFieldSum.SetText(format.FormatInt(account.Sum))
					t.updateSumChan <- struct{}{}
				})
			})
		}

		tree.AddAt(p, "addAccountButton", func(addAccButton *core.Button) {
			addAccButton.SetIcon(icons.Add).SetType(core.ButtonTonal)
			addAccButton.SetTooltip("Добавить счет")

			addAccButton.OnClick(func(e events.Event) {
				if e.MouseButton() == events.Left {
					addAccountDialog := dialogs.NewAddAccountDialog(t.logger, t.appBody, t.accountController, t.updateSumChan)
					addAccountDialog.Run()
				}
			})
		})
	})
}
