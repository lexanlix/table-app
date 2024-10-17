package undertable

import (
	"context"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
)

// drawLastUpdating отрисовка информации о последнем обновлении таблицы:
// - строка с датой и временем последнего обновления;
// - строка с текстовой записью о последних учтенных данных
func (t *UnderTable) drawLastUpdating() {
	mainUpdatingFrame := core.NewFrame(t.upperFrame)
	mainUpdatingFrame.SetName("mainUpdatingFrame")
	mainUpdatingFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	t.drawLastUpdatingTime(mainUpdatingFrame)
	t.drawLastRecord(mainUpdatingFrame)
}

func (t *UnderTable) drawLastUpdatingTime(mainUpdatingFrame *core.Frame) {
	updatingFrame := core.NewFrame(mainUpdatingFrame)
	updatingFrame.SetName("updatingFrame")

	leftUpdatingFrame := core.NewFrame(updatingFrame)
	leftUpdatingFrame.SetName("leftUpdatingFrame")
	leftUpdatingFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(100)
	})
	core.NewText(leftUpdatingFrame).SetText("Обновлено: ")

	updatedTime := t.updatingController.GetLastUpdated()

	rightUpdatingFrame := core.NewFrame(updatingFrame)
	rightUpdatingFrame.SetName("rightUpdatingFrame")
	rightUpdatingFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(300)
	})

	updatingText := core.NewText(rightUpdatingFrame)
	t.sumUpdater.AddUpdatedTimeText(updatingText)
	core.Bind(updatedTime, updatingText.SetText(*updatedTime))
}

func (t *UnderTable) drawLastRecord(mainUpdatingFrame *core.Frame) {
	lastRecordFrame := core.NewFrame(mainUpdatingFrame)
	lastRecordFrame.SetName("lastRecordFrame")
	lastRecordFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})

	leftLastRecordFrame := core.NewFrame(lastRecordFrame)
	leftLastRecordFrame.SetName("leftLastRecordFrame")
	leftLastRecordFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(160)
	})
	core.NewText(leftLastRecordFrame).SetText("Последняя запись: ")

	lastRecord := t.updatingController.GetLastRecord()

	rightLastRecordFrame := core.NewFrame(lastRecordFrame)
	rightLastRecordFrame.SetName("rightLastRecordFrame")

	recordTField := core.NewTextField(rightLastRecordFrame)
	recordTField.Styler(func(s *styles.Style) {
		s.Min.X.Dp(500)
		s.Border.Radius.Zero()
		s.Border.Width.Zero()
		s.Border.Offset.Zero()
	})
	recordTField.SetText(lastRecord)
	recordTField.OnInput(func(e events.Event) {
		inputText := recordTField.Text()
		err := t.updatingController.SetLastRecord(inputText)
		if err != nil {
			t.logger.Error(context.Background(), "set last record error: "+err.Error())
			core.MessageSnackbar(t.underTableFrame, "Ошибка ввода последней записи: "+err.Error())
		}
	})
}
