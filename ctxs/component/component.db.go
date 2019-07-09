package component

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/champly/hercules/configs"
	cmap "github.com/champly/lib4go/concurrent"
	"github.com/champly/lib4go/db"
)

type IComponentDB interface {
	GetDefDB() db.IDB
	GetDB(name string) (db.IDB, error)
}

type ComponentDB struct {
	defName string
	pool    cmap.ConcurrentMap
}

var (
	componentDB *ComponentDB
	lockDB      sync.Mutex
)

func NewComponentDB() *ComponentDB {
	if componentDB != nil {
		return componentDB
	}
	lockDB.Lock()
	defer lockDB.Unlock()

	if componentDB != nil {
		return componentDB
	}

	// if configs.DBInfo == nil || len(configs.DBInfo.List) < 1 {
	// 	panic("db config is empty, can't use db component")
	// }
	componentDB = &ComponentDB{
		pool: cmap.New(2),
	}
	for _, dbConf := range configs.DBInfo.List {
		if dbConf.Default {
			componentDB.defName = dbConf.Name
			break
		}
	}
	return componentDB
}

func (c *ComponentDB) GetDefDB() db.IDB {
	if c.defName == "" && len(configs.DBInfo.List) > 0 {
		c.defName = configs.DBInfo.List[0].Name
	}
	db, err := c.GetDB(c.defName)
	if err != nil {
		panic(err)
	}
	return db
}

func (c *ComponentDB) GetDB(name string) (db.IDB, error) {
	_, dbObj, err := c.pool.SetIfAbsentCb(name, func(key string, input ...interface{}) (interface{}, error) {
		if name == "" {
			panic("db config name is empty")
		}
		for _, conf := range configs.DBInfo.List {
			if strings.EqualFold(conf.Name, name) {
				return db.NewDB(conf.Provider, conf.ConnString, conf.MaxIdle, conf.MaxOpen, conf.MaxLifeTime)
			}
		}
		return nil, errors.New(name + " is not config")
	})
	if err != nil {
		return nil, fmt.Errorf("创建db失败: err:%s", err)
	}
	return dbObj.(db.IDB), nil
}
