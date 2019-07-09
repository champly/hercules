package cron

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
	"github.com/robfig/cron"
)

type CronServer struct {
	server   *cron.Cron
	schedule cron.Schedule
	routers  []servers.Router
	handing  func(*ctxs.Context) error
}

func NewCronServer(routers []servers.Router, h interface{}) (*CronServer, error) {
	handing, ok := h.(func(*ctxs.Context) error)
	if !ok {
		panic("handing function is not func(ctx *ctxs.Context)error")
	}

	c := &CronServer{routers: routers, handing: handing}
	c.server = cron.New()
	if err := c.AddFunc(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CronServer) AddFunc() error {
	taskListConfig := configs.CronServerInfo.TaskList

	for _, taskConf := range taskListConfig {
		isExist := false
		for _, task := range c.routers {
			handler, ok := task.Handler.(func(*ctxs.Context) error)
			if !ok {
				if reflect.TypeOf(task.Handler).Kind() != reflect.Ptr {
					panic(reflect.TypeOf(task.Handler).Elem().Name() + " handler is not func(ctx *ctxs.Context)error")
				}
				panic(reflect.TypeOf(task.Handler).Name() + " handler is not func(ctx *ctxs.Context)error")
			}

			if !strings.EqualFold(task.Name, taskConf.Name) {
				continue
			}

			isExist = true
			err := c.server.AddFunc(taskConf.Time, func() {
				ctx := ctxs.GetCronContext()
				defer ctx.Put()

				if c.handing != nil {
					if err := c.handing(ctx); err != nil {
						return
					}
				}
				if err := handler(ctx); err != nil {
					fmt.Println(err)
				}
			})
			if err != nil {
				return err
			}
		}
		if !isExist {
			return errors.New(taskConf.Name + " cron task is not exist")
		}
	}

	return nil
}

func (c *CronServer) Start() error {
	go c.server.Start()
	return nil
}

func (c *CronServer) ShutDown() {
	c.server.Stop()
	fmt.Println("cron shutdown")
}

func (c *CronServer) Restart() {
	fmt.Println("cron restart")
}
