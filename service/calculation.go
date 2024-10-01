package service

import (
	"strconv"
	"time"

	"table-app/conf"
	"table-app/domain"
	"table-app/repository"
	"table-app/utils"

	"github.com/pkg/errors"
)

type Calculation struct {
	cache         *repository.CalculationCache
	cellsCache    *repository.CellsCache
	categoryCache *repository.CategoryCache
	settings      conf.Setting
}

func NewCalculation(
	cache *repository.CalculationCache,
	cellsCache *repository.CellsCache,
	categoryCache *repository.CategoryCache,
	settings conf.Setting,
) *Calculation {
	return &Calculation{
		cache:         cache,
		cellsCache:    cellsCache,
		categoryCache: categoryCache,
		settings:      settings,
	}
}

func (s *Calculation) ConsumptionSum(month, year int) int {
	res := 0
	s.categoryCache.Lock()
	categories := s.categoryCache.GetCategoryArray()
	s.categoryCache.Unlock()

	s.cellsCache.Lock()
	valuesList := s.cellsCache.GetList()
	s.cellsCache.Unlock()

	idx := s.settings.MainCategoryOrder["Расходы"]

	for i, mainCategory := range categories {
		if i != idx {
			continue
		}
		for _, categ := range mainCategory {
			compositeId := categ.MainCategory + categ.Name + strconv.Itoa(month) + strconv.Itoa(year)
			cell, ok := valuesList[compositeId]
			if ok {
				res += cell.Value
			}
		}
	}

	s.cache.UpsertConsumption(month, year, res)

	return res
}

func (s *Calculation) BalanceSum(month, year int) (int, error) {
	res := 0

	s.categoryCache.Lock()
	categories := s.categoryCache.GetCategoryArray()
	s.categoryCache.Unlock()

	s.cellsCache.Lock()
	valuesList := s.cellsCache.GetList()
	s.cellsCache.Unlock()

	idx := s.settings.MainCategoryOrder["Доходы"]

	for i, mainCategory := range categories {
		if i != idx {
			continue
		}
		for _, categ := range mainCategory {
			compositeId := categ.MainCategory + categ.Name + strconv.Itoa(month) + strconv.Itoa(year)
			cell, ok := valuesList[compositeId]
			if ok {
				res += cell.Value
			}
		}
	}

	consumption, ok := s.cache.GetConsumption(month, year)
	if !ok {
		return 0, errors.Errorf("not found month consumption, %s %d", time.Month(month).String(), year)
	}

	prevBalance, err := s.getPreviousBalance(month, year)
	if err != nil {
		return 0, errors.Wrapf(err, "get previous balance")
	}

	res = res - consumption + prevBalance

	s.cache.UpsertBalance(month, year, res)

	return res, nil
}

func (s *Calculation) UpsertBalance(currentMonth, currentYear int) (map[string]int, error) {
	res := make(map[string]int)

	s.categoryCache.Lock()
	categories := s.categoryCache.GetCategoryArray()
	s.categoryCache.Unlock()

	s.cellsCache.Lock()
	valuesList := s.cellsCache.GetList()
	s.cellsCache.Unlock()

	idx := s.settings.MainCategoryOrder["Доходы"]

	var firstMonth, lastMonth time.Month
	for year := currentYear; year <= time.Now().Year(); year++ {
		if year != currentYear {
			// если год не текущий
			firstMonth = time.January
		} else {
			firstMonth = time.Month(currentMonth)
		}

		if year != time.Now().Year() {
			// если год не текущий
			lastMonth = time.December
		} else {
			lastMonth = time.Now().Month()
		}

		for month := firstMonth; month <= lastMonth; month++ {
			sum := 0
			for i, mainCategory := range categories {
				if i != idx {
					continue
				}
				for _, categ := range mainCategory {
					compositeId := categ.MainCategory + categ.Name + strconv.Itoa(int(month)) + strconv.Itoa(year)
					cell, ok := valuesList[compositeId]
					if ok {
						sum += cell.Value
					}
				}
			}

			consumption, ok := s.cache.GetConsumption(int(month), year)
			if !ok {
				return nil, errors.Errorf("not found month consumption, %s %d", month.String(), year)
			}

			prevBalance, err := s.getPreviousBalance(int(month), year)
			if err != nil {
				return nil, errors.Wrapf(err, "get previous balance")
			}

			sum = sum - consumption + prevBalance

			s.cache.UpsertBalance(int(month), year, sum)
			res[utils.GetCompositeDate(int(month), year)] = sum
		}
	}

	return res, nil
}

func (s *Calculation) getPreviousBalance(month, year int) (int, error) {
	if year == s.settings.StartYear {
		if month == s.settings.StartMonth {
			return s.settings.StartMoney, nil
		}

		balance, ok := s.cache.GetBalance(month-1, year)
		if !ok {
			return 0, errors.Errorf("not found month balance, %s %d", time.Month(month-1).String(), year)
		}

		return balance, nil
	}

	if month == 1 {
		// если Январь, то берем Декабрь предыдущего года
		balance, ok := s.cache.GetBalance(int(time.December), year-1)
		if !ok {
			return 0, errors.Errorf("not found month balance, %s %d", time.December, year-1)
		}
		return balance, nil
	}

	balance, ok := s.cache.GetBalance(month-1, year)
	if !ok {
		return 0, errors.Errorf("not found month balance, %s %d", time.Month(month-1).String(), year)
	}

	return balance, nil
}

func (s *Calculation) GetAnnualResult(year int) map[string]int {
	res := make(map[string]int)

	s.categoryCache.Lock()
	categories := s.categoryCache.GetCategoryArray()
	s.categoryCache.Unlock()

	s.cellsCache.Lock()
	valuesList := s.cellsCache.GetList()
	s.cellsCache.Unlock()

	for _, mainCategoryArr := range categories {
		for _, category := range mainCategoryArr {
			categoryResult := 0

			for month := 1; month <= int(time.December); month++ {
				compositeId := utils.GetCompositeId(category.MainCategory, category.Name, month, year)
				cell, ok := valuesList[compositeId]
				if ok {
					categoryResult += cell.Value
				}
			}

			compositeCategory := utils.GetCompositeCategory(category.MainCategory, category.Name)
			res[compositeCategory] = categoryResult
		}
	}

	consumptionResult := 0
	for month := 1; month <= int(time.December); month++ {
		consumption, ok := s.cache.GetConsumption(month, year)
		if ok {
			consumptionResult += consumption
		}
	}

	var balanceResult int
	var ok bool
	balanceResult, ok = s.cache.GetBalance(int(time.December), year)
	if !ok {
		balanceResult = 0
	}

	res[domain.ColumnConsumption] = consumptionResult
	res[domain.ColumnBalance] = balanceResult

	return res
}
