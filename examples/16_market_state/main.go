package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	state, err := client.GetMarketState(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetMarketState failed: %v", err)
	}
	fmt.Printf("MarketState: %d (1=pre, 2=open, 3=closed, 5=post)\n", state)
}
