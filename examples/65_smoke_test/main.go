// 99_smoke_test is a fast happy-path example intended for CI / smoke checks.
// It performs a single connection and three cheap read-only calls, exiting
// nonzero on any failure so it can gate automation.
//
// Checks:
//  1. Connect + GetGlobalState (server alive, quote login available)
//  2. GetAccountList   (trade subsystem reachable — may be empty in sim mode)
//  3. GetQuote NVDA    (market data round-trip)
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func step(name string, fn func() error) bool {
	start := time.Now()
	if err := fn(); err != nil {
		fmt.Printf("[FAIL] %s: %v\n", name, err)
		return false
	}
	fmt.Printf("[PASS] %s (%v)\n", name, time.Since(start).Round(time.Millisecond))
	return true
}

func main() {
	ok := true

	step("connect+global_state", func() error {
		mc := connect.MustConnect(context.Background())
		defer mc.Close()

		state, err := client.GetGlobalState(context.Background(), mc.Client)
		if err != nil {
			return fmt.Errorf("GetGlobalState: %w", err)
		}
		if state.ServerVer <= 0 {
			return fmt.Errorf("GetGlobalState: server version %d looks invalid", state.ServerVer)
		}
		fmt.Printf("       serverVer=%d  qotLogin=%v  trdLogin=%v\n",
			state.ServerVer, state.QotLogined, state.TrdLogined)
		return nil
	})

	ok = step("account_list", func() error {
		mc := connect.MustConnect(context.Background())
		defer mc.Close()

		accounts, err := client.GetAccountList(context.Background(), mc.Client)
		if err != nil {
			return fmt.Errorf("GetAccountList: %w", err)
		}
		fmt.Printf("       accounts=%d\n", len(accounts))
		return nil
	}) && ok

	ok = step("quote_NVDA", func() error {
		mc := connect.MustConnect(context.Background())
		defer mc.Close()

		if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
			return fmt.Errorf("Subscribe: %w", err)
		}
		quote, err := client.GetQuote(context.Background(), mc.Client, constant.Market_US, "NVDA")
		if err != nil {
			return fmt.Errorf("GetQuote: %w", err)
		}
		if quote.Name == "" {
			return fmt.Errorf("GetQuote: empty name for US.NVDA")
		}
		fmt.Printf("       NVDA %s = %.2f (volume %d)\n", quote.Name, quote.Price, quote.Volume)
		return nil
	}) && ok

	if !ok {
		fmt.Println("SMOKE FAIL")
		os.Exit(1)
	}
	fmt.Println("SMOKE OK")
}