// 00_connect_ha demonstrates HA (High Availability) connection with the following features:
//
//   - Multi-host failover with parallel TCP probe (fastest host wins)
//   - Per-host RSA encryption configuration
//   - Auto-reconnect on connection loss with exponential backoff
//   - Keep-alive monitoring every 30 seconds
//   - Connection state machine (Disconnected → Connecting → Connected → Reconnecting)
//
// It also demonstrates querying system-level APIs after connection:
//   - GetGlobalState: server version, connection ID, login status, market states
//   - GetUserInfo: user ID, nickname, quote rights per market, quota limits
//   - GetUsedQuota: current subscription and historical K-line quota usage
//   - GetDelayStatistics: network latency for quote push, request-reply, and order placement
//
// Configuration via .env file (or environment variables):
//
//	FUTU_OPEND_HOSTS=127.0.0.1:11111:false,172.18.208.88:11111:true
//	FUTU_ADDR=127.0.0.1:11111          # single host fallback
//	FUTU_RSA_KEY=/etc/futu/keys/private_key.pem
//	FUTU_TCP_TIMEOUT=3                  # TCP probe timeout in seconds
//
// Usage:
//
//	go run ./examples/00_connect_ha
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	fmt.Println("=== HA Gateway Selection Demo ===")
	fmt.Println()

	mc := &connect.ManagedConnection{}

	mc.OnStateChange = func(old, new connect.State) {
		fmt.Printf("[State] %v → %v\n", old, new)
	}

	mc.OnError = func(err error) {
		fmt.Printf("[Error] %v\n", err)
	}

	mc.OnConnect = func(info *connect.ConnectionInfo) {
		fmt.Printf("[OnConnect] Successfully connected to %s:%d (RSA=%v)\n",
			info.Host, info.Port, info.RSAUsed)
	}

	ctx := context.Background()
	*mc = *connect.MustConnect(ctx)

	defer mc.Close()

	fmt.Println()
	fmt.Println("Connected!")
	fmt.Printf("  Host:      %s:%d\n", mc.Info.Host, mc.Info.Port)
	fmt.Printf("  RSA Used:  %v\n", mc.Info.RSAUsed)
	fmt.Printf("  State:     %v\n", mc.State)

	fmt.Println()

	display.PrintAll(ctx, mc)

	fmt.Println("--- Connection will auto-reconnect if connection drops ---")
	fmt.Println("Press Ctrl+C to exit")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		<-ticker.C
		fmt.Printf("[Heartbeat] State=%v\n", mc.State)
	}

	fmt.Println("\nDemo complete.")
}
