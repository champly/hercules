package cron

import (
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

type cronServerAdapter struct{}

func (c *cronServerAdapter) Resolve(sConf *configs.ServerConfig, routers []configs.Router) (servers.IServers, error) {
	return NewCronServer(sConf, routers)
}

func init() {
	servers.Registry("cron", &cronServerAdapter{})
}
