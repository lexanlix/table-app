package updaters

import (
	"context"
	"strconv"
	"sync"

	"table-app/domain"
	"table-app/gui/styles/format"
	"table-app/internal/log"

	"cogentcore.org/core/core"
)

// Updater обновляет значения ячеек при вводе суммы через сумматор
type Updater struct {
	logger   log.Logger
	guiCells map[string]*core.TextField

	lock       sync.Mutex
	wgGroup    sync.WaitGroup
	updateChan chan domain.Cell
}

func NewUpdater(logger log.Logger) *Updater {
	return &Updater{
		logger:     logger,
		guiCells:   make(map[string]*core.TextField),
		lock:       sync.Mutex{},
		wgGroup:    sync.WaitGroup{},
		updateChan: make(chan domain.Cell),
	}
}

func (u *Updater) GetUpdateChan() chan domain.Cell {
	return u.updateChan
}

func (u *Updater) Start() {
	go u.start()
}

func (u *Updater) start() {
	for {
		select {
		case cell, isOpen := <-u.updateChan:
			{
				if !isOpen {
					return
				}

				compositeId := cell.MainCategory + cell.Category + strconv.Itoa(int(cell.Month)) + strconv.Itoa(cell.Year)
				u.lock.Lock()
				tField, ok := u.guiCells[compositeId]
				if !ok {
					u.logger.Error(context.Background(), "cell not found by id",
						log.String("compositeId", compositeId))
					continue
				}

				tField.SetText(format.FormatInt(cell.Value))
				u.guiCells[compositeId] = tField
				u.lock.Unlock()
			}
		}
	}
}

func (u *Updater) Close() {
	close(u.updateChan)
}

func (u *Updater) AddTextField(compositeId string, tField *core.TextField) {
	u.lock.Lock()
	u.guiCells[compositeId] = tField
	u.lock.Unlock()
}
