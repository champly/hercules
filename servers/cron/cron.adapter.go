package cron

import (
	"github.com/champly/hercules/servers"
)

type cronServerAdapter struct{}

func (c *cronServerAdapter) Resolve(routers []servers.Router, handing interface{}) (servers.IServers, error) {
	return NewCronServer(routers, handing)
}

func init() {
	servers.Registry("cron", &cronServerAdapter{})
}
