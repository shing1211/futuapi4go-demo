// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example: Getting Started - Basic usage of futuapi4go SDK
//
// This example demonstrates the fundamental features of the futuapi4go SDK:
// - Connect to OpenD
// - Query real-time quotes
// - Fetch historical K-lines
// - Subscribe to push notifications
//
// Run: go run examples/getting_started
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go/pkg/qot"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}

	fmt.Println("=== futuapi4go - Getting Started Example ===")
	fmt.Printf("Connecting to %s...\n", addr)

	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	fmt.Printf("Connected! ConnID=%d\n\n", cli.GetConnID())

	runQuoteExample(cli)
	runSnapshotExample(cli)
	runKLineExample(cli)
	runPushExample(cli)

	fmt.Println("\n=== Example Complete ===")
}

func runQuoteExample(cli *client.Client) {
	fmt.Println("--- Real-time Quote ---")

	security := &qotcommon.Security{
		Market: ptrInt32(constant.Market_HK),
		Code:   ptrStr("00700"),
	}

	quotes, err := qot.GetBasicQot(context.Background(), cli.Inner(), []*qotcommon.Security{security})
	if err != nil {
		log.Printf("GetBasicQot failed: %v", err)
		return
	}

	for _, q := range quotes {
		fmt.Printf("  %s (%s): %.2f\n", q.Name, q.Security.Code, q.CurPrice)
	}
}

func runSnapshotExample(cli *client.Client) {
	fmt.Println("\n--- Security Snapshot ---")

	securities := []*qotcommon.Security{
		{Market: ptrInt32(constant.Market_HK), Code: ptrStr("00700")},
		{Market: ptrInt32(constant.Market_HK), Code: ptrStr("09988")},
		{Market: ptrInt32(constant.Market_US), Code: ptrStr("AAPL")},
	}

	resp, err := qot.GetSecuritySnapshot(cli.Inner(), &qot.GetSecuritySnapshotRequest{
		SecurityList: securities,
	})
	if err != nil {
		log.Printf("GetSecuritySnapshot failed: %v", err)
		return
	}

	for _, s := range resp.SnapshotList {
		basic := s.Basic
		name := ""
		if basic.Name != nil {
			name = *basic.Name
		}
		code := ""
		if basic.Security != nil && basic.Security.Code != nil {
			code = *basic.Security.Code
		}
		curPrice := 0.0
		if basic.CurPrice != nil {
			curPrice = *basic.CurPrice
		}
		fmt.Printf("  %s: %.2f\n", code, curPrice)
		_ = name
	}
}

func runKLineExample(cli *client.Client) {
	fmt.Println("\n--- Historical K-lines ---")

	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -10).Format("2006-01-02")

	klines, err := client.RequestHistoryKL(
		cli,
		constant.Market_HK,
		"00700",
		int32(constant.KLType_K_Day),
		startDate,
		endDate,
	)
	if err != nil {
		log.Printf("RequestHistoryKL failed: %v", err)
		return
	}

	fmt.Printf("  Retrieved %d K-lines\n", len(klines))
	if len(klines) > 0 {
		fmt.Printf("  Latest: %s Close=%.2f\n",
			klines[len(klines)-1].Time, klines[len(klines)-1].Close)
	}
}

func runPushExample(cli *client.Client) {
	fmt.Println("\n--- Real-time Push (press Ctrl+C to stop) ---")

	security := &qotcommon.Security{
		Market: ptrInt32(constant.Market_HK),
		Code:   ptrStr("00700"),
	}

	_, err := qot.Subscribe(cli.Inner(), &qot.SubscribeRequest{
		SecurityList:     []*qotcommon.Security{security},
		SubTypeList:      []qot.SubType{qot.SubType(constant.SubType_Quote)},
		IsSubOrUnSub:     true,
		IsRegOrUnRegPush: true,
	})
	if err != nil {
		log.Printf("Subscribe failed: %v", err)
		return
	}

	cli.RegisterHandler(constant.ProtoID_Qot_UpdateBasicQot, func(protoID uint32, body []byte) {
		quote, err := client.ParsePushQuote(body)
		if err != nil || quote == nil {
			return
		}
		fmt.Printf("\r  [%s] %s: %.2f   ",
			time.Now().Format("15:04:05"),
			quote.Code,
			quote.CurPrice)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
	case <-time.After(5 * time.Second):
	}

	fmt.Println("\n  (Unsubscribing...)")
	_, _ = qot.Subscribe(cli.Inner(), &qot.SubscribeRequest{
		SecurityList:     []*qotcommon.Security{security},
		SubTypeList:      []qot.SubType{qot.SubType(constant.SubType_Quote)},
		IsSubOrUnSub:     false,
		IsRegOrUnRegPush: false,
	})
}

func ptrStr(s string) *string { return &s }
func ptrInt32(v int32) *int32 { return &v }
