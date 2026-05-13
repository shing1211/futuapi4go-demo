package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	quota, err := client.RequestHistoryKLQuota(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("RequestHistoryKLQuota failed: %v", err)
	}
	fmt.Printf("History KL Quota: used=%d remain=%d\n", quota.UsedQuota, quota.RemainQuota)
	for _, d := range quota.DetailList {
		fmt.Printf("  %s (%s): timestamp=%d\n", d.Name, d.Security.GetCode(), d.RequestTimestamp)
	}
}
