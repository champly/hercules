package cron

import (
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
	"github.com/robfig/cron"
	"k8s.io/klog/v2"
)

var (
	defaultTimeInterval = "@every 2s"
)

type CronServer struct {
	server   *cron.Cron
	schedule cron.Schedule
	routers  []servers.Router
	preHand  func(*ctxs.Context) error
}

func NewCronServer(routers []servers.Router, h interface{}) (*CronServer, error) {
	preHand, ok := h.(func(*ctxs.Context) error)
	if !ok {
		panic("prehand function is not func(ctx *ctxs.Context) error")
	}

	c := &CronServer{routers: routers, preHand: preHand}
	c.server = cron.New()
	if err := c.AddFunc(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CronServer) AddFunc() error {
	taskListConfig := configs.CronServerInfo.TaskList

	for _, task := range c.routers {
		handler, ok := task.Handler.(func(*ctxs.Context) error)
		if !ok {
			if reflect.TypeOf(task.Handler).Kind() != reflect.Ptr {
				panic(reflect.TypeOf(task.Handler).Elem().Name() + " handler is not func(ctx *ctxs.Context)error")
			}
			panic(reflect.TypeOf(task.Handler).Name() + " handler is not func(ctx *ctxs.Context)error")
		}

		var exists bool
		for _, taskConf := range taskListConfig {
			if !strings.EqualFold(task.Name, taskConf.Name) {
				continue
			}

			exists = true
			err := c.server.AddFunc(taskConf.Time, func() {
				c.do(handler)
			})
			if err != nil {
				return err
			}
		}
		if !exists {
			err := c.server.AddFunc(defaultTimeInterval, func() {
				c.do(handler)
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CronServer) do(handler func(*ctxs.Context) error) {
	ctx := ctxs.GetCronContext()
	ctx.Type = ctxs.ServerTypeCron
	defer ctx.Put()

	if c.preHand != nil {
		if err := c.preHand(ctx); err != nil {
			return
		}
	}
	if err := handler(ctx); err != nil {
		ctx.Log.Error(err)
	}
}

func (c *CronServer) Start() error {
	go c.server.Start()
	return nil
}

func (c *CronServer) ShutDown() {
	c.server.Stop()
	klog.Info("cron shutdown")
}

func (c *CronServer) Restart() {
	klog.Info("cron restart")
}
