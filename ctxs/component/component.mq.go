package component

import (
	"fmt"
	"sync"

	"github.com/champly/hercules/configs"
	"github.com/go-redis/redis"
)

type IComponentMQ interface {
	Produce(queueName, value string) error
}

type ComponentMQ struct {
	client *redis.Client
}

var (
	componentMQ *ComponentMQ
	lockMQ      sync.Mutex
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

func (m *ComponentMQ) getClient() {
	if m.client != nil {
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     configs.MQServer.Addr,
		Password: configs.MQServer.Password,
		DB:       configs.MQServer.DB,
	})
	// secret auth
	if configs.MQServer.Auth != "" {
		err := client.Do("AUTH", configs.MQServer.Auth).Err()
		if err != nil {
			panic("config component mq do auth failed:" + err.Error())
		}
	}
	_, err := client.Ping().Result()
	if err != nil {
		panic("config mqserver reture err:" + err.Error())
	}
	m.client = client
}

func (m *ComponentMQ) Produce(queueName, value string) error {
	m.getClient()

	cmd := m.client.LPush(queueName, value)
	_, err := cmd.Result()
	if err != nil {
		return fmt.Errorf("lpush %s %s fail:err:%+v", queueName, value, err)
	}
	return nil
}
