package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
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

	ch := make(chan *push.UpdateOrderBook, 100)
	chanpkg.SubscribeOrderBook(cli, constant.Market_US, "NVDA", ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA order book (Ctrl+C to exit)...")
	for {
		select {
		case ob := <-ch:
			for i, bid := range ob.OrderBookBidList {
				if i >= 5 {
					break
				}
				fmt.Printf("BID  [%d]: price=%.2f vol=%d\n", i, bid.GetPrice(), bid.GetVolume())
			}
			for i, ask := range ob.OrderBookAskList {
				if i >= 5 {
					break
				}
				fmt.Printf("ASK  [%d]: price=%.2f vol=%d\n", i, ask.GetPrice(), ask.GetVolume())
			}
			fmt.Println("---")
		case <-sig:
			return
		}
	}
}
