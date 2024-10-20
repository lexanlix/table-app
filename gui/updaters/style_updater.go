package updaters

import (
	"sync"

	"table-app/internal/log"

	"cogentcore.org/core/core"
)

// StyleUpdater обновляет стили отображения элементов по сигналу
type StyleUpdater struct {
	logger  log.Logger
	appBody *core.Body

	lock       sync.Mutex
	wgGroup    sync.WaitGroup
	updateChan chan struct{}
}

func NewStyleUpdater(logger log.Logger, appBody *core.Body) *StyleUpdater {
	return &StyleUpdater{
		logger:     logger,
		appBody:    appBody,
		lock:       sync.Mutex{},
		wgGroup:    sync.WaitGroup{},
		updateChan: make(chan struct{}),
	}
}

func (su *StyleUpdater) GetUpdateChan() chan struct{} {
	return su.updateChan
}

func (su *StyleUpdater) Start() {
	go su.start()
}

func (su *StyleUpdater) start() {
	for {
		select {
		case _, isOpen := <-su.updateChan:
			{
				if !isOpen {
					return
				}

				su.lock.Lock()

				su.appBody.WidgetBase.AsyncLock()
				su.appBody.Update()
				su.appBody.WidgetBase.AsyncUnlock()

				su.lock.Unlock()
			}
		}
	}
}

func (su *StyleUpdater) Close() {
	close(su.updateChan)
}
