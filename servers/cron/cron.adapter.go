package cron

import (
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

type cronServerAdapter struct{}

func (c *cronServerAdapter) Resolve(routers []configs.Router) (servers.IServers, error) {
	return NewCronServer(routers)
}

func init() {
	servers.Registry("cron", &cronServerAdapter{})
}
