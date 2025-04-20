package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func TestHttpGetData(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.GetData(ctx, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpCheckData(id, captchaKey, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.CheckData(ctx, id, captchaKey, value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpGetStatusInfo(captchaKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.GetStatusInfo(ctx, captchaKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get status info: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpDelStatusInfo(captchaKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.DelStatusInfo(ctx, captchaKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to del status info: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpGetResourceList() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.GetResourceList(ctx, "/gocaptcha")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpGetConfig() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.GetConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpUpdateConfig(jsonStr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := restapiCli.UpdateHotConfig(ctx, jsonStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestHttpUploadResource(dirname string, files []*os.File) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, someDone, err := restapiCli.UploadResource(ctx, dirname, files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	fmt.Println(">>>>>>>>", res)
	fmt.Println("Some files failed to be uploaded >>>>>>>>", someDone)
}
