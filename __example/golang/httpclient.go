package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
	"github.com/wenlng/go-captcha-service-sdk/golang/sdlb"
	"github.com/wenlng/go-service-discovery/loadbalancer"
	"github.com/wenlng/go-service-discovery/servicediscovery"
)

// setupHttpClient .
func setupHttpClient() (resetapi.Client, error) {
	sdlbInstance, err := sdlb.NewServiceDiscoveryLB(sdlb.ClientConfig{
		ServiceDiscoveryType: servicediscovery.ServiceDiscoveryTypeConsul,
		Addrs:                "localhost:8500",
		LoadBalancerType:     loadbalancer.LoadBalancerTypeRandom,
		ServiceName:          "go-captcha-service",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to new sdlb: %v", err)
	}

	return resetapi.NewHTTPClient(resetapi.ClientConfig{
		APIKey: "my-secret-key-123",
	}, sdlbInstance)
}

func TestHttpGetData(id string) {
	client, err := setupHttpClient()
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

func TestHttpCheckData(id, captchaKey, value string) {
	client, err := setupHttpClient()
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

func TestHttpGetStatusInfo(captchaKey string) {
	client, err := setupHttpClient()
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

func TestHttpDelStatusInfo(captchaKey string) {
	client, err := setupHttpClient()
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

func TestHttpGetResourceList() {
	client, err := setupHttpClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetResourceList(ctx, "/gocaptcha")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpGetConfig() {
	client, err := setupHttpClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpUpdateConfig(jsonStr string) {
	client, err := setupHttpClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.UpdateHotConfig(ctx, jsonStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpUploadResource(dirname string, files []*os.File) {
	client, err := setupHttpClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, someDone, err := client.UploadResource(ctx, dirname, files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	fmt.Println(">>>>>>>>", res)
	fmt.Println("Some files failed to be uploaded >>>>>>>>", someDone)
}
