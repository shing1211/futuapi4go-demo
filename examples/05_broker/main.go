package main

import (
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

	if err := client.Subscribe(cli, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Broker}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	ch := make(chan *push.UpdateBroker, 100)
	chanpkg.SubscribeBroker(cli, constant.Market_US, "NVDA", ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA broker data (Ctrl+C to exit)...")
	for {
		select {
		case b := <-ch:
			for _, bid := range b.BidBrokerList {
				fmt.Printf("BID BROKER: name=%s pos=%d\n",
					bid.GetName(), bid.GetPos())
			}
			for _, ask := range b.AskBrokerList {
				fmt.Printf("ASK BROKER: name=%s pos=%d\n",
					ask.GetName(), ask.GetPos())
			}
		case <-sig:
			return
		}
	}
}
