package styles

import (
	"image"

	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
)

const (
	yearFrameMinYDp       = 83
	tableHeadTextFrameYDp = 81
)

type FrameStyle struct {
	sizeDpX float32
	sizeDpY float32
}

func NewFrameStyle(sizeDpX float32, sizeDpY float32) *FrameStyle {
	return &FrameStyle{
		sizeDpX: sizeDpX,
		sizeDpY: sizeDpY,
	}
}

// TableFrameStyle стиль таблицы - 4 фрейма сеткой 2 на 2
func (f *FrameStyle) TableFrameStyle() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Display = styles.Grid
		s.Columns = 2
		s.Border.Width.Set(units.Dp(1))
		s.Border.Radius.Zero()
		s.Gap.Zero()
		s.CenterAll()
	}
}

func (f *FrameStyle) YearFrameStyle() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Min.X.Dp(f.sizeDpX)
		s.Min.Y.Dp(yearFrameMinYDp)
		s.Border.Width.Right.Dp(1)
		s.Border.Width.Bottom.Dp(2)
		s.Gap.Zero()
		s.CenterAll()
	}
}

func (f *FrameStyle) HeadFrameStyle() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Border.Width.Left.Dp(1)
		s.Border.Width.Bottom.Dp(1)
		s.Gap.Zero()
		s.CenterAll()
	}
}

func (f *FrameStyle) MonthFrameStyle() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Direction = styles.Column
		s.Border.Width.Right.Dp(1)
		s.Border.Width.Top.Dp(1)
		s.Gap.Zero()
	}
}

func (f *FrameStyle) ValuesFrameStyle() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Direction = styles.Column
		s.Border.Width.Left.Dp(1)
		s.Border.Width.Top.Dp(1)
		s.Border.Width.Right.Dp(0.5)
		s.Gap.Zero()
	}
}

// StandardSizeWithAllBorders
// Фрейм с заданными размерами для ячейки и всеми границами, можно указать цвет фона
func (f *FrameStyle) StandardSizeWithAllBorders(color ...image.Image) func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Min.X.Dp(f.sizeDpX)
		s.Min.Y.Dp(f.sizeDpY)
		s.Border.Width.SetAll(units.Dp(1))
		s.CenterAll()

		if len(color) != 0 {
			s.Background = color[0]
		}
	}
}

// TableHeadTextFrame
// Фрейм шапки со специально выбранным sizeY, стандартным sizeX и всеми границами
func (f *FrameStyle) TableHeadTextFrame() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Min.X.Dp(f.sizeDpX)
		s.Min.Y.Dp(tableHeadTextFrameYDp)
		s.Border.Width.SetAll(units.Dp(1))
		s.CenterAll()
	}
}
