// 51_unsubscribe_all demonstrates Subscribe (multiple types) and UnsubscribeAll
// to clear all active subscriptions at once.
//
// Uses HK.00700 (Tencent) to avoid US stock subscription duration restrictions.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	var market constant.Market = constant.Market_HK
	code := "00700" // Tencent

	subs := []constant.SubType{
		constant.SubType_Quote,
		constant.SubType_Ticker,
		constant.SubType_K_Day,
	}

	if err := client.Subscribe(context.Background(), mc.Client, market, code, subs); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}
	fmt.Printf("Subscribed to %s.%s: Quote, Ticker, Day K-line.\n", market, code)

	// Brief pause to confirm subscriptions are active
	fmt.Println("Waiting 3 seconds to confirm subscriptions are active...")
	time.Sleep(3 * time.Second)

	// Query current subscription state before unsubscribing
	subInfo, err := client.QuerySubscription(context.Background(), mc.Client)
	if err != nil {
		fmt.Printf("⚠️  Could not query subscription state: %v\n", err)
	} else if subInfo != nil && len(subInfo.ConnSubInfoList) == 0 {
		fmt.Println("⚠️  No active subscriptions found — may have expired or server returned empty list.")
		fmt.Println("   Attempting UnsubscribeAll anyway...")
	} else if subInfo != nil {
		fmt.Printf("Active subscriptions: %d\n", len(subInfo.ConnSubInfoList))
	}

	if err := client.UnsubscribeAll(context.Background(), mc.Client); err != nil {
		log.Fatalf("UnsubscribeAll failed: %v", err)
	}
	fmt.Println("Unsubscribed from all active subscriptions.")

	// Verify all are cleared
	subInfo, err = client.QuerySubscription(context.Background(), mc.Client)
	if err != nil {
		fmt.Printf("⚠️  Could not verify cleanup: %v\n", err)
	} else if subInfo != nil {
		fmt.Printf("Remaining active subscriptions: %d\n", len(subInfo.ConnSubInfoList))
	}
}