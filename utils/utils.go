package utils

import "strconv"

func GetCompositeDate(month, year int) string {
	return strconv.Itoa(month) + strconv.Itoa(year)
}

func GetCompositeId(mainCategory, category string, month, year int) string {
	return mainCategory + category + strconv.Itoa(month) + strconv.Itoa(year)
}

func GetCompositeCategory(mainCategory, category string) string {
	return mainCategory + category
}
