package nsqhelper

import (
	"config"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	nsq "github.com/EncoreJiang/go-nsq"
	"github.com/cihub/seelog"
	// consul "github.com/hashicorp/consul/api"
)

var nsqProducer *nsq.Producer
var nsqProducerMutex = sync.RWMutex{}
var nsqProducerFail = true

func PublishMessage(topic string, payload []byte) error {
	var err error
	for retry := 0; retry < 5; retry++ {
		func() {
			defer func() {
				if rerr := recover(); rerr != nil {
					if rerr, ok := rerr.(error); ok {
						err = rerr
						return
					}
					err = errors.New(fmt.Sprint(rerr))
				}
				err = GetNSQProducer().Publish(topic, payload)
			}()
		}()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration((1 + retry)))
	}
	seelog.Error("publish message failed:", err)
	return err
}

func GetNSQProducer() *nsq.Producer {
	nsqProducerMutex.RLock()
	if nsqProducer != nil {
		if nsqProducerFail {
			if err := nsqProducer.Ping(); err != nil {
				seelog.Errorf("Ping nsqd server %v failed:%v", nsqProducer.String(), err)
			} else {
				nsqProducerFail = false
			}
		}
		if nsqProducerFail == false {
			nsqProducerMutex.RUnlock()
			return nsqProducer
		}
	}
	nsqProducerMutex.RUnlock()
	nsqProducerMutex.Lock()
	defer nsqProducerMutex.Unlock()

	nsqConfig := nsq.NewConfig()

	var addr string
	// if len(config.NSQDServiceName) != 0 {
	// 	consulConfig := consul.DefaultConfig()
	// 	consulClient, err := consul.NewClient(consulConfig)
	// 	if err != nil {
	// 		seelog.Error("consul.NewClient failed:", err)
	// 		panic(err)
	// 	}

	// 	q := &consul.QueryOptions{}
	// 	services, _, err := consulClient.Health().Service(config.NSQDServiceName, "", true, q)
	// 	if err != nil {
	// 		seelog.Error("query catalog failed:", err)
	// 		panic(err)
	// 	}
	// 	if len(services) > 0 {
	// 		service := services[rand.Intn(len(services))]
	// 		if len(service.Service.Address) != 0 {
	// 			addr = fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
	// 		} else {
	// 			addr = fmt.Sprintf("%v:%v", service.Node.Address, service.Service.Port)
	// 		}
	// 	}
	// }
	if len(addr) == 0 {
		addrs := strings.Split(config.NSQDAddress, ";")
		if len(addrs) != 0 {
			addr = addrs[rand.Intn(len(addrs))]
		}
	}
	nsqProducer, err := nsq.NewProducer(addr, nsqConfig)

	if err != nil {
		seelog.Error("get new produce failed:", err)
		panic("Could not connect to nsq")
	}
	seelog.Info("get new produce ok:", addr)
	return nsqProducer
}
