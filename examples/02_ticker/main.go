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
				fmt.Printf("TICKER: time=%s seq=%d dir=%d price=%.2f vol=%d turnover=%.2f recvTime=%.0f type=%d typeSign=%d ts=%.0f\n",
					tick.GetTime(), tick.GetSequence(), tick.GetDir(), tick.GetPrice(), tick.GetVolume(),
					tick.GetTurnover(), tick.GetRecvTime(), tick.GetType(), tick.GetTypeSign(), tick.GetTimestamp())
			}
		case <-sig:
			return
		}
	}
}
