package component

import (
	"fmt"
	"sync"

	"github.com/champly/hercules/configs"
	"github.com/go-redis/redis"
	"k8s.io/klog/v2"
)

type IComponentMQ interface {
	Produce(queueName, value string) error
}

type ComponentMQ struct {
}

var (
	componentMQ *ComponentMQ
	lockMQ      sync.Mutex
	client      *redis.Client
	l           sync.Mutex
)

func NewComponentMQ() *ComponentMQ {
	if componentMQ != nil {
		return componentMQ
	}
	lockMQ.Lock()
	defer lockMQ.Unlock()

	if componentMQ != nil {
		return componentMQ
	}
	componentMQ = &ComponentMQ{}
	return componentMQ
}

func (m *ComponentMQ) Produce(queueName, value string) error {
	if client == nil {
		GetSingleClient()
	}

	cmd := client.LPush(queueName, value)
	_, err := cmd.Result()
	if err != nil {
		return fmt.Errorf("lpush %s %s fail:err:%+v", queueName, value, err)
	}
	return nil
}

func GetSingleClient() *redis.Client {
	if client != nil {
		return client
	}

	l.Lock()
	defer l.Unlock()

	if client != nil {
		return client
	}
	cli := redis.NewClient(&redis.Options{
		Addr:     configs.MQServer.Addr,
		Password: configs.MQServer.Password,
		DB:       configs.MQServer.DB,
	})
	// secret auth
	if configs.MQServer.Auth != "" {
		err := cli.Do("AUTH", configs.MQServer.Auth).Err()
		if err != nil {
			panic("config component mq do auth failed:" + err.Error())
		}
	}
	_, err := cli.Ping().Result()
	if err != nil {
		panic("config mqserver reture err:" + err.Error())
	}

	klog.Infof("connect redis {%s} succ.", configs.MQServer.Addr)

	client = cli
	return client
}
