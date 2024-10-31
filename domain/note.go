package domain

import (
	"cmp"
	"slices"
	"time"
)

type Note struct {
	Id        string
	Text      string
	UpdatedAt time.Time
	Deleted   bool
}

type NoteList []Note

func (n NoteList) SortByDateAsc() {
	slices.SortFunc(n, func(a, b Note) int {
		return cmp.Compare(a.UpdatedAt.Unix(), b.UpdatedAt.Unix())
	})
}
