package servers

import (
	"errors"
)

var resolvers = make(map[string]IServersResolver)

type Router struct {
	Name    string
	Method  string
	Handler interface{}
}

type IServersResolver interface {
	Resolve(router []Router, handing interface{}) (IServers, error)
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

func NewRegistryServer(serverType string, router []Router, handing interface{}) (IServers, error) {
	if resolver, ok := resolvers[serverType]; ok {
		return resolver.Resolve(router, handing)
	}
	return nil, errors.New("server type is not support:" + serverType)
}
