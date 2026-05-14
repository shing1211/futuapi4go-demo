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
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Ticker}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	ch := make(chan *push.UpdateTicker, 100)
	chanpkg.SubscribeTicker(context.Background(), mc.Client, constant.Market_US, "NVDA", ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA tick trades (Ctrl+C to exit)...")
	for {
		select {
		case t := <-ch:
			for _, tick := range t.TickerList {
				fmt.Printf("TICKER: price=%.2f vol=%d dir=%d\n",
					tick.GetPrice(), tick.GetVolume(), tick.GetDir())
			}
		case <-sig:
			return
		}
	}
}
