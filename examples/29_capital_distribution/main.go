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

	dist, err := client.GetCapitalDistribution(context.Background(), cli, int32(constant.Market_US), "NVDA")
	if err != nil {
		log.Fatalf("GetCapitalDistribution failed: %v", err)
	}
	fmt.Printf("Main: %.2f  Large: %.2f  Medium: %.2f  Small: %.2f\n",
		dist.MainInflow, dist.BigInflow, dist.MidInflow, dist.SmallInflow)
}
