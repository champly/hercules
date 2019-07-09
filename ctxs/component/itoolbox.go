package component

import (
	"github.com/champly/lib4go/db"
)

type IToolBox interface {
	GetDefDB() db.IDB
	GetDB(name string) (db.IDB, error)

	Produce(queueName, value string) error
}

type ToolBox struct {
	*ComponentDB
	*ComponentMQ
}

func NewToolBox() IToolBox {
	return &ToolBox{
		ComponentDB: NewComponentDB(),
		ComponentMQ: NewComponentMQ(),
	}
}
