package domain

import (
	"table-app/conf"
)

type GuiTableData struct {
	// Categories массив основных категорий, каждый элемент которого содержит массив обычных категорий данной основной категории
	Categories [][]Category

	// ValuesList мап ячеек таблицы, ключом является compositeId, состоящий из года, месяца, категории и основной категории
	ValuesList map[string]Cell

	// MainCategoryOrder порядок приоритетов основных категорий
	MainCategoryOrder conf.Order

	// Accounts указатель на массив счетов
	Accounts *[]Account
}
