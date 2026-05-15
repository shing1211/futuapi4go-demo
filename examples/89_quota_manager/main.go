package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Subscription Quota Manager Demo ===")
	fmt.Println("Manage subscription quotas with batch operations")
	fmt.Println()

	fmt.Println("--- Step 1: Check Initial Quota ---")
	quota, err := mc.Client.System().GetUsedQuota(ctx)
	if err != nil {
		fmt.Printf("GetUsedQuota: %v\n", err)
	} else {
		fmt.Printf("  Used Sub Quota:   %d\n", quota.UsedSubQuota)
		fmt.Printf("  Used KLine Quota: %d\n", quota.UsedKLineQuota)
	}
	fmt.Println()

	fmt.Println("--- Step 2: Check Current Subscriptions ---")
	subInfo, err := client.GetSubInfo(ctx, mc.Client)
	if err != nil {
		fmt.Printf("GetSubInfo: %v\n", err)
	} else {
		fmt.Printf("  IsSub: %v\n", subInfo.IsSub)
		fmt.Printf("  SubTypes: %v\n", subInfo.SubTypes)
		fmt.Printf("  Quota: %s\n", subInfo.Security)
	}

	subInfo2, err := client.QuerySubscription(ctx, mc.Client)
	if err != nil {
		fmt.Printf("QuerySubscription: %v\n", err)
	} else {
		fmt.Printf("  Query: %+v\n", subInfo2)
	}

	fmt.Println()
	fmt.Println("--- Step 3: Batch Subscribe Symbols ---")
	symbols := []string{"NVDA", "AAPL", "TSLA", "MSFT", "GOOGL"}
	subTypes := []constant.SubType{constant.SubType_Quote}

	err = client.SubscribeSymbols(ctx, mc.Client, constant.Market_US, symbols, subTypes)
	if err != nil {
		fmt.Printf("SubscribeSymbols: %v\n", err)
	} else {
		fmt.Printf("  Subscribed %d symbols in one call\n", len(symbols))
	}

	fmt.Println()
	fmt.Println("--- Step 4: Verify Subscriptions After Subscribe ---")
	subInfo3, err := client.GetSubInfo(ctx, mc.Client)
	if err != nil {
		fmt.Printf("GetSubInfo: %v\n", err)
	} else {
		fmt.Printf("    Subscribed: %v\n", subInfo3.IsSub)
		fmt.Printf("    SubTypes:   %v\n", subInfo3.SubTypes)
	}

	fmt.Println()
	fmt.Println("--- Step 5: Check Updated Quota ---")
	quota2, err := mc.Client.System().GetUsedQuota(ctx)
	if err != nil {
		fmt.Printf("GetUsedQuota: %v\n", err)
	} else {
		fmt.Printf("  Used Sub Quota:   %d (delta: %+d)\n",
			quota2.UsedSubQuota, quota2.UsedSubQuota-quota.UsedSubQuota)
		fmt.Printf("  Used KLine Quota: %d\n", quota2.UsedKLineQuota)
	}

	fmt.Println()
	fmt.Println("--- Step 6: Batch Unsubscribe Symbols ---")
	allSubTypes := []constant.SubType{constant.SubType_Quote}
	err = client.Unsubscribe(ctx, mc.Client, constant.Market_US, "TSLA", allSubTypes)
	if err != nil {
		fmt.Printf("Unsubscribe TSLA: %v\n", err)
	} else {
		fmt.Printf("  Unsubscribed TSLA (single)\n")
	}

	err = client.Unsubscribe(ctx, mc.Client, constant.Market_US, "GOOGL", allSubTypes)
	if err != nil {
		fmt.Printf("Unsubscribe GOOGL: %v\n", err)
	} else {
		fmt.Printf("  Unsubscribed GOOGL (single)\n")
	}

	fmt.Println()
	fmt.Println("--- Step 7: Unsubscribe All ---")
	err = client.UnsubscribeAll(ctx, mc.Client)
	if err != nil {
		fmt.Printf("UnsubscribeAll: %v\n", err)
	} else {
		fmt.Printf("  All subscriptions cleared\n")
	}

	fmt.Println()
	fmt.Println("--- Step 8: Final Quota Check ---")
	quota3, err := mc.Client.System().GetUsedQuota(ctx)
	if err != nil {
		fmt.Printf("GetUsedQuota: %v\n", err)
	} else {
		fmt.Printf("  Used Sub Quota:   %d\n", quota3.UsedSubQuota)
		fmt.Printf("  Used KLine Quota: %d\n", quota3.UsedKLineQuota)
	}

	fmt.Println()
	fmt.Println("=== Quota Management Complete ===")
}
