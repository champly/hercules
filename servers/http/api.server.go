package http

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
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
		fmt.Println("cors enable")
	}
	return a, nil
}

func (a *ApiServer) getRouter(routers []servers.Router) error {
	for _, r := range routers {
		handler, ok := r.Handler.(func(*ctxs.Context) error)
		if !ok {
			v := reflect.TypeOf(r.Handler)
			panic(v.Elem().Name() + " handler is not func(ctx *ctxs.Context)error")
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

		h := a.GetRouter(c.Request.URL.String(), c.Request.Method)
		if h == nil {
			c.Status(http.StatusNotFound)
			return
		}
		ctx := ctxs.GetContext(c)
		defer ctx.Put()

		if a.handing != nil {
			if err := a.handing(ctx); err != nil {
				return
			}
		}
		if err := h(ctx); err != nil {
			fmt.Println(err)
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
	if !strings.EqualFold(configs.HttpServerInfo.Status, "start") {
		return fmt.Errorf("http server config is: %s", configs.HttpServerInfo.Status)
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}

func (a *ApiServer) ShutDown() {
	a.server.Shutdown(context.TODO())
	fmt.Println("http shutdown")
}

func (a *ApiServer) Restart() {
	fmt.Println("restart")
}

func (a *ApiServer) SetAddr(addr string) {
	a.server.Addr = addr
}
