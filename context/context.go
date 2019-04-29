package context

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type Context struct {
	// Request  *Request
	// Response *Response
	*gin.Context
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
	return ctx
}
