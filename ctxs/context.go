package ctxs

import (
	"sync"

	"github.com/champly/hercules/component"
	"github.com/gin-gonic/gin"
)

type Context struct {
	*herContext
	Log     ILog
	ToolBox component.IToolBox
}

var contextPool *sync.Pool

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Context{
				ToolBox: component.NewToolBox(),
			}
		},
	}
}

func GetContext(c *gin.Context) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.herContext = NewHerContext(c)
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
