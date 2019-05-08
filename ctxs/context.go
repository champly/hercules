package ctxs

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Name    string
	Method  string
	ToolBox interface{}
	Handler Handler
}

type Handler func(*Context) error

type Context struct {
	*gin.Context
	Log     ILog
	ToolBox interface{}
}

var contextPool *sync.Pool

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}
}

func GetContext(c *gin.Context) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Context = c
	ctx.Log = NewLogger()
	return ctx
}

func GetDContext() *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Log = NewLogger()
	return ctx
}

func (c *Context) Put() {
	contextPool.Put(c)
}
