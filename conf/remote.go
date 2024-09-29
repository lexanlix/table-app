package conf

import (
	db "table-app/internal/db/client"
	"table-app/internal/log"
)

type Remote struct {
	LogLevel log.Level `schemaGen:"logLevel" schema:"Уровень логирования"`
	Storage  Storage
	Settings Setting
}

type Storage struct {
	Files    *Files
	Database *db.StorageConfig
}

type Files struct {
	TableFilePath    string
	CategoryFilePath string
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
