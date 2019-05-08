package cron

import (
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
)

type cronServerAdapter struct{}

func (c *cronServerAdapter) Resolve(routers []ctxs.Router, handing func(*ctxs.Context) error) (servers.IServers, error) {
	return NewCronServer(routers, handing)
}

func init() {
	servers.Registry("cron", &cronServerAdapter{})
}
