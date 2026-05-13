// 50_unsubscribe demonstrates Subscribe and Unsubscribe for real-time data.
//
// IMPORTANT: US and HK stocks require a minimum 1-minute subscription before unsubscribing.
// This example shows the full flow with proper error handling for the subscription
// duration requirement.
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

	if err := client.Subscribe(context.Background(), mc.Client, market, code, []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}
	fmt.Printf("Subscribed to %s.%s quote.\n", market, code)

	fmt.Println("Waiting 3 seconds to receive quote updates...")
	time.Sleep(3 * time.Second)

	if err := client.Unsubscribe(context.Background(), mc.Client, market, code, []constant.SubType{constant.SubType_Quote}); err != nil {
		fe, ok := constant.AsFutuError(err)
		if ok && fe.Code == -1 {
			fmt.Printf("⚠️  Unsubscribe failed: %s\n", fe.Message)
			fmt.Println("   This happens when subscription is less than 1 minute.")
			fmt.Println("   The subscription will auto-expire after 1 minute.")
			fmt.Println("   For immediate unsubscribe, wait 60+ seconds after subscribing.")
		} else {
			log.Fatalf("Unsubscribe failed: %v", err)
		}
	} else {
		fmt.Printf("Unsubscribed from %s.%s quote.\n", market, code)
	}

	fmt.Println("\n✓ Example complete. Note: HK stocks also require 1-minute minimum subscription.")
}