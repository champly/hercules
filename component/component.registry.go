package component

import (
	"reflect"

	"github.com/champly/hercules/context"
)

type Router struct {
	Name   string
	Method string
	Handle interface{}
}

type IServiceRegistry interface {
	API(pattern string, handle interface{})
	GetRouter(router string, method string) interface{}
	GetRouters() []Router
}

type ServiceRegistry struct {
	services map[string]map[string]interface{}
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{services: map[string]map[string]interface{}{}}
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
			if _, ok := s.services[pattern]["get"]; !ok {
				s.services[pattern]["get"] = h
			}
			if _, ok := s.services[pattern]["post"]; !ok {
				s.services[pattern]["post"] = h
			}
			if _, ok := s.services[pattern]["put"]; !ok {
				s.services[pattern]["put"] = h
			}
			if _, ok := s.services[pattern]["delete"]; !ok {
				s.services[pattern]["delete"] = h
			}
		case "GetHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern]["get"] = h
		case "PostHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern]["post"] = h
		case "PutHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern]["put"] = h
		case "DeleteHandler":
			h, ok := v.Method(i).Interface().(func(*context.Context) error)
			if !ok {
				panic(t.Method(i).Name + " is not func(*context.Context) error method")
			}
			s.services[pattern]["delete"] = h
		}
	}
}

func (s *ServiceRegistry) GetRouters() []Router {
	r := []Router{}
	for p, ss := range s.services {
		for m, h := range ss {
			r = append(r, Router{p, m, h})
		}
	}
	return r
}

func (s *ServiceRegistry) GetRouter(router string, method string) interface{} {
	r, ok := s.services[router]
	if !ok {
		return nil
	}
	return r[method]
}
