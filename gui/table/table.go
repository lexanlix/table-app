package table

import (
	"strconv"

	"table-app/conf"
	"table-app/domain"
	"table-app/gui/iface"
	custom "table-app/gui/styles"
	"table-app/gui/updaters"
	"table-app/internal/log"

	"cogentcore.org/core/core"
)

type Table struct {
	logger     log.Logger
	year       int
	data       *domain.GuiTableData
	settings   conf.Setting
	controller iface.TableController

	updater    *updaters.Updater
	sumUpdater *updaters.SumUpdater

	tableFrame *core.Frame
	yearFrame  *core.Frame
	styler     *custom.FrameStyle
}

func NewTable(
	frame *core.Frame,
	year int,
	data *domain.GuiTableData,
	settings conf.Setting,
	controller iface.TableController,
	updater *updaters.Updater,
	sumUpdater *updaters.SumUpdater,
) *Table {
	styler := custom.NewFrameStyle(settings.Gui.CellSizeDpX, settings.Gui.CellSizeDpY)

	tableFrame := core.NewFrame(frame)
	tableFrame.SetName("table_" + strconv.Itoa(year) + "_Frame")
	tableFrame.Styler(styler.TableFrameStyle())

	yearFrame := core.NewFrame(tableFrame)
	yearFrame.SetName("yearFrame")
	yearFrame.Styler(styler.YearFrameStyle())
	core.NewText(yearFrame).SetText(strconv.Itoa(year) + " год")

	return &Table{
		year:       year,
		data:       data,
		settings:   settings,
		controller: controller,
		updater:    updater,
		sumUpdater: sumUpdater,
		tableFrame: tableFrame,
		yearFrame:  yearFrame,
		styler:     styler,
	}
}

// Draw отрисовывает табличку каждого года
func (t Table) Draw() {
	t.drawTableHead()
	t.drawMonthsColumn()
	t.drawValuesGrid()
}

// getCellSizeDpX возвращает длину ячейки в зависимости от длины текста
func (t Table) getCellSizeDpX(nameLen int) float32 {
	if nameLen < 8 {
		return 80
	}

	if nameLen > 12 {
		return t.settings.Gui.CellSizeDpX
	}

	return float32(nameLen * 11)
}
