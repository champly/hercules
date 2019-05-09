package component

import (
	"reflect"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

func (s *ServiceRegistry) buildAPIRouterByObj(pattern string, constrObj interface{}) {
	routers, ok := s.services[configs.ServerTypeAPI]
	if !ok {
		routers = []servers.Router{}
	}

	v := reflect.ValueOf(constrObj)
	t := reflect.TypeOf(constrObj)

	for i := 0; i < t.NumMethod(); i++ {
		switch t.Method(i).Name {
		case "Handler":
			routers = append(routers, servers.Router{Name: pattern, Method: configs.HttpMethodALL, Handler: v.Method(i).Interface()})
		case "GetHandler":
			routers = append(routers, servers.Router{Name: pattern, Method: configs.HttpMethodGet, Handler: v.Method(i).Interface()})
		case "PostHandler":
			routers = append(routers, servers.Router{Name: pattern, Method: configs.HttpMethodPost, Handler: v.Method(i).Interface()})
		case "PutHandler":
			routers = append(routers, servers.Router{Name: pattern, Method: configs.HttpMethodPut, Handler: v.Method(i).Interface()})
		case "DeleteHandler":
			routers = append(routers, servers.Router{Name: pattern, Method: configs.HttpMethodDelete, Handler: v.Method(i).Interface()})
		}
	}
	s.services[configs.ServerTypeAPI] = routers
	return
}

func (s *ServiceRegistry) buildAPIRouterByFunc(pattern string, f interface{}) {
	routers, ok := s.services[configs.ServerTypeAPI]
	if !ok {
		routers = []servers.Router{}
	}
	routers = append(routers, servers.Router{Name: pattern, Method: configs.HttpMethodALL, Handler: f})
	s.services[configs.ServerTypeAPI] = routers
	return
}
