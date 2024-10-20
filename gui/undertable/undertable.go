package undertable

import (
	"table-app/domain"
	"table-app/gui/iface"
	custom "table-app/gui/styles"
	"table-app/gui/updaters"
	"table-app/internal/log"

	"cogentcore.org/core/core"
	"cogentcore.org/core/styles"
)

// UnderTable
// Область под таблицей, включающая записи о последнем обновлении,
// общую сумму по счетам и разницу с остатком,
// информацию о счетах с возможностью добавления счета
type UnderTable struct {
	logger             log.Logger
	updatingController iface.UpdatingController
	accountController  iface.AccountController

	sumUpdater      *updaters.SumUpdater
	textFieldStyler *custom.TextFieldStyle

	appBody         *core.Body
	underTableFrame *core.Frame
	upperFrame      *core.Frame
	accountList     *[]domain.Account
	updateSumChan   chan struct{}
}

func NewUnderTable(
	logger log.Logger,
	appBody *core.Body,
	frame *core.Frame,
	updatingController iface.UpdatingController,
	accountController iface.AccountController,
	sumUpdater *updaters.SumUpdater,
	accountList *[]domain.Account,
) *UnderTable {
	underTableFrame := core.NewFrame(frame)
	underTableFrame.SetName("underTableFrame")
	underTableFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	upperFrame := core.NewFrame(underTableFrame)
	upperFrame.SetName("upperFrame")
	upperFrame.Styler(func(s *styles.Style) {
		s.CenterAll()

		tableFrame := appBody.Child(1).AsTree().Child(0).AsTree().This.(*core.Frame)
		sizeX := tableFrame.Geom.Size.Actual.Total.X

		if sizeX > custom.MinUpperUnderTableFrameSize {
			s.Min.X.Dp(sizeX / custom.PixToDp)
		}
	})

	core.NewSpace(underTableFrame).Styler(func(s *styles.Style) {
		s.Min.X.Dp(20)
	})

	return &UnderTable{
		logger:             logger,
		updatingController: updatingController,
		accountController:  accountController,

		sumUpdater:      sumUpdater,
		textFieldStyler: custom.NewTextFieldStyle(),

		appBody:         appBody,
		underTableFrame: underTableFrame,
		upperFrame:      upperFrame,
		accountList:     accountList,
		updateSumChan:   sumUpdater.GetUpdateAccountsChan(),
	}
}

// Draw рисует область под таблицей
func (t *UnderTable) Draw() {
	t.drawLastUpdating()
	t.drawAccounts()
	core.NewStretch(t.upperFrame)
	t.drawSum()
}
