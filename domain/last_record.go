package domain

import "time"

type LastRecord struct {
	Id        string
	UpdatedAt time.Time
	Note      string
	Sum       int
	IsPlus    bool
}
