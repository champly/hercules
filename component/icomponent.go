package component

import (
	"github.com/champly/hercules/context"
)

type IComponent interface {
	GetRouter() []Router
	Close() error
}

type IHandler interface {
	Handler(ctx *context.Context) error
}

type IGetHandler interface {
	GetHandler(ctx *context.Context) error
}

type IPostHandler interface {
	PostHandler(ctx *context.Context) error
}

type IPutHandler interface {
	PutHandler(ctx *context.Context) error
}

type IDeleteHandler interface {
	DeleteHandler(ctx *context.Context) error
}
