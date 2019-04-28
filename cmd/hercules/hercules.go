package hercules

import (
	"fmt"
	"net/http"

	"github.com/champly/hercules/component"
	"github.com/gin-gonic/gin"
)

type Hercules struct {
	*option
	server *http.Server
	engine *gin.Engine
	component.IServiceRegistry
}

func New(opts ...Option) *Hercules {
	h := &Hercules{option: &option{Mode: "debug"}, IServiceRegistry: component.NewServiceRegistry()}
	for _, opt := range opts {
		opt(h.option)
	}

	h.server = &http.Server{
		Addr: h.Addr,
	}
	h.server.Handler = h.getHandler(h.Mode)
	return h
}

func (h *Hercules) getHandler(mode string) http.Handler {
	gin.SetMode(mode)
	engine := gin.New()
	engine.Any("/*name", h.GeneralHandler())
	return engine
}

func (h *Hercules) GeneralHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hh := h.GetRouter(ctx.Request.URL.String())
		if hh == nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		handler, ok := hh.(func(*gin.Context) error)
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		if err := handler(ctx); err != nil {
			fmt.Println(err)
		}
	}
}

func (h *Hercules) Start() {
	h.server.ListenAndServe()
}
