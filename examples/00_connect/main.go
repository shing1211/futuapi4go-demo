// 00_connect is the most basic connection example.
//
// It demonstrates a simple HA connection and displays a full enriched report
// including server info, user info, market rights, quota usage, and delay stats.
//
// Usage:
//
//	go run ./examples/00_connect
package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	fmt.Println("=== Basic Connect Demo ===")
	fmt.Println()

	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	fmt.Println("Connected!")
	fmt.Printf("  Host:   %s:%d\n", mc.Info.Host, mc.Info.Port)
	fmt.Printf("  RSA:    %v\n", mc.Info.RSAUsed)
	fmt.Println()

	display.PrintAll(context.Background(), mc)
}
