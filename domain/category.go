package domain

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

const startCategoryName = "Категория"

type Category struct {
	Id           string
	Name         string
	MainCategory string
	Priority     int
}

func (c Category) CellCompositeId(month time.Month, year int) string {
	return c.MainCategory + c.Name + strconv.Itoa(int(month)) + strconv.Itoa(year)
}

func GetStartingCategories(mainCategories map[string]int) []Category {
	result := make([]Category, 0)
	for mainCategory, priority := range mainCategories {
		category := Category{
			Id:           uuid.New().String(),
			Name:         startCategoryName + " " + strconv.Itoa(priority+1),
			MainCategory: mainCategory,
			Priority:     1,
		}

		result = append(result, category)
	}

	return result
}
