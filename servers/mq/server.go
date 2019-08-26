package mq

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/champly/hercules/configs"
	"github.com/champly/hercules/ctxs"
	"github.com/champly/hercules/servers"
	"github.com/go-redis/redis"
)

type MQServer struct {
	client   *redis.Client
	handing  func(*ctxs.Context) error
	servers  map[string]func(*ctxs.Context) error
	stopCh   chan struct{}
	stopSucc chan struct{}
}

func NewMQServer(routers []servers.Router, h interface{}) (*MQServer, error) {
	handing, ok := h.(func(*ctxs.Context) error)
	if !ok {
		panic("handing function is not func(ctx *ctxs.Context) error")
	}
	mq := &MQServer{
		handing:  handing,
		servers:  make(map[string]func(*ctxs.Context) error),
		stopCh:   make(chan struct{}),
		stopSucc: make(chan struct{}),
	}

	mq.client = redis.NewClient(&redis.Options{
		Addr:     configs.MQServer.Addr,
		Password: configs.MQServer.Password,
		DB:       configs.MQServer.DB,
	})
	_, err := mq.client.Ping().Result()
	if err != nil {
		panic("config mqserver reture err:" + err.Error())
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
	wg := sync.WaitGroup{}
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
	for {
		msgCh := make(chan messgae)

		go func() {
			cmd := m.client.BRPop(time.Second*1, queueName)
			msg, err := cmd.Result()
			hasData := err == nil && len(msg) > 0
			ndata := ""
			if hasData {
				ndata = msg[len(msg)-1]
			}
			msgCh <- messgae{Data: ndata, HasData: hasData}
		}()

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

			go m.do(msg.Data, callback)
		}
	}
}

func (m *MQServer) do(data string, callback func(*ctxs.Context) error) (err error) {
	ctx := ctxs.GetMQContext(data)
	ctx.Type = ctxs.ServerTypeMQ
	defer ctx.Put()
	if m.handing != nil {
		if err = m.handing(ctx); err != nil {
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
	fmt.Println("mq shutdown")
}

func (m *MQServer) Restart() {
	fmt.Println("mq restart")
}
