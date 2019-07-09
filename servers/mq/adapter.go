package mq

import (
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

type mqServerAdapter struct{}

func (a *mqServerAdapter) Resolve(router []servers.Router, handing interface{}) (servers.IServers, error) {
	return NewMQServer(router, handing)
}

func init() {
	servers.Registry(configs.ServerTypeMQ, &mqServerAdapter{})
}
