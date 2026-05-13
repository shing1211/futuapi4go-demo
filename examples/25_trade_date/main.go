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

	dates, err := client.GetTradeDate(context.Background(), mc.Client, constant.Market_HK, "2026-01-01", "2026-04-24")
	if err != nil {
		log.Fatalf("GetTradeDate failed: %v", err)
	}
	for _, d := range dates {
		fmt.Println("TRADE DATE:", d)
	}
}
