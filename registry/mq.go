package registry

import (
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

func (s *ServiceRegistry) buildMQRouterByObj(pattern string, constrObj interface{}) {
	routers, ok := s.services[configs.ServerTypeMQ]
	if !ok {
		routers = []servers.Router{}
	}

	v := reflect.ValueOf(constrObj)
	t := reflect.TypeOf(constrObj)

	for i := 0; i < t.NumMethod(); i++ {
		if strings.EqualFold(t.Method(i).Name, "Handler") {
			routers = append(routers, servers.Router{Name: pattern, Handler: v.Method(i).Interface()})
			break
		}
	}
	s.services[configs.ServerTypeMQ] = routers
	return
}

func (s *ServiceRegistry) buildMQRouterByFunc(pattern string, f interface{}) {
	routers, ok := s.services[configs.ServerTypeMQ]
	if !ok {
		routers = []servers.Router{}
	}
	routers = append(routers, servers.Router{Name: pattern, Handler: f})
	s.services[configs.ServerTypeMQ] = routers
	return
}
