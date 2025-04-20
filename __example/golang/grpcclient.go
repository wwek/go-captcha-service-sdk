package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func TestGrpcGetData(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := grpcapiCli.GetData(ctx, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestGrpcCheckData(id, captchaKey, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := grpcapiCli.CheckData(ctx, id, captchaKey, value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check data: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestGrpcGetStatusInfo(captchaKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := grpcapiCli.GetStatusInfo(ctx, captchaKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get status info: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}

func TestGrpcDelStatusInfo(captchaKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := grpcapiCli.DelStatusInfo(ctx, captchaKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to del status info: %v\n", err)
		return
	}

	resStr, _ := json.Marshal(res)
	fmt.Println(">>>>>>>>", string(resStr))
}
