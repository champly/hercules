package cron

import (
	"fmt"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/robfig/cron"
)

type CronServer struct {
	server   *cron.Cron
	schedule cron.Schedule
	routers  []configs.Router
}

func NewCronServer(routers []configs.Router) (*CronServer, error) {
	c := &CronServer{routers: routers}
	c.server = cron.New()
	c.AddFunc()
	return c, nil
}

func (c *CronServer) AddFunc() {
	for _, t := range c.routers {
		c.server.AddFunc(t.Cron, func() {
			ctx := ctxs.GetDContext()
			defer ctx.Put()
			if err := t.Handler(ctx); err != nil {
				fmt.Println(err)
			}
		})
	}
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
