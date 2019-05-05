package http

import (
	ctxs "context"
	"fmt"
	"net/http"
	"strings"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/context"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	*configs.ServerConfig
	services map[string]map[string]interface{}
	server   *http.Server
	engine   *gin.Engine
}

func NewApiServer(sConf *configs.ServerConfig, routers []configs.Router) (*ApiServer, error) {
	a := &ApiServer{ServerConfig: sConf, services: make(map[string]map[string]interface{})}
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
				a.services[r.Name] = map[string]interface{}{}
				a.services[r.Name][m] = r.Handler

			}
		}
	}
	return nil
}

func (a *ApiServer) getHandler(mode string) http.Handler {
	gin.SetMode(mode)
	engine := gin.New()
	engine.Any("/*router", a.GeneralHandler())
	return engine
}

func (a *ApiServer) GeneralHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		hh := a.GetRouter(c.Request.URL.String(), c.Request.Method)
		if hh == nil {
			c.Status(http.StatusNotFound)
			return
		}
		ctx := context.GetContext(c)
		defer ctx.Close()
		handler, ok := hh.(func(*context.Context) error)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		if err := handler(ctx); err != nil {
			fmt.Println(err)
		}
	}
}

func (a *ApiServer) GetRouter(router string, method string) interface{} {
	method = strings.ToUpper(method)
	r, ok := a.services[router]
	if !ok {
		return nil
	}
	return r[method]
}

func (a *ApiServer) Start() error {
	go a.server.ListenAndServe()
	return nil
}

func (a *ApiServer) ShutDown() {
	a.server.Shutdown(ctxs.TODO())
	fmt.Println("http shutdown")
}

func (a *ApiServer) Restart() {
	fmt.Println("restart")
}
