package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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

	quota, err := client.RequestHistoryKLQuota(context.Background(), cli)
	if err != nil {
		log.Fatalf("RequestHistoryKLQuota failed: %v", err)
	}
	fmt.Printf("History KL Quota: used=%d remain=%d\n", quota.UsedQuota, quota.RemainQuota)
	for _, d := range quota.DetailList {
		fmt.Printf("  %s (%s): timestamp=%d\n", d.Name, d.Security.GetCode(), d.RequestTimestamp)
	}
}
