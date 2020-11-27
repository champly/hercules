package hercules

import (
	"strings"
	"time"

	"github.com/champly/hercules/cmd/hercules/initialize"
	"github.com/champly/hercules/cmd/hercules/status"
	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	_ "github.com/champly/hercules/init"
	"github.com/champly/hercules/registry"
	"github.com/champly/hercules/servers"
	"github.com/champly/hercules/servers/http"
	"github.com/champly/lib4go/signal"
	"k8s.io/klog/v2"
)

// Hercules hercules struct
type Hercules struct {
	*option
	registry.IServiceRegistry
	// component.IComponentDB
	cl       chan bool
	services map[string]servers.IServers
	handing  func(ctx *ctxs.Context) error
	initf    func(ctx *ctxs.Context) error
}

// New build Hercules obj
func New(opts ...Option) *Hercules {
	h := &Hercules{
		option:           &option{},
		IServiceRegistry: registry.NewServiceRegistry(),
		// IComponentDB:     component.NewComponentDB(),
		cl:       make(chan bool),
		services: map[string]servers.IServers{},
	}
	for _, opt := range opts {
		opt(h.option)
	}
	return h
}

// Start start Hercules service
func (h *Hercules) Start() {
	initialize.Setpprof()

	// start services
	h.startService()

	// start health server
	if configs.SystemInfo.Health {
		h.startHealthService()
	}

	stop := signal.SetupSignalHandler()
	<-stop
	h.ShutDown()
	klog.Info("close success")
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
			klog.Errorf("[%s] start fail: %v", t, err)
			continue
		}
		klog.Infof("[%s] start success", t)
		h.services[t] = server
	}
	if h.initf == nil {
		return
	}

	ctx := ctxs.GetCronContext()
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
		klog.Errorf("health server start fail: %v", err)
		return
	}
	klog.Info("health server start success")
}

// ShutDown shutdown Hercules all service
func (h *Hercules) ShutDown() {
	klog.Info("Closing service, please wait a moment, if you need force close, prese 'Ctrl+c'")
	go func() {
		for _, server := range h.services {
			server.ShutDown()
		}
		close(h.cl)
	}()

	select {
	case <-time.After(time.Second * 30):
		klog.Error("Closing 30s timeout, force close!")
		return
	case <-h.cl:
		return
	}
}

// Init set Hercules init func
func (h *Hercules) Init(f func(*ctxs.Context) error) {
	h.initf = f
}

// Handing set all request first exec func
func (h *Hercules) Handing(f func(*ctxs.Context) error) {
	h.handing = f
}
