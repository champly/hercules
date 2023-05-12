package registry

import (
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

func (s *ServiceRegistry) buildCronRouterByObj(pattern string, constrObj interface{}) {
	routers, ok := s.services[configs.ServerTypeCron]
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
	s.services[configs.ServerTypeCron] = routers
}

func (s *ServiceRegistry) buildCronRouterByFunc(pattern string, f interface{}) {
	routers, ok := s.services[configs.ServerTypeCron]
	if !ok {
		routers = []servers.Router{}
	}
	routers = append(routers, servers.Router{Name: pattern, Handler: f})
	s.services[configs.ServerTypeCron] = routers
}
