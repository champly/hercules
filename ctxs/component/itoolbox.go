package component

import (
	"github.com/champly/lib4go/db"
)

type IToolBox interface {
	GetDefDB() db.IDB
	GetDB(name string) (db.IDB, error)
}

type ToolBox struct {
	*ComponentDB
}

func NewToolBox() IToolBox {
	return &ToolBox{
		ComponentDB: NewComponentDB(),
	}
}
