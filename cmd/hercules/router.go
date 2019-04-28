package hercules

import "github.com/gin-gonic/gin"

type Router struct {
	pattern string
	handle  func(ctx *gin.Context) error
}

type IHandleFunc interface {
	Handle(ctx *gin.Context) error
}

func (h *Hercules) HandleFunc(pattern string, handle IHandleFunc) {
	h.routers = append(h.routers, Router{pattern, handle.Handle})
}

func errHandleFunc(handle func(ctx *gin.Context) error) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		err := handle(ctx)
		if err != nil {
			ctx.Status(500)
		}
		return
	}
}
