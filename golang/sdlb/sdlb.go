package sdlb

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/wenlng/go-service-link/foundation/common"
	"github.com/wenlng/go-service-link/servicediscovery"
	"github.com/wenlng/go-service-link/servicediscovery/balancer"
	"github.com/wenlng/go-service-link/servicediscovery/instance"
)

// ServiceDiscoveryType .
const (
	ServiceDiscoveryTypeEtcd      = servicediscovery.ServiceDiscoveryTypeEtcd
	ServiceDiscoveryTypeZookeeper = servicediscovery.ServiceDiscoveryTypeZookeeper
	ServiceDiscoveryTypeConsul    = servicediscovery.ServiceDiscoveryTypeConsul
	ServiceDiscoveryTypeNacos     = servicediscovery.ServiceDiscoveryTypeNacos
	ServiceDiscoveryTypeNone      = servicediscovery.ServiceDiscoveryTypeNone
)

// LoadBalancerType .
const (
	LoadBalancerTypeRandom         = balancer.LoadBalancerTypeRandom
	LoadBalancerTypeRoundRobin     = balancer.LoadBalancerTypeRoundRobin
	LoadBalancerTypeConsistentHash = balancer.LoadBalancerTypeConsistentHash
)

// SDLB ..
type SDLB struct {
	discovery *servicediscovery.DiscoveryWithLB
	config    ClientConfig
	apiKey    string

	isSDLBActive     bool
	isSDLBActiveLock sync.RWMutex

	stop     chan struct{}
	isClosed bool
}

// ClientConfig ..
type ClientConfig struct {
	ServiceDiscoveryType servicediscovery.ServiceDiscoveryType
	LoadBalancerType     balancer.LoadBalancerType
	ServiceName          string
	Addrs                string // 127.0.0.1:8080,192.168.0.1:8080
	Username             string
	Password             string

	TlsConfig *common.TLSConfig

	TTL                   time.Duration
	KeepAlive             time.Duration // Heartbeat interval
	LogOutputHookCallback servicediscovery.OutputLogCallback
}

// NewServiceDiscoveryLB ..
func NewServiceDiscoveryLB(cnf ClientConfig) (*SDLB, error) {
	config := servicediscovery.Config{
		Type:        cnf.ServiceDiscoveryType,
		Addrs:       cnf.Addrs,
		TTL:         cnf.TTL,
		KeepAlive:   cnf.KeepAlive,
		ServiceName: cnf.ServiceName,
		Username:    cnf.Username,
		Password:    cnf.Password,
		TlsConfig:   cnf.TlsConfig,
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

	discovery.SetOutputLogCallback(cnf.LogOutputHookCallback)

	c := &SDLB{
		discovery: discovery,
		config:    cnf,
	}

	c.start()

	return c, nil
}

func (c *SDLB) start() {
	go c.watchInstances()
}

func (c *SDLB) watchInstances() {
	if c.discovery == nil {
		return
	}

	if c.stop != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := c.discovery.Watch(ctx, c.config.ServiceName)
	if err != nil {
		if c.config.LogOutputHookCallback != nil {
			c.config.LogOutputHookCallback(servicediscovery.OutputLogTypeError, fmt.Sprintf("[GoCaptchaServiceSDK-SDLB] Failed to discovery watch: %v", err))
		}

		if !c.isClosed {
			c.start()
		}

		return
	}

	c.stop = make(chan struct{})
	for {
		select {
		case <-c.stop:
			return
		case instances, ok := <-ch:
			if !ok {
				return
			}
			instancesStr, _ := json.Marshal(instances)
			if c.config.LogOutputHookCallback != nil {
				c.config.LogOutputHookCallback(servicediscovery.OutputLogTypeInfo, fmt.Sprintf("[GoCaptchaServiceSDK-SDLB] Discovered instances: %d, list: %v", len(instances), string(instancesStr)))
			}
		}
	}
}

// Select .
func (c *SDLB) Select(key string) (instance.ServiceInstance, error) {
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
func (c *SDLB) Close() error {
	if c.stop != nil {
		c.stop <- struct{}{}
	}
	c.isClosed = true
	return c.discovery.Close()
}
