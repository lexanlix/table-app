package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrNoAccountName = errors.New("account name is required")

type Account struct {
	Id        string    `display:"-"`
	Name      string    `label:"Название"`
	Sum       int       `label:"Сумма"`
	Note      string    `label:"Комментарий"`
	IsInSum   bool      `label:"Учитывать в сумме" default:"true"`
	UpdatedAt time.Time `label:"Изменен" display:"-" table:"+" edit:"-"`
	Deleted   bool      `display:"-"`
}

func (a Account) Validate() error {
	if len(a.Name) == 0 {
		return ErrNoAccountName
	}

	return nil
}

func GetStartingAccounts() []Account {
	return []Account{
		{
			Id:      uuid.New().String(),
			Name:    "Счет 1",
			Sum:     0,
			Note:    "",
			IsInSum: true,
		},
	}
}
