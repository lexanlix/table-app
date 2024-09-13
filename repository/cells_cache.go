package repository

import (
	"sync"
	"time"

	"table-app/domain"

	"github.com/google/uuid"
)

// CellsCache
// use mutex functions outside
type CellsCache struct {
	cache map[string]domain.Cell
	mutex sync.Mutex
}

func NewCellsCache() *CellsCache {
	return &CellsCache{
		cache: make(map[string]domain.Cell),
		mutex: sync.Mutex{},
	}
}

func (r *CellsCache) InitCache(cells []domain.Cell) {
	for _, cell := range cells {
		compositeId := cell.CompositeId()
		cell.IsUpdated = false
		r.cache[compositeId] = cell
	}
}

func (r *CellsCache) Upsert(newCell domain.Cell) {
	compositeId := newCell.CompositeId()
	cell, isExist := r.cache[compositeId]
	if !isExist {
		newCell.IsUpdated = true
		newCell.Id = uuid.New().String()
		r.cache[compositeId] = newCell
		return
	}

	cell.Value = newCell.Value
	cell.IsUpdated = true
	r.cache[compositeId] = cell
}

func (r *CellsCache) Insert(newCell domain.Cell) {
	compositeId := newCell.CompositeId()
	newCell.IsUpdated = true
	r.cache[compositeId] = newCell
}

func (r *CellsCache) Get(compositeId string) (domain.Cell, bool) {
	cell, ok := r.cache[compositeId]
	return cell, ok
}

func (r *CellsCache) ReadAll() []domain.Cell {
	all := make([]domain.Cell, 0)
	for _, cell := range r.cache {
		all = append(all, cell)
	}
	return all
}

func (r *CellsCache) Delete(compositeId string) {
	_, ok := r.cache[compositeId]
	if ok {
		delete(r.cache, compositeId)
	}
}

func (r *CellsCache) GetList() map[string]domain.Cell {
	return r.cache
}

// UpdateCategoryName
// обновляем compositeId в кеше, так как меняется название категории
func (r *CellsCache) UpdateCategoryName(oldCategory, newCategory domain.Category, startMonth, startYear int) {
	currentYear := time.Now().Year()

	for year := startYear; year <= currentYear; year++ {
		currentMonth := time.December
		if year == currentYear {
			currentMonth = time.Now().Month()
		}

		for month := time.Month(startMonth); month <= currentMonth; month++ {
			compositeId := oldCategory.CellCompositeId(month, year)
			cell, ok := r.cache[compositeId]
			if !ok {
				continue
			}

			cell.Category = newCategory.Name
			cell.IsUpdated = true

			newCompositeId := newCategory.CellCompositeId(month, year)
			r.cache[newCompositeId] = cell
			delete(r.cache, compositeId)
		}
	}

}

func (r *CellsCache) Lock() {
	r.mutex.Lock()
}

func (r *CellsCache) Unlock() {
	r.mutex.Unlock()
}
