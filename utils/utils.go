package utils

import "strconv"

func GetCompositeDate(month, year int) string {
	return strconv.Itoa(month) + strconv.Itoa(year)
}
