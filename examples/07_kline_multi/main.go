package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	ch := make(chan *push.UpdateKL, 100)
	stop, err := chanpkg.SubscribeKLines(ctx, mc.Client, constant.Market_US, "NVDA",
		[]constant.KLType{constant.KLType_K_1Min, constant.KLType_K_5Min, constant.KLType_K_Day}, ch)
	if err != nil {
		fmt.Printf("SubscribeKLines: %v\n", err)
		return
	}
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA K-lines 1m/5m/day (Ctrl+C to exit)...")
	for {
		select {
		case kl := <-ch:
			prefix := ""
			switch kl.KlType {
			case int32(constant.KLType_K_1Min):
				prefix = "[1min]"
			case int32(constant.KLType_K_5Min):
				prefix = "[5min]"
			case int32(constant.KLType_K_Day):
				prefix = "[day]"
			}
			for _, bar := range kl.KLList {
				if bar == nil {
					continue
				}
				fmt.Printf("%s  %s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
					prefix, bar.Time, bar.OpenPrice, bar.HighPrice, bar.LowPrice, bar.ClosePrice, bar.Volume)
			}
		case <-sig:
			fmt.Println("Shutting down...")
			return
		}
	}
}
