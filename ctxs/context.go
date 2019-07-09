package ctxs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/champly/hercules/ctxs/component"
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

func GetCronContext() *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Log = NewLogger()
	return ctx
}

func GetMQContext(msg string) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.herContext = NewHerContext(&gin.Context{
		Request: &http.Request{
			Body: ioutil.NopCloser(bytes.NewBuffer([]byte(msg))),
		},
	})
	ctx.Log = NewLogger()
	return ctx
}

func (c *Context) Put() {
	contextPool.Put(c)
}
