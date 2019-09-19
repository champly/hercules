package http

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

type Handler func(*ctxs.Context) error

type ApiServer struct {
	services map[string]map[string]func(*ctxs.Context) error
	server   *http.Server
	engine   *gin.Engine
	handing  func(*ctxs.Context) error
}

func NewApiServer(routers []servers.Router, h interface{}) (*ApiServer, error) {
	// nil not panic
	a := &ApiServer{services: make(map[string]map[string]func(*ctxs.Context) error)}
	if h != nil {
		handing, ok := h.(func(*ctxs.Context) error)
		if !ok {
			panic("handing function is not func(ctx *ctxs.Context)error")
		}
		a.handing = handing
	}
	a.server = &http.Server{
		Addr: configs.HttpServerInfo.Address,
	}
	a.server.Handler = a.getHandler(configs.SystemInfo.Mode)
	if err := a.getRouter(routers); err != nil {
		return nil, err
	}
	if configs.HttpServerInfo.Cors.Enable {
		color.HiYellow("cors enable")
	}
	return a, nil
}

func (a *ApiServer) getRouter(routers []servers.Router) error {
	for _, r := range routers {
		handler, ok := r.Handler.(func(*ctxs.Context) error)
		if !ok {
			if reflect.TypeOf(r.Handler).Kind() != reflect.Ptr {
				panic(reflect.TypeOf(r.Handler).Elem().Name() + " handler is not func(ctx *ctxs.Context)error")
			}
			panic(reflect.TypeOf(r.Handler).Name() + " handler is not func(ctx *ctxs.Context)error")
		}
		for _, m := range strings.Split(r.Method, "|") {
			if _, ok := a.services[r.Name]; !ok {
				a.services[r.Name] = map[string]func(*ctxs.Context) error{}
			}
			a.services[r.Name][m] = handler
		}
	}
	return nil
}

func (a *ApiServer) getHandler(mode string) http.Handler {
	gin.SetMode(mode)

	engine := gin.New()
	if strings.EqualFold(mode, "debug") {
		engine.Use(gin.Logger())
	}

	engine.Use(gin.Recovery())
	engine.Any("/*router", a.GeneralHandler())
	return engine
}

func (a *ApiServer) GeneralHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if configs.HttpServerInfo.Cors.Enable {
			for k, v := range configs.HttpServerInfo.Cors.Header {
				c.Writer.Header().Add(k, v)
			}
			if strings.EqualFold(c.Request.Method, configs.HttpMethodOptions) {
				return
			}
		}

		h := a.GetRouter(c.Request.URL.Path, c.Request.Method)
		if h == nil {
			c.Status(http.StatusNotFound)
			return
		}
		ctx := ctxs.GetContext(c)
		ctx.Type = ctxs.ServerTypeHTTP
		defer ctx.Put()

		if a.handing != nil {
			// handing
			if err := a.handing(ctx); err != nil {
				ctx.Log.Error(err.Error())

				if ctx.Writer.Status() != 200 {
					return
				}

				respMsg := "system busy"
				if strings.EqualFold(configs.SystemInfo.Mode, "debug") {
					respMsg = err.Error()
				}
				ctx.JSON(http.StatusInternalServerError, respMsg)
				return
			}
		}

		// Handler
		if err := h(ctx); err != nil {
			ctx.Log.Error(err.Error())

			if ctx.Writer.Status() != 200 {
				return
			}

			respMsg := "system busy"
			if strings.EqualFold(configs.SystemInfo.Mode, "debug") {
				respMsg = err.Error()
			}
			ctx.JSON(http.StatusInternalServerError, respMsg)
			return
		}
		return
	}
}

func (a *ApiServer) GetRouter(router string, method string) func(*ctxs.Context) error {
	method = strings.ToUpper(method)
	r, ok := a.services[router]
	if !ok {
		return nil
	}
	return r[method]
}

func (a *ApiServer) Start() error {
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if !strings.EqualFold(err.Error(), "http: Server closed") {
				color.HiRed(err.Error())
			}
		}
	}()
	return nil
}

func (a *ApiServer) ShutDown() {
	a.server.Shutdown(context.TODO())
	color.HiYellow("http shutdown")
}

func (a *ApiServer) Restart() {
	color.HiYellow("http restart")
}

func (a *ApiServer) SetAddr(addr string) {
	a.server.Addr = addr
}
