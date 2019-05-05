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
}

func NewCronServer(sConf *configs.ServerConfig, routers []configs.Router) (*CronServer, error) {
	c := &CronServer{ServerConfig: sConf}
	c.server = cron.New()
	c.AddFunc(routers)
	return c, nil
}

func (c *CronServer) AddFunc(routers []configs.Router) {
	for _, t := range routers {
		c.server.AddFunc(t.Cron, func() {
			ctx := context.GetDContext()
			defer ctx.Close()
			handler, _ := t.Handler.(func(*context.Context) error)
			if err := handler(ctx); err != nil {
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
