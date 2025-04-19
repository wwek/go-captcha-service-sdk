package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/wenlng/go-captcha-service-sdk/golang/grpcapi"
	"github.com/wenlng/go-captcha-service-sdk/golang/sdlb"
	"github.com/wenlng/go-service-discovery/loadbalancer"
	"github.com/wenlng/go-service-discovery/servicediscovery"
)

// setupGrpcClient .
func setupGrpcClient() (grpcapi.Client, error) {
	sdlbInstance, err := sdlb.NewServiceDiscoveryLB(sdlb.ClientConfig{
		ServiceDiscoveryType: servicediscovery.ServiceDiscoveryTypeConsul,
		Addrs:                "localhost:8500",
		LoadBalancerType:     loadbalancer.LoadBalancerTypeRandom,
		ServiceName:          "go-captcha-service",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to new sdlb: %v", err)
	}

	return grpcapi.NewGRPCClient(grpcapi.ClientConfig{
		APIKey: "my-secret-key-123",
	}, sdlbInstance)
}

func TestGrpcGetData(id string) {
	client, err := setupGrpcClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetData(ctx, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestGrpcCheckData(id, captchaKey, value string) {
	client, err := setupGrpcClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CheckData(ctx, id, captchaKey, value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestGrpcGetStatusInfo(captchaKey string) {
	client, err := setupGrpcClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetStatusInfo(ctx, captchaKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get status info: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestGrpcDelStatusInfo(captchaKey string) {
	client, err := setupGrpcClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.DelStatusInfo(ctx, captchaKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to del status info: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}
