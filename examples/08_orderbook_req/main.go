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

	if err := client.Subscribe(context.Background(), cli, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_OrderBook}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	book, err := client.GetOrderBook(context.Background(), cli, constant.Market_US, "NVDA", 10)
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
