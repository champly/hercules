package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	*configs.ServerConfig
	services map[string]map[string]ctxs.Handler
	server   *http.Server
	engine   *gin.Engine
}

func NewApiServer(sConf *configs.ServerConfig, routers []configs.Router) (*ApiServer, error) {
	a := &ApiServer{ServerConfig: sConf, services: make(map[string]map[string]ctxs.Handler)}
	a.server = &http.Server{
		Addr: a.Addr,
	}
	a.server.Handler = a.getHandler(a.Mode)
	if err := a.getRouter(routers); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *ApiServer) getRouter(routers []configs.Router) error {
	for _, r := range routers {
		for _, m := range strings.Split(r.Method, "|") {
			if _, ok := a.services[r.Name]; !ok {
				a.services[r.Name] = map[string]ctxs.Handler{}
			}
			a.services[r.Name][m] = r.Handler
		}
	}
	return nil
}

func (a *ApiServer) getHandler(mode string) http.Handler {
	gin.SetMode(mode)
	engine := gin.New()
	if mode == "debug" {
		engine.Use(gin.Logger())
	}
	engine.Use(gin.Recovery())
	engine.Any("/*router", a.GeneralHandler())
	return engine
}

func (a *ApiServer) GeneralHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "http://localhost:8081")
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "__jwt__")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		c.Writer.Header().Add("Access-Control-Expose-Headers", "__jwt__")

		if strings.EqualFold(c.Request.Method, configs.HttpMethodOptions) {
			return
		}

		handler := a.GetRouter(c.Request.URL.String(), c.Request.Method)
		if handler == nil {
			c.Status(http.StatusNotFound)
			return
		}
		ctx := ctxs.GetContext(c)
		defer ctx.Put()

		if err := handler(ctx); err != nil {
			fmt.Println(err)
		}
		return
	}
}

func (a *ApiServer) GetRouter(router string, method string) ctxs.Handler {
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
