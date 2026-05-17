package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	resp, err := client.QuerySubscription(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("QuerySubscription failed: %v", err)
	}
	fmt.Printf("RemainQuota: %d\n", resp.RemainQuota)
	for _, si := range resp.ConnSubInfoList {
		fmt.Printf("  UsedQuota: %d\n", si.GetUsedQuota())
		for _, sub := range si.GetSubInfoList() {
			fmt.Printf("    SubType: %d\n", sub.GetSubType())
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(resp)
}
