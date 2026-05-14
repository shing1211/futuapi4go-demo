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

	stop := chanpkg.SubscribeKLines(context.Background(), mc.Client, constant.Market_US, "NVDA", map[constant.KLType]func(*push.UpdateKL){
		constant.KLType_K_1Min: func(kl *push.UpdateKL) {
			for _, bar := range kl.KLList {
				fmt.Printf("[1min]  %s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
					*bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice, *bar.Volume)
			}
		},
		constant.KLType_K_5Min: func(kl *push.UpdateKL) {
			for _, bar := range kl.KLList {
				fmt.Printf("[5min]  %s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
					*bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice, *bar.Volume)
			}
		},
		constant.KLType_K_Day: func(kl *push.UpdateKL) {
			for _, bar := range kl.KLList {
				fmt.Printf("[day]   %s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
					*bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice, *bar.Volume)
			}
		},
	})
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA K-lines 1m/5m/day (Ctrl+C to exit)...")
	<-sig
}
