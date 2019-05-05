package component

import (
	"reflect"

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
	GetRouters() map[string]map[string]interface{}
}

type ServiceRegistry struct {
	services map[string]map[string]interface{}
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{services: make(map[string]map[string]interface{})}
}

func (s *ServiceRegistry) API(pattern string, r interface{}) {
	constructor, ok := r.(func(c IContainer) interface{})
	if !ok {
		panic("constructor is not func(container component.IContainer) interface{}")
	}
	constrObj := constructor(NewContainer())

	if _, ok := s.services[pattern]; ok {
		panic("router registry repeat:" + pattern)
	}
	s.services[pattern] = map[string]interface{}{}

	v := reflect.ValueOf(constrObj)
	t := reflect.TypeOf(constrObj)
	for i := 0; i < t.NumMethod(); i++ {
		switch t.Method(i).Name {
		case "Handler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			if _, ok := s.services[pattern][configs.HttpMethodGet]; !ok {
				s.services[pattern][configs.HttpMethodGet] = h
			}
			if _, ok := s.services[pattern][configs.HttpMethodPost]; !ok {
				s.services[pattern][configs.HttpMethodPost] = h
			}
			if _, ok := s.services[pattern][configs.HttpMethodPut]; !ok {
				s.services[pattern][configs.HttpMethodPut] = h
			}
			if _, ok := s.services[pattern][configs.HttpMethodDelete]; !ok {
				s.services[pattern][configs.HttpMethodDelete] = h
			}
		case "GetHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern][configs.HttpMethodGet] = h
		case "PostHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern][configs.HttpMethodPost] = h
		case "PutHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern][configs.HttpMethodPut] = h
		case "DeleteHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern][configs.HttpMethodDelete] = h
		}
	}
}

func (s *ServiceRegistry) GetRouters() map[string]map[string]interface{} {
	return s.services
}
