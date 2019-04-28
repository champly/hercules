package component

import "github.com/champly/hercules/servers"

type IServiceRegistry interface {
}

type ServiceRegistry struct {
	services map[string]map[string]interface{}
	tags     map[string]interface{}
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{}
}

func (s *ServiceRegistry) API(pattern string, handle servers.IHandle) {
	if _, ok := s.tags[pattern]; ok {
		panic("router registry repeat:" + pattern)
	}
	s.tags[pattern] = handle
}

// func (s *ServiceRegistry) RAPI(pattern string, handle service.IHandle) {

// }
