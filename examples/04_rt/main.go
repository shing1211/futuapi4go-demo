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

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_RT}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	ch := make(chan *push.UpdateRT, 100)
	chanpkg.SubscribeRT(context.Background(), mc.Client, constant.Market_US, "NVDA", ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA tick-by-tick data (Ctrl+C to exit)...")
	for {
		select {
		case rt := <-ch:
			for _, r := range rt.RTList {
				fmt.Printf("RT: time=%s price=%.2f vol=%d avg=%.2f\n",
					r.GetTime(), r.GetPrice(), r.GetVolume(), r.GetAvgPrice())
			}
		case <-sig:
			return
		}
	}
}
