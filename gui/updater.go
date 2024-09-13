package gui

import (
	"strconv"
	"sync"

	"table-app/domain"

	"cogentcore.org/core/core"
)

type Updater struct {
	guiCells map[string]*core.TextField

	lock       sync.Mutex
	wgGroup    sync.WaitGroup
	updateChan chan domain.Cell
}

func NewUpdater() *Updater {
	return &Updater{
		guiCells:   make(map[string]*core.TextField),
		lock:       sync.Mutex{},
		wgGroup:    sync.WaitGroup{},
		updateChan: make(chan domain.Cell),
	}
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
					continue
				}

				tField.SetText(FormatInt(cell.Value))
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
