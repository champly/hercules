package servers

import "github.com/champly/hercules/context"

type IServersAdapter interface {
	Resolve()
}

type IHandlerAdapter interface {
	Handler(ctx *context.Context) error
}

type IServers interface {
	Start()
	Restart()
	ShutDown()
}
