package http

import (
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
)

type apiServerAdapter struct{}

func (a *apiServerAdapter) Resolve(router []ctxs.Router, handing func(*ctxs.Context) error) (servers.IServers, error) {
	return NewApiServer(router, handing)
}

func init() {
	servers.Registry("api", &apiServerAdapter{})
}
