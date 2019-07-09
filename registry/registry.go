package registry

import (
	"github.com/champly/hercules/ctxs/component"
	"github.com/champly/hercules/servers"
)

type IServiceRegistry interface {
	API(pattern string, handle interface{})
	Cron(name string, handle interface{})
	GetRouters(st string) []servers.Router
}

type ServiceRegistry struct {
	services map[string][]servers.Router
	toolBox  component.IToolBox
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{services: make(map[string][]servers.Router), toolBox: component.NewToolBox()}
}

func (s *ServiceRegistry) API(pattern string, r interface{}) {
	constructor, ok := r.(func() interface{})
	if ok {
		s.buildAPIRouterByObj(pattern, constructor())
		return
	}

	s.buildAPIRouterByFunc(pattern, r)
	return
}

func (s *ServiceRegistry) Cron(pattern string, r interface{}) {
	constructor, ok := r.(func() interface{})
	if ok {
		s.buildCronRouterByObj(pattern, constructor())
		return
	}

	s.buildCronRouterByFunc(pattern, r)
	return
}

func (s *ServiceRegistry) GetRouters(st string) []servers.Router {
	return s.services[st]
}
