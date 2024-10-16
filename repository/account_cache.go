package repository

import (
	"sync"
	"time"

	"table-app/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// AccountCache
// use mutex functions outside
type AccountCache struct {
	accountIndexById map[string]int
	accounts         []domain.Account

	mutex sync.Mutex
}

func NewAccountCache() *AccountCache {
	return &AccountCache{
		accountIndexById: make(map[string]int),
		accounts:         make([]domain.Account, 0),
		mutex:            sync.Mutex{},
	}
}

func (r *AccountCache) InitCache(accounts []domain.Account) {
	for _, account := range accounts {
		r.accountIndexById[account.Id] = len(r.accounts)
		r.accounts = append(r.accounts, account)
	}
}

func (r *AccountCache) Insert(newAcc domain.Account) error {
	newAcc.Id = uuid.New().String()
	newAcc.UpdatedAt = time.Now()

	r.accountIndexById[newAcc.Id] = len(r.accounts)
	r.accounts = append(r.accounts, newAcc)
	return nil
}

func (r *AccountCache) ReadAll() []domain.Account {
	all := make([]domain.Account, 0)
	for _, acc := range r.accounts {
		all = append(all, acc)
	}

	return all
}

func (r *AccountCache) GetListPtr() *[]domain.Account {
	return &r.accounts
}

func (r *AccountCache) IsInCache(account domain.Account) bool {
	_, ok := r.accountIndexById[account.Id]
	return ok
}

func (r *AccountCache) UpdateAccount(account domain.Account) error {
	account.UpdatedAt = time.Now()

	idx, ok := r.accountIndexById[account.Id]
	if !ok {
		return errors.Errorf("account %s %s not found", account.Id, account.Name)
	}

	r.accounts[idx] = account
	return nil
}

func (r *AccountCache) Delete(account domain.Account) {
	idx, ok := r.accountIndexById[account.Id]
	if !ok {
		return
	}

	r.accounts[idx] = account
}

func (r *AccountCache) Lock() {
	r.mutex.Lock()
}

func (r *AccountCache) Unlock() {
	r.mutex.Unlock()
}
