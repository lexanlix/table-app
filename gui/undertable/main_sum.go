package undertable

import (
	"context"

	"table-app/domain"
	"table-app/gui/styles/format"

	"cogentcore.org/core/core"
	"cogentcore.org/core/styles"
)

const (
	mainSumTooltip = "Сумма по всем счетам, за исключением не учитываемых"
	diffSumTooltip = "Считается как 'Общий бюджет минус остаток в последнем месяце'"

	textSumFrameMinDpX  = 150
	valueSumFrameMinDpX = 100
	sumFrameMinDpY      = 30
)

// drawSum Отрисовка информации об общем бюджете (mainSum) и разнице (diffSum)
// - Первая строка: Общий бюджет, содержит сумму по всем счетам
// - Вторая строка: Разница - разность между общим бюджетом и последним остатком на последний месяц
func (t *UnderTable) drawSum() {
	sums, err := t.accountController.GetSum(context.Background())
	if err != nil {
		t.logger.Error(context.Background(), "get sum error: "+err.Error())
		return
	}

	mainFrame := core.NewFrame(t.upperFrame)
	mainFrame.SetName("mainSumFrame")
	mainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	mainSumText := t.drawMainSum(mainFrame, sums)
	diffSumText := t.drawDiffSum(mainFrame, sums)

	t.sumUpdater.AddSumTexts(mainSumText, diffSumText)
}

func (t *UnderTable) drawMainSum(mainFrame *core.Frame, sums domain.Sum) *core.Text {
	mainSumFrame := core.NewFrame(mainFrame)
	mainSumFrame.Styler(func(s *styles.Style) {
		s.Gap.Zero()
	})
	mainSumFrame.SetName("textMainSumFrame")

	textMainSumFrame := core.NewFrame(mainSumFrame)
	textMainSumFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(textSumFrameMinDpX)
		s.Min.Y.Dp(sumFrameMinDpY)
		s.CenterAll()
		s.Gap.Zero()
	})
	core.NewText(textMainSumFrame).SetText("Общий бюджет:").SetTooltip(mainSumTooltip)

	valueMainSumFrame := core.NewFrame(mainSumFrame)
	valueMainSumFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(valueSumFrameMinDpX)
		s.Min.Y.Dp(sumFrameMinDpY)
		s.CenterAll()
	})
	mainSumText := core.NewText(valueMainSumFrame).SetText(format.FormatInt(sums.MainSum))
	mainSumText.SetTooltip(mainSumTooltip)

	return mainSumText
}

func (t *UnderTable) drawDiffSum(mainFrame *core.Frame, sums domain.Sum) *core.Text {
	mainDiffFrame := core.NewFrame(mainFrame)
	mainDiffFrame.Styler(func(s *styles.Style) {
		s.Gap.Zero()
	})
	mainDiffFrame.SetName("mainDiffFrame")

	textMainDiffFrame := core.NewFrame(mainDiffFrame)
	textMainDiffFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(textSumFrameMinDpX)
		s.Min.Y.Dp(sumFrameMinDpY)
		s.CenterAll()
	})
	core.NewText(textMainDiffFrame).SetText("Разница:").SetTooltip(diffSumTooltip)

	valueMainDiffFrame := core.NewFrame(mainDiffFrame)
	valueMainDiffFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(valueSumFrameMinDpX)
		s.Min.Y.Dp(sumFrameMinDpY)
		s.CenterAll()
	})
	diffSumText := core.NewText(valueMainDiffFrame).SetText(format.FormatInt(sums.DiffSum))
	diffSumText.SetTooltip(diffSumTooltip)

	return diffSumText
}
