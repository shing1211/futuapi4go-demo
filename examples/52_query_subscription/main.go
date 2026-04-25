package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	if err := client.Subscribe(context.Background(), cli, int32(constant.Market_US), "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	resp, err := client.QuerySubscription(cli)
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
}
