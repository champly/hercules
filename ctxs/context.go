package ctxs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/champly/hercules/ctxs/component"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

const (
	ServerTypeHTTP = "http"
	ServerTypeCron = "cron"
	ServerTypeMQ   = "mq"
)

type Context struct {
	*herContext
	Log     ILog
	ToolBox component.IToolBox
	Type    string
}

var (
	contextPool *sync.Pool
	validate    *validator.Validate
)

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Context{
				ToolBox: component.NewToolBox(),
			}
		},
	}
	validate = validator.New()
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

func (c *Context) CheckStruct(data interface{}) error {
	return validate.Struct(data)
}
