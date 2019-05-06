package component

import (
	"github.com/champly/hercules/ctxs"
)

type IComponent interface {
	GetRouter() []Router
	Close() error
}

type IHandler interface {
	Handler(ctx *ctxs.Context) error
}

type IGetHandler interface {
	GetHandler(ctx *ctxs.Context) error
}

type IPostHandler interface {
	PostHandler(ctx *ctxs.Context) error
}

type IPutHandler interface {
	PutHandler(ctx *ctxs.Context) error
}

type IDeleteHandler interface {
	DeleteHandler(ctx *ctxs.Context) error
}
