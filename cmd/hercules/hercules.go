package hercules

import (
	"fmt"

	"github.com/champly/hercules/component"
	"github.com/champly/hercules/configs"
	_ "github.com/champly/hercules/init"
	"github.com/champly/hercules/servers"
)

type Hercules struct {
	*option
	component.IServiceRegistry
}

func New(opts ...Option) *Hercules {
	h := &Hercules{option: &option{ServerConfig: &configs.ServerConfig{Mode: "debug"}}, IServiceRegistry: component.NewServiceRegistry()}
	for _, opt := range opts {
		opt(h.option)
	}
	return h
}

func (h *Hercules) Start() {
	for _, t := range h.ServiceType {
		server, err := servers.NewRegistryServer(t, h.ServerConfig, h.GetRouters())
		if err != nil {
			panic(err)
		}
		if err = server.Start(); err != nil {
			fmt.Println(err)
			fmt.Println("启动失败")
			return
		}
		fmt.Println("启动成功")
	}
}
