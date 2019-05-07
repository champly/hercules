package component

import (
	"fmt"
	"sync"

	"github.com/champly/hercules/configs"
	cmap "github.com/champly/lib4go/concurrent"
	"github.com/champly/lib4go/db"
)

type IComponentDB interface {
	SetDefault(name string, conf configs.DBConfig)
	GetDefDB() (db.IDB, error)
	GetDB(name string, conf configs.DBConfig) (db.IDB, error)
}

type ComponentDB struct {
	defName   string
	defConfig configs.DBConfig
	pool      cmap.ConcurrentMap
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
	return componentDB
}

func (c *ComponentDB) SetDefault(name string, conf configs.DBConfig) {
	c.defName = name
	c.defConfig = conf
}

func (c *ComponentDB) GetDefDB() (db.IDB, error) {
	return c.GetDB(c.defName, c.defConfig)
}

func (c *ComponentDB) GetDB(name string, conf configs.DBConfig) (db.IDB, error) {
	_, dbObj, err := c.pool.SetIfAbsentCb(name, func(key string, input ...interface{}) (interface{}, error) {
		return db.NewDB(conf.Provider, conf.ConnString, conf.MaxIdle, conf.MaxOpen, conf.MaxLifeTime)
	})
	if err != nil {
		return nil, fmt.Errorf("创建db失败:config:%+v, err:%s", conf, err)
	}
	return dbObj.(db.IDB), nil
}
