package repository

import (
	"sort"
	"sync"

	"table-app/conf"
	"table-app/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CategoryCache
// use mutex functions outside
type CategoryCache struct {
	// mainCategoryPriorityByName - это
	// мап приоритетов основных категорий
	mainCategoryPriorityByName map[string]int

	// orderArr
	// используется для определения порядка
	orderArr [][]domain.Category

	// categoryIndexByName - map[string][]int, где []int = {mainCategoryIndex, categoryIndex};
	// используется для поиска категорий в кеше
	categoryIndexByName map[string][]int

	mutex sync.Mutex
}

func NewCategoryCache(order conf.Order) *CategoryCache {
	orderArr := make([][]domain.Category, 0)
	for i := 0; i < len(order); i++ {
		orderArr = append(orderArr, make([]domain.Category, 0))
	}

	return &CategoryCache{
		mainCategoryPriorityByName: order,
		orderArr:                   orderArr,
		categoryIndexByName:        make(map[string][]int),
		mutex:                      sync.Mutex{},
	}
}

func (r *CategoryCache) InitCache(categories []domain.Category) {
	for _, cat := range categories {
		priority, ok := r.mainCategoryPriorityByName[cat.MainCategory]
		if ok {
			r.orderArr[priority] = append(r.orderArr[priority], cat)
		}
	}

	for k := range r.orderArr {
		sort.Slice(r.orderArr[k], func(i, j int) bool {
			return r.orderArr[k][i].Priority < r.orderArr[k][j].Priority
		})
	}

	for i := range r.orderArr {
		for j := range r.orderArr[i] {
			category := r.orderArr[i][j]
			r.categoryIndexByName[category.MainCategory+category.Name] = []int{i, j}
		}
	}
}

func (r *CategoryCache) Insert(newCategory domain.Category) error {
	// находим приоритет основной категории
	mainPriority, ok := r.mainCategoryPriorityByName[newCategory.MainCategory]
	if !ok {
		return errors.Errorf("main category %s not found", newCategory.MainCategory)
	}

	// находим приоритет категории:
	// порядковый номер в массиве категорий данной основной категории + 1
	priority := len(r.orderArr[mainPriority]) + 1
	newCategory.Priority = priority
	newCategory.Id = uuid.New().String()

	r.orderArr[mainPriority] = append(r.orderArr[mainPriority], newCategory)
	r.categoryIndexByName[newCategory.MainCategory+newCategory.Name] = []int{mainPriority, priority}
	return nil
}

func (r *CategoryCache) ReadAll() []domain.Category {
	all := make([]domain.Category, 0)
	for _, catArr := range r.orderArr {
		for _, cat := range catArr {
			all = append(all, cat)
		}
	}

	return all
}

func (r *CategoryCache) GetCategoryArray() [][]domain.Category {
	return r.orderArr
}

func (r *CategoryCache) IsInCache(category domain.Category) bool {
	_, ok := r.categoryIndexByName[category.MainCategory+category.Name]
	return ok
}

func (r *CategoryCache) UpdateCategory(old, new domain.Category) error {
	idxs, ok := r.categoryIndexByName[old.MainCategory+old.Name]
	if !ok {
		return errors.Errorf("category %s %s not found", old.MainCategory, old.Name)
	}

	r.categoryIndexByName[new.MainCategory+new.Name] = idxs
	delete(r.categoryIndexByName, old.MainCategory+old.Name)

	return nil
}

func (r *CategoryCache) Lock() {
	r.mutex.Lock()
}

func (r *CategoryCache) Unlock() {
	r.mutex.Unlock()
}
