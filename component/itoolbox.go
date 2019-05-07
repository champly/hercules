package component

import (
	"github.com/champly/hercules/configs"
	"github.com/champly/lib4go/db"
)

type IToolBox interface {
	SetDefault(name string, conf configs.DBConfig)
	GetDefDB() (db.IDB, error)
	GetDB(name string, conf configs.DBConfig) (db.IDB, error)
}

type ToolBox struct {
	*ComponentDB
}

func NewToolBox() IToolBox {
	return &ToolBox{ComponentDB: NewComponentDB()}
}
