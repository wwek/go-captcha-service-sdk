package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGrpc(t *testing.T) {
	err := setupGrpcClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}
	defer closeSDLB()

	TestGrpcGetData("click_dark_ch")

	//TestGrpcCheckData("click_dark_ch", "25011d90-1cc8-11f0-b41e-8c85907c8cf5", "10,25,63,57")

	//TestGrpcGetStatusInfo("25011d90-1cc8-11f0-b41e-8c85907c8cf5")

	//TestGrpcDelStatusInfo("bba172f4-1c73-11f0-95b9-8c85907c8cf5")
}
