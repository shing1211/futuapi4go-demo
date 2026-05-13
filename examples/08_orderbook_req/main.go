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

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_OrderBook}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	book, err := client.GetOrderBook(context.Background(), mc.Client, constant.Market_US, "NVDA", 10)
	if err != nil {
		log.Fatalf("GetOrderBook failed: %v", err)
	}
	for i, bid := range book.Bids {
		if i >= 5 {
			break
		}
		fmt.Printf("BID  [%d]: price=%.2f vol=%d\n", i, bid.Price, bid.Volume)
	}
	for i, ask := range book.Asks {
		if i >= 5 {
			break
		}
		fmt.Printf("ASK  [%d]: price=%.2f vol=%d\n", i, ask.Price, ask.Volume)
	}
}
