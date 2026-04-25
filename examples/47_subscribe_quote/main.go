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

	if err := client.Subscribe(context.Background(), cli, int32(constant.Market_US), "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	ch := make(chan *push.UpdateBasicQot, 100)
	chanpkg.SubscribeQuote(cli, int32(constant.Market_US), "NVDA", ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA quotes (Ctrl+C to exit)...")
	for {
		select {
		case q := <-ch:
			fmt.Printf("QUOTE: price=%.2f open=%.2f high=%.2f low=%.2f vol=%d\n",
				q.CurPrice, q.OpenPrice, q.HighPrice, q.LowPrice, q.Volume)
		case <-sig:
			return
		}
	}
}
