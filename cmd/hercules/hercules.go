package hercules

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/champly/hercules/cmd/hercules/status"
	"github.com/champly/hercules/component"
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	_ "github.com/champly/hercules/init"
	"github.com/champly/hercules/servers"
	"github.com/champly/hercules/servers/http"
)

type Hercules struct {
	*option
	component.IServiceRegistry
	component.IComponentDB
	cl       chan bool
	services map[string]servers.IServers
	handing  func(ctx *ctxs.Context) error
	initf    func(ctx *ctxs.Context) error
}

func New(opts ...Option) *Hercules {
	h := &Hercules{
		option:           &option{},
		IServiceRegistry: component.NewServiceRegistry(),
		IComponentDB:     component.NewComponentDB(),
		cl:               make(chan bool),
		services:         map[string]servers.IServers{},
	}
	for _, opt := range opts {
		opt(h.option)
	}
	return h
}

func (h *Hercules) Start() {
	// start services
	h.startService()

	// start health server
	if configs.SystemInfo.Health {
		h.startHealthService()
	}

	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-sign:
		h.ShutDown()
	}
	fmt.Println("关闭成功")
}

func (h *Hercules) startService() {
	for _, t := range strings.Split(configs.SystemInfo.Type, "|") {
		rl := h.GetRouters(t)
		if len(rl) == 0 {
			continue
		}
		server, err := servers.NewRegistryServer(t, rl, h.handing)
		if err != nil {
			panic(err)
		}
		if err = server.Start(); err != nil {
			fmt.Println(t+" start fail:", err)
			continue
		}
		fmt.Println(t + " start success")
		h.services[t] = server
	}
	if h.initf == nil {
		return
	}

	ctx := ctxs.GetDContext()
	if err := h.initf(ctx); err != nil {
		panic("Init service filed:" + err.Error())
	}
	ctx.Put()
}

func (h *Hercules) startHealthService() {
	routers := []servers.Router{}
	routers = append(routers, servers.Router{Name: "/status", Method: configs.HttpMethodALL, Handler: status.GetServerStatus})
	statusServer, err := http.NewApiServer(routers, nil)
	if err != nil {
		panic(err)
	}
	statusServer.SetAddr(":16666")
	if err = statusServer.Start(); err != nil {
		fmt.Println("health server start fail:", err)
		return
	}
	fmt.Println("health server start success")
}

func (h *Hercules) ShutDown() {
	fmt.Println("正在关闭服务器")
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

func (h *Hercules) Init(f func(*ctxs.Context) error) {
	h.initf = f
}

func (h *Hercules) Handing(f func(*ctxs.Context) error) {
	h.handing = f
}
