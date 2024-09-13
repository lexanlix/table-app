package domain

import (
	"strconv"
	"time"
)

type Category struct {
	Id           int
	Name         string
	MainCategory string
	Priority     int
}

func (c Category) CellCompositeId(month time.Month, year int) string {
	return c.MainCategory + c.Name + strconv.Itoa(int(month)) + strconv.Itoa(year)
}
