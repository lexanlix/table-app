package domain

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Cell struct {
	Id           string
	MainCategory string
	Category     string
	Value        int
	Month        time.Month
	Year         int
	IsUpdated    bool
	IsDeleted    bool
}

func (c Cell) CompositeId() string {
	return c.MainCategory + c.Category + strconv.Itoa(int(c.Month)) + strconv.Itoa(c.Year)
}

func (c Cell) Validate() error {
	if len(c.MainCategory) == 0 || len(c.Category) == 0 {
		return errors.New("category is empty")
	}

	if c.Month > 12 || c.Month < 1 {
		return errors.New("invalid month")
	}

	return nil
}
