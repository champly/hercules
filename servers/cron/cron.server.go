package cron

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/champly/hercules/component"
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/robfig/cron"
)

type CronServer struct {
	server   *cron.Cron
	schedule cron.Schedule
	routers  []ctxs.Router
	handing  func(*ctxs.Context) error
}

func NewCronServer(routers []ctxs.Router, handing func(*ctxs.Context) error) (*CronServer, error) {
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
			if !strings.EqualFold(task.Name, taskConf.Name) {
				continue
			}

			isExist = true
			toolBox, ok := task.ToolBox.(component.IToolBox)
			if !ok {
				v := reflect.TypeOf(task.ToolBox)
				return errors.New(v.Elem().Name() + " constructor is not assignment component.IToolBox")
			}

			err := c.server.AddFunc(taskConf.Time, func() {
				ctx := ctxs.GetDContext()
				defer ctx.Put()

				ctx.ToolBox = toolBox
				if c.handing != nil {
					if err := c.handing(ctx); err != nil {
						return
					}
				}

				if err := task.Handler(ctx); err != nil {
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
	if !strings.EqualFold(configs.CronServerInfo.Status, "start") {
		return fmt.Errorf("cron server config is: %s", configs.CronServerInfo.Status)
	}

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
