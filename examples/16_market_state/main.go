package main

import (
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

	state, err := client.GetMarketState(cli, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetMarketState failed: %v", err)
	}
	fmt.Printf("MarketState: %d (1=pre, 2=open, 3=closed, 5=post)\n", state)
}
