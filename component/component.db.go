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
	GetDefDB() (db.IDB, error)
	GetDB(name string) (db.IDB, error)
}

type ComponentDB struct {
	defName string
	pool    cmap.ConcurrentMap
}

var (
	componentDB *ComponentDB
	lock        sync.Mutex
)

func NewComponentDB() *ComponentDB {
	if componentDB != nil {
		return componentDB
	}
	lock.Lock()
	defer lock.Unlock()

	if componentDB != nil {
		return componentDB
	}
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

func (c *ComponentDB) GetDefDB() (db.IDB, error) {
	return c.GetDB(c.defName)
}

func (c *ComponentDB) GetDB(name string) (db.IDB, error) {
	_, dbObj, err := c.pool.SetIfAbsentCb(name, func(key string, input ...interface{}) (interface{}, error) {
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
