package sdlb

import (
	"sync"
	"time"

	"github.com/wenlng/go-service-discovery/base"
	"github.com/wenlng/go-service-discovery/loadbalancer"
	"github.com/wenlng/go-service-discovery/servicediscovery"
)

// ServiceDiscoveryType .
const (
	ServiceDiscoveryTypeEtcd      = servicediscovery.ServiceDiscoveryTypeEtcd
	ServiceDiscoveryTypeZookeeper = servicediscovery.ServiceDiscoveryTypeZookeeper
	ServiceDiscoveryTypeConsul    = servicediscovery.ServiceDiscoveryTypeConsul
	ServiceDiscoveryTypeNacos     = servicediscovery.ServiceDiscoveryTypeNone
	ServiceDiscoveryTypeNone      = servicediscovery.ServiceDiscoveryTypeNacos
)

// LoadBalancerType .
const (
	LoadBalancerTypeRandom         = loadbalancer.LoadBalancerTypeRandom
	LoadBalancerTypeRoundRobin     = loadbalancer.LoadBalancerTypeRoundRobin
	LoadBalancerTypeConsistentHash = loadbalancer.LoadBalancerTypeConsistentHash
)

// SDLB ..
type SDLB struct {
	discovery *servicediscovery.DiscoveryWithLB
	config    ClientConfig
	apiKey    string

	isSDLBActive     bool
	isSDLBActiveLock sync.RWMutex

	done      chan bool
	reconnect chan bool
}

// ClientConfig ..
type ClientConfig struct {
	ServiceDiscoveryType  servicediscovery.ServiceDiscoveryType
	LoadBalancerType      loadbalancer.LoadBalancerType
	ServiceName           string
	Addrs                 string // 127.0.0.1:8080,192.168.0.1:8080
	TTL                   time.Duration
	KeepAlive             time.Duration // Heartbeat interval
	LogOutputHookCallback servicediscovery.LogOutputHookFunc
}

// NewServiceDiscoveryLB ..
func NewServiceDiscoveryLB(cnf ClientConfig) (*SDLB, error) {
	config := servicediscovery.Config{
		Type:        cnf.ServiceDiscoveryType,
		Addrs:       cnf.Addrs,
		TTL:         cnf.TTL,
		KeepAlive:   cnf.KeepAlive,
		ServiceName: cnf.ServiceName,
	}

	if config.TTL < 0 {
		config.TTL = time.Second * 10
	}

	if config.KeepAlive < 0 {
		config.KeepAlive = time.Second * 3
	}

	discovery, err := servicediscovery.NewDiscoveryWithLB(config, cnf.LoadBalancerType)
	if err != nil {
		return nil, err
	}

	discovery.SetLogOutputHookFunc(cnf.LogOutputHookCallback)

	return &SDLB{
		discovery: discovery,
		config:    cnf,
	}, nil
}

// Select .
func (c *SDLB) Select(key string) (base.ServiceInstance, error) {
	return c.discovery.Select(c.config.ServiceName, key)
}

// IsActive .
func (c *SDLB) IsActive() bool {
	c.isSDLBActiveLock.RLock()
	defer c.isSDLBActiveLock.RUnlock()
	return c.isSDLBActive
}

// setActive .
func (c *SDLB) setActive(active bool) {
	c.isSDLBActiveLock.Lock()
	defer c.isSDLBActiveLock.Unlock()
	c.isSDLBActive = active
}

// Close ..
func (c *SDLB) close() error {
	return c.discovery.Close()
}

//func (c *SDLB) start() error {
//	defer c.close()
//	c.reconnect = make(chan bool)
//
//	go c.watch()
//	for {
//		select {
//		case <-p.reconnect:
//			if !p.reconn() {
//				go p.delayReconn()
//			} else {
//				go p.receive(false)
//			}
//		case <-p.done:
//			return nil
//		}
//	}
//}
//
//func (c *SDLB) watch(hasConn bool) {
//
//}

// Watch ..
//func (c *SDLB) watch(ctx context.Context, name string) (chan []base.ServiceInstance, error) {
//	return c.discovery.Watch(ctx, c.config.ServiceName)
//}
