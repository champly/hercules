package component

import (
	"github.com/champly/hercules/context"
)

type IComponent interface {
	Handle(c *context.Context) interface{}
	Close() error
}
