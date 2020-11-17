package registry

import (
	"reflect"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/servers"
)

const (
	DefaultHandlerName       = "Handler"
	RESTfulHandlerGetName    = "GetHandler"
	RESTfulHandlerPostName   = "PostHandler"
	RESTfulHandlerPutName    = "PutHandler"
	RESTfulHandlerDeleteName = "DeleteHandler"
)

var (
	methodMap map[string]string
)

func init() {
	methodMap = map[string]string{
		DefaultHandlerName:       configs.HttpMethodALL,
		RESTfulHandlerGetName:    configs.HttpMethodGet,
		RESTfulHandlerPostName:   configs.HttpMethodPost,
		RESTfulHandlerPutName:    configs.HttpMethodPut,
		RESTfulHandlerDeleteName: configs.HttpMethodDelete,
	}
}

func (s *ServiceRegistry) buildAPIRouterByObj(pattern string, constrObj interface{}) {
	routers, ok := s.services[configs.ServerTypeAPI]
	if !ok {
		routers = []servers.Router{}
	}

	v := reflect.ValueOf(constrObj)
	t := reflect.TypeOf(constrObj)

	for i := 0; i < t.NumMethod(); i++ {
		method, ok := methodMap[t.Method(i).Name]
		if ok {
			routers = append(routers, servers.Router{Name: pattern, Method: method, Handler: v.Method(i).Interface()})
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
