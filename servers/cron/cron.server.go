package cron

import (
	"fmt"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/context"
	"github.com/robfig/cron"
)

type CronServer struct {
	*configs.ServerConfig
	server   *cron.Cron
	schedule cron.Schedule
	routers  []configs.Router
}

func NewCronServer(sConf *configs.ServerConfig, routers []configs.Router) (*CronServer, error) {
	c := &CronServer{ServerConfig: sConf, routers: routers}
	c.server = cron.New()
	c.AddFunc()
	return c, nil
}

func (c *CronServer) AddFunc() {
	for _, t := range c.routers {
		c.server.AddFunc(t.Cron, func() {
			ctx := context.GetDContext()
			defer ctx.Close()
			if err := t.Handler(ctx); err != nil {
				fmt.Println(err)
			}
		})
	}
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
