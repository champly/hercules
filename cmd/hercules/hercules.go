package hercules

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/champly/hercules/component"
	"github.com/champly/hercules/context"
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
	return func(c *gin.Context) {
		hh := h.GetRouter(c.Request.URL.String(), strings.ToLower(c.Request.Method))
		if hh == nil {
			c.Status(http.StatusNotFound)
			return
		}
		ctx := context.GetContext(c)
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

func (h *Hercules) Start() {
	h.server.ListenAndServe()
}
