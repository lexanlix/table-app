package domain

import (
	"strconv"

	"table-app/conf"
)

type GuiTableData struct {
	Categories        [][]Category
	ValuesList        map[string]Cell
	MainCategoryOrder conf.Order
}

func (t GuiTableData) GetConsumptionSum(month, year int) int {
	res := 0
	idx := t.MainCategoryOrder["Расходы"]

	for i, mainCategory := range t.Categories {
		if i != idx {
			continue
		}
		for _, categ := range mainCategory {
			compositeId := categ.MainCategory + categ.Name + strconv.Itoa(month) + strconv.Itoa(year)
			cell, ok := t.ValuesList[compositeId]
			if ok {
				res += cell.Value
			}
		}
	}

	return res
}
