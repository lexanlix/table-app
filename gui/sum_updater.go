package gui

import (
	"context"
	"sync"

	"table-app/entity"
	"table-app/internal/log"
	"table-app/utils"

	"cogentcore.org/core/core"
)

type SumUpdater struct {
	logger            log.Logger
	consumptionFields map[string]*core.Text
	balanceFields     map[string]*core.Text
	controller        TableController

	lock       sync.Mutex
	wgGroup    sync.WaitGroup
	updateChan chan entity.MonthYear
}

func NewSumUpdater(logger log.Logger, controller TableController) *SumUpdater {
	return &SumUpdater{
		logger:            logger,
		consumptionFields: make(map[string]*core.Text),
		balanceFields:     make(map[string]*core.Text),
		controller:        controller,
		lock:              sync.Mutex{},
		wgGroup:           sync.WaitGroup{},
		updateChan:        make(chan entity.MonthYear),
	}
}

func (u *SumUpdater) Start() {
	go u.start()
}

func (u *SumUpdater) start() {
	for {
		select {
		case date, isOpen := <-u.updateChan:
			{
				if !isOpen {
					return
				}

				// сначала изменяются расходы, затем остаток, так как он пересчитывается с учетом расходов
				consumption := u.controller.GetConsumptionSum(date.Month, date.Year)

				balanceById, err := u.controller.UpsertBalance(date.Month, date.Year)
				if err != nil {
					u.logger.Error(context.Background(), "get balance sum", log.Any("err", err))
					continue
				}

				compositeDate := utils.GetCompositeDate(date.Month, date.Year)
				u.lock.Lock()

				consumptionField, ok := u.consumptionFields[compositeDate]
				if !ok {
					u.logger.Error(context.Background(), "not found consumption field",
						log.String("compositeDate", compositeDate))
					continue
				}

				consumptionField.SetText(FormatInt(consumption, addMinus))

				for id, balance := range balanceById {
					balanceField, ok := u.balanceFields[id]
					if !ok {
						u.logger.Error(context.Background(), "not found balance field",
							log.String("compositeDate", compositeDate))
						continue
					}
					balanceField.SetText(FormatInt(balance))
					u.balanceFields[id] = balanceField
				}

				u.consumptionFields[compositeDate] = consumptionField

				u.lock.Unlock()
			}
		}
	}
}

func (u *SumUpdater) Close() {
	close(u.updateChan)
}

func (u *SumUpdater) AddConsumptionText(month, year int, tField *core.Text) {
	u.lock.Lock()

	sum := u.controller.GetConsumptionSum(month, year)
	tField.SetText(FormatInt(sum, addMinus))

	compositeDate := utils.GetCompositeDate(month, year)
	u.consumptionFields[compositeDate] = tField

	u.lock.Unlock()
}

func (u *SumUpdater) AddBalanceText(month, year int, tField *core.Text) {
	u.lock.Lock()

	var sum int
	var err error

	sum, err = u.controller.GetBalanceSum(month, year)
	if err != nil {
		u.logger.Error(context.Background(), "get balance sum", log.Any("err", err))
		sum = 0
	}

	tField.SetText(FormatInt(sum))

	compositeDate := utils.GetCompositeDate(month, year)
	u.balanceFields[compositeDate] = tField

	u.lock.Unlock()
}
