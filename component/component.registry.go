package component

import (
	"github.com/gin-gonic/gin"
)

type IServiceRegistry interface {
	API(pattern string, handle func(*gin.Context) error)
	GetRouter(router string) interface{}
	GetRouters() []Router
}

type ServiceRegistry struct {
	services map[string]map[string]interface{}
	tags     map[string]interface{}
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{tags: map[string]interface{}{}}
}

func (s *ServiceRegistry) API(pattern string, handle func(ctx *gin.Context) (err error)) {
	if _, ok := s.tags[pattern]; ok {
		panic("router registry repeat:" + pattern)
	}
	s.tags[pattern] = handle
}

func (s *ServiceRegistry) GetRouters() []Router {
	r := []Router{}
	for n, h := range s.tags {
		r = append(r, Router{n, h})
	}
	return r
}

func (s *ServiceRegistry) GetRouter(router string) interface{} {
	return s.tags[router]
}

// func (s *ServiceRegistry) RAPI(pattern string, handle service.IHandle) {

// }
