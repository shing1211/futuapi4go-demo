// 00_ws_connect demonstrates connecting to FutuOpenD via WebSocket.
//
// WebSocket connections are useful for browser-based or JavaScript clients.
// The Go SDK supports both ws:// and wss:// (TLS) protocols.
//
// Configuration via environment variables:
//
//	FUTU_WS_ADDR=127.0.0.1:11113    # WebSocket OpenD address
//	FUTU_WS_SECRET=xxx              # WebSocket secret key (optional)
//
// Usage:
//
//	go run ./examples/00_ws_connect
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
)

func main() {
	wsAddr := os.Getenv("FUTU_WS_ADDR")
	if wsAddr == "" {
		wsAddr = "127.0.0.1:11113"
	}
	secretKey := os.Getenv("FUTU_WS_SECRET")

	cli := client.New()
	if secretKey != "" {
		err := cli.ConnectWS(wsAddr, secretKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ConnectWS failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		err := cli.ConnectWS(wsAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ConnectWS failed: %v\n", err)
			os.Exit(1)
		}
	}
	defer cli.Close()

	host, portStr, _ := net.SplitHostPort(wsAddr)
	port, _ := strconv.Atoi(portStr)

	mc := &connect.ManagedConnection{
		Client: cli,
		Info: &connect.ConnectionInfo{
			Host: host,
			Port: port,
		},
	}

	fmt.Println("Connected via WebSocket!")
	fmt.Printf("  Host:      %s\n", wsAddr)
	fmt.Println()

	display.PrintAll(context.Background(), mc)
}
