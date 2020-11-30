package mq

import (
	"context"
	"reflect"
	"runtime"
	ssync "sync"
	"time"

	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/ctxs/component"
	"github.com/champly/hercules/servers"
	"github.com/champly/lib4go/sync"
	"github.com/go-redis/redis/v8"
	"k8s.io/klog/v2"
)

type MQServer struct {
	client     *redis.Client
	preHand    func(*ctxs.Context) error
	servers    map[string]func(*ctxs.Context) error
	stopCh     chan struct{}
	stopSucc   chan struct{}
	workerPool sync.WorkerPool
}

func NewMQServer(routers []servers.Router, h interface{}) (*MQServer, error) {
	preHand, ok := h.(func(*ctxs.Context) error)
	if !ok {
		panic("preHand function is not func(ctx *ctxs.Context) error")
	}
	mq := &MQServer{
		preHand:    preHand,
		servers:    make(map[string]func(*ctxs.Context) error),
		stopCh:     make(chan struct{}),
		stopSucc:   make(chan struct{}),
		client:     component.GetSingleClient(),
		workerPool: sync.NewWorkerPool(runtime.GOMAXPROCS(0) * 10),
	}

	mq.getRouter(routers)
	return mq, nil
}

func (m *MQServer) getRouter(routers []servers.Router) {
	for _, r := range routers {
		handler, ok := r.Handler.(func(*ctxs.Context) error)
		if !ok {
			if reflect.TypeOf(r.Handler).Kind() != reflect.Ptr {
				panic(reflect.TypeOf(r.Handler).Elem().Name() + " handler is not func(ctx *ctxs.Context)error")
			}
			panic(reflect.TypeOf(r.Handler).Name() + " handler is not func(ctx *ctxs.Context)error")
		}
		m.servers[r.Name] = handler
	}
}

func (m *MQServer) startServer() {
	wg := ssync.WaitGroup{}
	wg.Add(len(m.servers))
	for queue, handler := range m.servers {
		go func(queue string, handler func(*ctxs.Context) error) {
			m.Consume(queue, handler)
			wg.Done()
		}(queue, handler)
	}
	wg.Wait()
	close(m.stopSucc)
}

func (m *MQServer) Consume(queueName string, callback func(*ctxs.Context) error) {
	msgCh := make(chan messgae)
	defer close(msgCh)

	for {
		m.workerPool.ScheduleAuto(func() {
			cmd := m.client.BRPop(context.TODO(), time.Second*1, queueName)
			msg, err := cmd.Result()
			hasData := err == nil && len(msg) > 0
			ndata := ""
			if hasData {
				ndata = msg[len(msg)-1]
			}
			msgCh <- messgae{Data: ndata, HasData: hasData}
		})

		select {
		case <-m.stopCh:
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			if !msg.HasData {
				continue
			}
			m.workerPool.ScheduleAuto(func() {
				m.do(msg.Data, callback)
			})

		}
	}
}

func (m *MQServer) do(data string, callback func(*ctxs.Context) error) (err error) {
	ctx := ctxs.GetMQContext(data)
	ctx.Type = ctxs.ServerTypeMQ
	defer ctx.Put()

	if m.preHand != nil {
		if err = m.preHand(ctx); err != nil {
			return
		}
	}
	if err = callback(ctx); err != nil {
		ctx.Log.Error(err)
	}
	return err
}

func (m *MQServer) Start() error {
	go m.startServer()
	return nil
}

func (m *MQServer) ShutDown() {
	close(m.stopCh)
	m.client.Close()
	<-m.stopSucc
	klog.Info("mq shutdown")
}

func (m *MQServer) Restart() {
	klog.Info("mq restart")
}
