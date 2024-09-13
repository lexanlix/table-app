package gui

import (
	"cogentcore.org/core/core"
	"cogentcore.org/core/tree"
)

func (a *App) withFrame(parent tree.Node) *core.Frame {
	frame := core.NewFrame(parent)
	a.frames = append(a.frames, frame)

	return frame
}

func (a *App) withTextField(parent tree.Node) *core.TextField {
	tField := core.NewTextField(parent)
	a.txtFields = append(a.txtFields, tField)

	return tField
}

func (a *App) withText(parent tree.Node, text string) *core.Text {
	t := core.NewText(parent).SetText(text)
	a.texts = append(a.texts, t)

	return t
}
