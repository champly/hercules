package http

import (
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

type apiServerAdapter struct{}

func (a *apiServerAdapter) Resolve(router []servers.Router, handing interface{}) (servers.IServers, error) {
	return NewApiServer(router, handing)
}

func init() {
	servers.Registry(configs.ServerTypeAPI, &apiServerAdapter{})
}
