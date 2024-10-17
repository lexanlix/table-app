package styles

import (
	custom "table-app/gui/styles/colors"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/styles"
)

const (
	tFieldAccountXMinDp = 200
	tFieldAccountYMaxDp = 15
)

type TextFieldStyle struct {
}

func NewTextFieldStyle() *TextFieldStyle {
	return &TextFieldStyle{}
}

// AccountName стиль текстового поля названия счета
func (tf *TextFieldStyle) AccountName(isInSum *bool) func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Min.X.Dp(tFieldAccountXMinDp)
		s.Max.Y.Dp(tFieldAccountYMaxDp)
		s.Border.Radius.Zero()
		s.Border.Width.Zero()
		s.Border.Offset.Zero()
		s.Background = colors.Uniform(colors.White)
		if !(*isInSum) {
			s.Background = custom.ColorWhiteRed
		}
	}
}

// AccountSum стиль текстового поля суммы счета
func (tf *TextFieldStyle) AccountSum() func(s *styles.Style) {
	return func(s *styles.Style) {
		s.Min.X.Dp(tFieldAccountXMinDp)
		s.Max.Y.Dp(tFieldAccountYMaxDp)
		s.Border.Radius.Zero()
		s.Border.Width.Zero()
		s.Border.Offset.Zero()
	}
}
