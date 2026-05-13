package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	secretKey := os.Getenv("FUTU_WS_SECRET")
	if secretKey == "" {
		fmt.Println("FUTU_WS_SECRET environment variable not set, using default")
		secretKey = "test-secret"
	}

	mc := connect.MustConnectWS(context.Background(), secretKey)
	defer mc.Close()

	fmt.Println("Connected!")
	fmt.Printf("ConnID: %d\n", mc.Client.GetConnID())
}