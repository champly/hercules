package servers

import (
	"github.com/champly/hercules/context"
)

type IHandle interface {
	Handle(ctx *context.Context) (err error)
}
