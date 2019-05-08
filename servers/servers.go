package servers

import (
	"errors"
	"fmt"

	"github.com/champly/hercules/ctxs"
)

var resolvers = make(map[string]IServersResolver)

type IServersResolver interface {
	Resolve(router []ctxs.Router, handing func(*ctxs.Context) error) (IServers, error)
}

type IServers interface {
	Start() error
	Restart()
	ShutDown()
}

func Registry(serverType string, resolver IServersResolver) {
	if _, ok := resolvers[serverType]; ok {
		panic("services: " + serverType + " is repeat registry")
	}
	resolvers[serverType] = resolver
}

func NewRegistryServer(serverType string, router []ctxs.Router, handing func(*ctxs.Context) error) (IServers, error) {
	if resolver, ok := resolvers[serverType]; ok {
		return resolver.Resolve(router, handing)
	}
	return nil, errors.New(fmt.Sprintf("server type is not support:%s", serverType))
}
