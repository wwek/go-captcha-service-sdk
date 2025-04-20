package main

import (
	"fmt"

	"github.com/wenlng/go-captcha-service-sdk/golang/grpcapi"
	"github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
	"github.com/wenlng/go-captcha-service-sdk/golang/sdlb"
	"github.com/wenlng/go-service-link/servicediscovery"
	"github.com/wenlng/go-service-link/servicediscovery/balancer"
)

var restapiCli resetapi.Client
var grpcapiCli grpcapi.Client
var sdlbInts *sdlb.SDLB

const serviceName = "go-captcha-service"
const Addrs = "localhost:8500"
const ApiKey = "my-secret-key-123"

// setupHttpClient .
func setupHttpClient() error {
	sdlbInstance, err := sdlb.NewServiceDiscoveryLB(sdlb.ClientConfig{
		ServiceDiscoveryType: servicediscovery.ServiceDiscoveryTypeConsul,
		Addrs:                Addrs,
		LoadBalancerType:     balancer.LoadBalancerTypeRandom,
		ServiceName:          serviceName,
	})

	if err != nil {
		return fmt.Errorf("failed to new sdlb: %v", err)
	}

	sdlbInts = sdlbInstance
	restapiCli, err = resetapi.NewHTTPClient(resetapi.ClientConfig{
		APIKey: ApiKey,
	}, sdlbInstance)

	return err
}

// setupGrpcClient .
func setupGrpcClient() error {
	sdlbInstance, err := sdlb.NewServiceDiscoveryLB(sdlb.ClientConfig{
		ServiceDiscoveryType: servicediscovery.ServiceDiscoveryTypeConsul,
		Addrs:                Addrs,
		LoadBalancerType:     balancer.LoadBalancerTypeRandom,
		ServiceName:          serviceName,
	})

	if err != nil {
		return fmt.Errorf("failed to new sdlb: %v", err)
	}

	grpcapiCli, err = grpcapi.NewGRPCClient(grpcapi.ClientConfig{
		APIKey: ApiKey,
	}, sdlbInstance)
	return err
}

func closeSDLB() {
	if sdlbInts != nil {
		sdlbInts.Close()
	}
}
