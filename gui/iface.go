package gui

import (
	"context"

	"table-app/domain"
)

type TableController interface {
	UpsertValue(ctx context.Context, cell domain.Cell) error
	AddCategory(ctx context.Context, category domain.Category) error
	UpdateCategoryName(ctx context.Context, old, new domain.Category) error
	CategoryIsExist(ctx context.Context, category domain.Category) bool
	SaveAll(ctx context.Context) error
}
