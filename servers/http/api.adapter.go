package http

import (
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

type apiServerAdapter struct{}

func (a *apiServerAdapter) Resolve(sConf *configs.ServerConfig, router []configs.Router) (servers.IServers, error) {
	return NewApiServer(sConf, router)
}

func init() {
	servers.Registry("api", &apiServerAdapter{})
}
