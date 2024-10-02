package table

import (
	"strconv"
	"time"

	"table-app/domain"
	"table-app/gui/styles/colors"

	"cogentcore.org/core/core"
	"cogentcore.org/core/styles"
)

func (t Table) drawMonthsColumn() *core.Frame {
	mainFrame := core.NewFrame(t.tableFrame)
	mainFrame.SetName("monthsFrame")
	mainFrame.Styler(t.styler.MonthFrameStyle())

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
			s.Min.X.Dp(t.settings.Gui.CellSizeDpX)
			s.Min.Y.Dp(t.settings.Gui.CellSizeDpY)
			s.Gap.Zero()
			s.Border.Width.Top.Dp(1)
			s.Border.Width.Bottom.Dp(1)
			s.CenterAll()

			if t.year == time.Now().Year() && month == int(lastMonth) {
				s.Background = colors.ColorYellow
			}
		})

		monthName, ok := domain.RusMonths[month]
		if !ok {
			monthName = time.Month(month).String()
		}

		core.NewText(monthFrame).SetText(monthName)
	}

	if t.year == time.Now().Year() {
		return mainFrame
	}

	// строка итогов года
	resultFrame := core.NewFrame(mainFrame)
	resultFrame.SetName("resultFrame")
	resultFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(t.settings.Gui.CellSizeDpX)
		s.Min.Y.Dp(t.settings.Gui.CellSizeDpY)
		s.Gap.Zero()
		s.Border.Width.Top.Dp(1)
		s.Border.Width.Bottom.Dp(1)
		s.CenterAll()
		s.Background = colors.ColorPurple
	})

	resultText := core.NewText(resultFrame).SetText(strconv.Itoa(t.year))
	resultText.Styler(func(s *styles.Style) {
		s.Font.Weight = styles.WeightBold
	})

	return mainFrame
}
