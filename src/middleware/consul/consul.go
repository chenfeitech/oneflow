package utils

import (
	"conv"
	"fmt"
	"net"
	"strings"
	"time"

	"net/http"

	"github.com/cihub/seelog"
	consul "github.com/hashicorp/consul/api"
	"github.com/miekg/dns"
)

type consulServiceDialer struct {
	resolver net.Addr
	client   *dns.Client
	dailer   *net.Dialer
}

// defaultServiceDialer use for dial to consul service
var defaultServiceDialer *consulServiceDialer

// ConsulServiceClient use for request to consul service
var ConsulServiceClient *http.Client

func init() {

	defaultServiceDialer = &consulServiceDialer{
		&net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 8600,
		},
		new(dns.Client),
		&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	}

	ConsulServiceClient = &http.Client{
		Transport: &http.Transport{
			Dial: defaultServiceDialer.Dial,
		},
	}
}

func (dd *consulServiceDialer) Dial(network, address string) (net.Conn, error) {
	tokens := strings.Split(address, ":")
	hostname := tokens[0]
	answer, extra, err := dd.getAnswer(hostname)
	if len(answer) == 0 {
		return nil, fmt.Errorf("DNS lookup failed for %s", hostname)
	}
	var conn net.Conn
	for i, ans := range answer {
		if len(extra) > i {
			if a, ok := extra[i].(*dns.A); ok {
				tcpAddr := &net.TCPAddr{IP: a.A, Port: int(ans.(*dns.SRV).Port)}
				conn, err = dd.dailer.Dial(network, tcpAddr.String())
				if err == nil {
					return conn, nil
				}
			}
		}
	}
	return nil, err
}

func (dd *consulServiceDialer) getAnswer(hostname string) ([]dns.RR, []dns.RR, error) {
	query := new(dns.Msg).SetQuestion(hostname+".", dns.TypeSRV)
	resp, _, err := dd.client.Exchange(query, dd.resolver.String())
	if err != nil {
		return nil, nil, err
	}
	return resp.Answer, resp.Extra, nil
}

func ConsulGetAnswer(hostname string) (ip string, port string, err error) {
	answer, extra, err := defaultServiceDialer.getAnswer(hostname)
	for i, ans := range answer {
		if len(extra) > i {
			if a, ok := extra[i].(*dns.A); ok {
				return a.A.String(), conv.String(ans.(*dns.SRV).Port), nil
			}
		}
	}
	return "", "", seelog.Error("DNS lookup failed for %s", hostname)
}

const (
	Productive         = "prod"
	ProductiveStandby  = "prod_standby"
	Development        = "dev"
	DevelopmentStandby = "dev_standby"
)

type ConsulLeader struct {
	Client    *consul.Client
	sessionId string
	name      string
	tag       string
	port      int
}

func CreateConsulLeader(name string, tag string, port int) (*ConsulLeader, error) {
	config := consul.DefaultConfig()
	client, err := consul.NewClient(config)
	return &ConsulLeader{Client: client, name: name, port: port, tag: tag}, err
}

//注册服务
func (c *ConsulLeader) ServiceRegister(tag string) error {
	check := &consul.AgentServiceCheck{
		TCP: fmt.Sprintf("127.0.0.1:%d", c.port),
		DeregisterCriticalServiceAfter: "1m",
		Interval:                       "10s",
	}
	reg := &consul.AgentServiceRegistration{
		Name:              c.name,
		Tags:              []string{tag},
		Port:              c.port,
		EnableTagOverride: true,
		Check:             check,
	}
	return c.Client.Agent().ServiceRegister(reg)
}

//服务查询
func (c *ConsulLeader) ServiceConfirm(tag string) (bool, error) {
	services, err := c.Client.Agent().Services()
	if err != nil {
		return false, err
	}
	service, ok := services[c.name]
	if ok {
		for _, service_tag := range service.Tags {
			if service_tag == tag {
				return true, nil
			}
		}
		return false, nil
	}
	return false, nil
}

//服务取消
func (c *ConsulLeader) ServiceDeregister() error {
	return c.Client.Agent().ServiceDeregister(c.name)
}

//创建Session
func (c *ConsulLeader) CreateSession(ch chan struct{}) error {
	entry := &consul.SessionEntry{
		TTL:      "30s",
		Behavior: "release",
	}
	session := c.Client.Session()
	id, _, err := session.Create(entry, &consul.WriteOptions{})
	if err != nil {
		return err
	}
	c.sessionId = id
	go session.RenewPeriodic("5s", c.sessionId, &consul.WriteOptions{}, ch)
	return nil
}

//死循环获取Key
func (c *ConsulLeader) AcquireKey() error {
	node_name, err := c.Client.Agent().NodeName()
	if err != nil {
		return err
	}

	kv := c.Client.KV()
	kvpair := &consul.KVPair{
		Key:     fmt.Sprintf("service/%s.%s/leader", c.tag, c.name),
		Value:   []byte(fmt.Sprintf("service %s.%s running on %s, %d", c.tag, c.name, node_name, c.port)),
		Session: c.sessionId,
	}
	for {
		time.Sleep(time.Second)
		result, _, err := kv.Acquire(kvpair, &consul.WriteOptions{})
		if err != nil {
			continue
		}
		if result == true {
			return nil
		}
	}
}

//释放已经注册的Key
func (c *ConsulLeader) CheckKey() (bool, error) {
	kv := c.Client.KV()
	result, _, err := kv.Get(fmt.Sprintf("service/%s.%s/leader", c.tag, c.name), &consul.QueryOptions{})
	if err != nil {
		return false, err
	}
	if result == nil || result.Session == c.sessionId {
		return true, nil
	}
	return false, nil
}

//释放已经注册的Key
func (c *ConsulLeader) ReleaseKey() (bool, error) {
	kv := c.Client.KV()
	kvpair := &consul.KVPair{
		Key:     fmt.Sprintf("service/%s.%s/leader", c.tag, c.name),
		Session: c.sessionId,
	}

	result, _, err := kv.Release(kvpair, &consul.WriteOptions{})
	if err != nil {
		return false, err
	}
	return result, nil
}

//开始竞选流程
func (c *ConsulLeader) StartElection(tag string, tagStandby string) error {
	retian_chan := make(chan struct{})
	defer close(retian_chan)
	for {
		err := c.ServiceRegister(tagStandby)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.CreateSession(retian_chan)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.AcquireKey()
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.ServiceDeregister()
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.ServiceRegister(tag)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err_count := 0
		for {
			if err_count > 10 {
				seelog.Error("find error to consul")
				return err
			}
			is_leader, err := c.CheckKey()
			if err != nil {
				seelog.Warn(err)
				err_count = err_count + 1
				continue
			}
			if is_leader {
				is_tag_exist, err := c.ServiceConfirm(tagStandby)
				if err != nil {
					seelog.Warn(err)
					err_count = err_count + 1
					continue
				}
				if is_tag_exist {
					err = c.ServiceDeregister()
					if err != nil {
						seelog.Warn(err)
						err_count = err_count + 1
						continue
					}
					err = c.ServiceRegister(tag)
					if err != nil {
						seelog.Warn(err)
						err_count = err_count + 1
						continue
					}
				}
			} else {
				err = c.ServiceDeregister()
				if err != nil {
					seelog.Warn(err)
					err_count = err_count + 1
					continue
				}
				break
			}
			time.Sleep(5 * time.Second)
			err_count = 0
		}
	}
}
