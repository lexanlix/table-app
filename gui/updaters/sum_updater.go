package updaters

import (
	"context"
	"sync"

	"table-app/entity"
	"table-app/gui/iface"
	"table-app/gui/styles/format"
	"table-app/internal/log"
	"table-app/utils"

	"cogentcore.org/core/core"
)

// SumUpdater Обновляет данные таблицы при изменении ее содержимого,
// также обновляет данные о суммах и время последнего изменения в таблице
type SumUpdater struct {
	logger            log.Logger
	consumptionFields map[string]*core.Text
	balanceFields     map[string]*core.Text
	mainSumText       *core.Text
	diffSumText       *core.Text
	tableController   iface.TableController
	accountController iface.AccountController

	updatedTimeText *core.Text

	lock               sync.Mutex
	wgGroup            sync.WaitGroup
	updateChan         chan entity.MonthYear
	updateAccountsChan chan struct{}
}

func NewSumUpdater(
	logger log.Logger,
	tableController iface.TableController,
	accountController iface.AccountController,
) *SumUpdater {
	return &SumUpdater{
		logger:             logger,
		consumptionFields:  make(map[string]*core.Text),
		balanceFields:      make(map[string]*core.Text),
		tableController:    tableController,
		accountController:  accountController,
		lock:               sync.Mutex{},
		wgGroup:            sync.WaitGroup{},
		updateChan:         make(chan entity.MonthYear),
		updateAccountsChan: make(chan struct{}),
	}
}

func (u *SumUpdater) GetUpdateChan() chan entity.MonthYear {
	return u.updateChan
}

func (u *SumUpdater) GetUpdateAccountsChan() chan struct{} {
	return u.updateAccountsChan
}

func (u *SumUpdater) SendToChannel(obj entity.MonthYear) {
	u.updateChan <- obj
	return
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
				consumption := u.tableController.GetConsumptionSum(date.Month, date.Year)

				balanceById, err := u.tableController.UpsertBalance(date.Month, date.Year)
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
					u.lock.Unlock()
					continue
				}

				consumptionField.SetText(format.FormatInt(consumption, format.AddMinus))
				consumptionField.Update()

				for id, balance := range balanceById {
					balanceField, ok := u.balanceFields[id]
					if !ok {
						u.logger.Error(context.Background(), "not found balance field",
							log.String("compositeDate", compositeDate))
						u.lock.Unlock()
						continue
					}
					balanceField.SetText(format.FormatInt(balance))
					balanceField.Update()
					u.balanceFields[id] = balanceField
				}

				u.consumptionFields[compositeDate] = consumptionField

				u.updatedTimeText.Update()
				u.updateSumTexts()

				u.lock.Unlock()
			}
		case _, isOpen := <-u.updateAccountsChan:
			{
				if !isOpen {
					return
				}
				u.lock.Lock()
				u.updateSumTexts()
				u.lock.Unlock()
			}
		}
	}
}

func (u *SumUpdater) Close() {
	close(u.updateChan)
	close(u.updateAccountsChan)
}

func (u *SumUpdater) AddConsumptionText(month, year int, tField *core.Text) {
	u.lock.Lock()

	sum := u.tableController.GetConsumptionSum(month, year)
	tField.SetText(format.FormatInt(sum, format.AddMinus))

	compositeDate := utils.GetCompositeDate(month, year)
	u.consumptionFields[compositeDate] = tField

	u.lock.Unlock()
}

func (u *SumUpdater) AddBalanceText(month, year int, tField *core.Text) {
	u.lock.Lock()

	var sum int
	var err error

	sum, err = u.tableController.GetBalanceSum(month, year)
	if err != nil {
		u.logger.Error(context.Background(), "get balance sum", log.Any("err", err))
		sum = 0
	}

	tField.SetText(format.FormatInt(sum))

	compositeDate := utils.GetCompositeDate(month, year)
	u.balanceFields[compositeDate] = tField

	u.lock.Unlock()
}

func (u *SumUpdater) AddUpdatedTimeText(updatedText *core.Text) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.updatedTimeText = updatedText
}

func (u *SumUpdater) AddSumTexts(mainSumText, diffSumText *core.Text) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.mainSumText = mainSumText
	u.diffSumText = diffSumText
}

func (u *SumUpdater) updateSumTexts() {
	sums, err := u.accountController.GetSum(context.Background())
	if err != nil {
		u.logger.Error(context.Background(), "get account sums: "+err.Error())
		return
	}

	u.mainSumText.SetText(format.FormatInt(sums.MainSum))
	u.mainSumText.Update()

	u.diffSumText.SetText(format.FormatInt(sums.DiffSum))
	u.diffSumText.Update()
}
