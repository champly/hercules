package cron

import (
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
)

type cronServerAdapter struct{}

func (c *cronServerAdapter) Resolve(routers []ctxs.Router) (servers.IServers, error) {
	return NewCronServer(routers)
}

func init() {
	servers.Registry("cron", &cronServerAdapter{})
}
