package component

import (
	"github.com/champly/hercules/configs"
	cmap "github.com/champly/lib4go/concurrent"
	"github.com/champly/lib4go/db"
)

type IDBComponent interface {
}

type NewDBComponent struct {
	pool cmap.ConcurrentMap
}

func NewDBComponent(conf configs.DBConfig) (db.IDB, error) {
	d, err := db.NewDB(conf.Provider, conf.ConnString, conf.MaxOpen, conf.MaxIdle, conf.MaxLifeTime)
	if err != nil {
		return nil, err
	}
}
