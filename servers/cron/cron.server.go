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
	for _, r := range c.routers {
		toolBox, ok := r.ToolBox.(component.IToolBox)
		if !ok {
			v := reflect.TypeOf(r.ToolBox)
			return errors.New(v.Elem().Name() + " constructor is not assignment component.IToolBox")
		}

		c.server.AddFunc(r.Cron, func() {
			ctx := ctxs.GetDContext()
			defer ctx.Put()

			ctx.ToolBox = toolBox
			if c.handing != nil {
				if err := c.handing(ctx); err != nil {
					return
				}
			}

			if err := r.Handler(ctx); err != nil {
				fmt.Println(err)
			}
		})
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
