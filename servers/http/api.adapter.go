package http

import (
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
)

type apiServerAdapter struct{}

func (a *apiServerAdapter) Resolve(router []ctxs.Router) (servers.IServers, error) {
	return NewApiServer(router)
}

func init() {
	servers.Registry("api", &apiServerAdapter{})
}
