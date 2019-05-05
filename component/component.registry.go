package component

import (
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/context"
)

type Router struct {
	Name   string
	Method string
	Handle interface{}
}

type IServiceRegistry interface {
	API(pattern string, handle interface{})
	Cron(serverName string, timespan string, handle interface{})
	GetRouters(st string) []configs.Router
}

type ServiceRegistry struct {
	services map[string][]configs.Router
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{services: make(map[string][]configs.Router)}
}
func (s *ServiceRegistry) API(pattern string, r interface{}) {
	constructor, ok := r.(func(c IContainer) interface{})
	if !ok {
		panic("constructor is not func(container component.IContainer) interface{}")
	}
	constrObj := constructor(NewContainer())

	routers, ok := s.services[configs.ServerTypeAPI]
	if !ok {
		routers = []configs.Router{}
	}

	v := reflect.ValueOf(constrObj)
	t := reflect.TypeOf(constrObj)

	for i := 0; i < t.NumMethod(); i++ {
		switch t.Method(i).Name {
		case "Handler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			routers = append(routers, configs.Router{Name: pattern, Method: configs.HttpMethodALL, Handler: h})
		case "GetHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			routers = append(routers, configs.Router{Name: pattern, Method: configs.HttpMethodGet, Handler: h})
		case "PostHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			routers = append(routers, configs.Router{Name: pattern, Method: configs.HttpMethodPost, Handler: h})
		case "PutHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			routers = append(routers, configs.Router{Name: pattern, Method: configs.HttpMethodPut, Handler: h})
		case "DeleteHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			routers = append(routers, configs.Router{Name: pattern, Method: configs.HttpMethodDelete, Handler: h})
		}
	}

	s.services[configs.ServerTypeAPI] = routers
}

func (s *ServiceRegistry) Cron(sn string, tn string, r interface{}) {
	constructor, ok := r.(func(c IContainer) interface{})
	if !ok {
		panic("constructor is not func(container component.IContainer) interface{}")
	}
	constrObj := constructor(NewContainer())

	routers, ok := s.services[configs.ServerTypeCron]
	if !ok {
		routers = []configs.Router{}
	}

	v := reflect.ValueOf(constrObj)
	t := reflect.TypeOf(constrObj)
	for i := 0; i < t.NumMethod(); i++ {
		if strings.EqualFold(t.Method(i).Name, "Handler") {
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			routers = append(routers, configs.Router{Name: sn, Cron: tn, Handler: h})
			break
		}
	}
	s.services[configs.ServerTypeCron] = routers
}

func (s *ServiceRegistry) GetRouters(st string) []configs.Router {
	return s.services[st]
}
