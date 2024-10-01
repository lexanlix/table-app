package repository

import (
	"strconv"
	"sync"
	"time"

	"table-app/conf"
	"table-app/domain"
	"table-app/utils"

	"github.com/pkg/errors"
)

type CalculationCache struct {
	consumptionByDate map[string]int
	balanceByDate     map[string]int
	mutex             sync.Mutex
	settings          conf.Setting
}

func NewCalculationCache(settings conf.Setting) *CalculationCache {
	return &CalculationCache{
		consumptionByDate: make(map[string]int),
		balanceByDate:     make(map[string]int),
		mutex:             sync.Mutex{},
		settings:          settings,
	}
}

func (r *CalculationCache) InitCache(valuesList map[string]domain.Cell, categories [][]domain.Category) error {
	// заполнение consumptionByDate
	for _, cell := range valuesList {
		if cell.MainCategory != "Расходы" {
			continue
		}

		compositeDate := utils.GetCompositeDate(int(cell.Month), cell.Year)

		consumption, isExist := r.consumptionByDate[compositeDate]
		if !isExist {
			r.consumptionByDate[compositeDate] = cell.Value
			continue
		}

		consumption += cell.Value
		r.consumptionByDate[compositeDate] = consumption
	}

	mainCategIdx := r.settings.MainCategoryOrder["Доходы"]

	var lastMonth time.Month
	for year := r.settings.StartYear; year <= time.Now().Year(); year++ {
		if year != time.Now().Year() {
			// если год не текущий
			lastMonth = time.December
		} else {
			lastMonth = time.Now().Month()
		}

		for month := 1; month <= int(lastMonth); month++ {
			balanceSum := 0

			// считаем сумму доходов
			for i, mainCategory := range categories {
				if i != mainCategIdx {
					continue
				}

				for _, categ := range mainCategory {
					compositeId := categ.MainCategory + categ.Name + strconv.Itoa(month) + strconv.Itoa(year)
					cell, ok := valuesList[compositeId]
					if ok {
						balanceSum += cell.Value
					}
				}
			}

			// берем сумму расходов
			compositeDate := utils.GetCompositeDate(month, year)
			var consumption int
			var ok bool
			consumption, ok = r.consumptionByDate[compositeDate]
			if !ok {
				consumption = 0
			}

			// и остаток предыдущего месяца
			prevBalance, err := r.getPreviousBalance(month, year)
			if err != nil {
				return errors.WithMessage(err, "get previous balance")
			}

			balanceSum = balanceSum - consumption + prevBalance
			r.balanceByDate[compositeDate] = balanceSum
		}
	}

	return nil
}

func (r *CalculationCache) UpsertConsumption(month, year, newValue int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	compositeDate := utils.GetCompositeDate(month, year)
	r.consumptionByDate[compositeDate] = newValue
}

func (r *CalculationCache) UpsertBalance(month, year, newValue int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	compositeDate := utils.GetCompositeDate(month, year)
	r.balanceByDate[compositeDate] = newValue
}

func (r *CalculationCache) GetConsumption(month, year int) (int, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	compositeDate := utils.GetCompositeDate(month, year)
	value, ok := r.consumptionByDate[compositeDate]
	return value, ok
}

func (r *CalculationCache) GetBalance(month, year int) (int, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	compositeDate := utils.GetCompositeDate(month, year)
	value, ok := r.balanceByDate[compositeDate]
	return value, ok
}

func (r *CalculationCache) getPreviousBalance(month, year int) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if year == r.settings.StartYear {
		if month == r.settings.StartMonth {
			return r.settings.StartMoney, nil
		}

		balance, ok := r.balanceByDate[utils.GetCompositeDate(month-1, year)]
		if !ok {
			return 0, errors.Errorf("not found month balance, %s %d", time.Month(month-1).String(), year)
		}

		return balance, nil
	}

	if month == 1 {
		// если Январь, то берем Декабрь предыдущего года
		balance, ok := r.balanceByDate[utils.GetCompositeDate(int(time.December), year-1)]
		if !ok {
			return 0, errors.Errorf("not found month balance, %s %d", time.December, year-1)
		}
		return balance, nil
	}

	balance, ok := r.balanceByDate[utils.GetCompositeDate(month-1, year)]
	if !ok {
		return 0, errors.Errorf("not found month balance, %s %d", time.Month(month-1).String(), year)
	}

	return balance, nil
}
