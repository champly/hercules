package hercules

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Hercules struct {
	*option
	server  *http.Server
	engine  *gin.Engine
	routers []Router
}

func New(opts ...Option) *Hercules {
	h := &Hercules{option: &option{Mode: "debug"}, routers: []Router{}}
	for _, opt := range opts {
		opt(h.option)
	}

	h.server = &http.Server{
		Addr: h.Addr,
	}

	// h.server.Handler = getHandler(h.Mode, h.routers)
	return h
}

func getHandler(mode string, routes []Router) http.Handler {
	gin.SetMode(mode)
	engine := gin.New()
	for _, r := range routes {
		engine.Any(r.pattern, errHandleFunc(r.handle))
	}
	return engine
}

func (h *Hercules) Start() {
	h.server.Handler = getHandler(h.Mode, h.routers)
	h.server.ListenAndServe()
}
