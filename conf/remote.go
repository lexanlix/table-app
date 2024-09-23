package conf

import (
	db "table-app/internal/db/client"
	"table-app/internal/log"
)

type Remote struct {
	LogLevel log.Level `schemaGen:"logLevel" schema:"Уровень логирования"`
	Database db.StorageConfig
	Settings Setting
}

type Setting struct {
	StartYear         int
	StartMonth        int
	StartMoney        int
	Gui               Gui
	MainCategoryOrder Order
}

type Gui struct {
	CellSizeDpX float32
	CellSizeDpY float32
}

type Order map[string]int
