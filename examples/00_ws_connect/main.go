package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_WS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11113"
	}

	secretKey := os.Getenv("FUTU_WS_SECRET")
	if secretKey == "" {
		log.Fatal("FUTU_WS_SECRET environment variable is required")
	}

	fmt.Printf("Connecting via WebSocket to %s ...\n", addr)

	err := cli.ConnectWS(addr, secretKey)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	fmt.Println("Connected!")
	fmt.Printf("ConnID: %d\n", cli.GetConnID())
}