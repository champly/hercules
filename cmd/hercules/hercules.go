package hercules

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/champly/hercules/component"
	"github.com/champly/hercules/configs"
	_ "github.com/champly/hercules/init"
	"github.com/champly/hercules/servers"
)

type Hercules struct {
	*option
	component.IServiceRegistry
	cl       chan bool
	services map[string]servers.IServers
}

func New(opts ...Option) *Hercules {
	h := &Hercules{
		option:           &option{ServerConfig: &configs.ServerConfig{Mode: "debug"}},
		IServiceRegistry: component.NewServiceRegistry(),
		cl:               make(chan bool),
		services:         map[string]servers.IServers{},
	}
	for _, opt := range opts {
		opt(h.option)
	}
	return h
}

func (h *Hercules) Start() {
	for _, t := range h.ServiceType {
		server, err := servers.NewRegistryServer(t, h.ServerConfig, h.GetRouters(t))
		if err != nil {
			panic(err)
		}
		if err = server.Start(); err != nil {
			fmt.Println(err)
			fmt.Println(t + "启动失败")
			return
		}
		fmt.Println(t + "启动成功")
		h.services[t] = server
	}

	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-sign:
		fmt.Println("正在关闭服务器")
		h.ShutDown()
	}

	fmt.Println("关闭成功")
}

func (h *Hercules) ShutDown() {
	go func() {
		for _, server := range h.services {
			server.ShutDown()
		}
		close(h.cl)
	}()

	select {
	case <-time.After(time.Second * 30):
		return
	case <-h.cl:
		return
	}
}
