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

	ch := make(chan *push.UpdateKL, 100)
	stop := chanpkg.SubscribeKLine(cli, constant.Market_US, "NVDA", constant.KLType_K_5Min, ch)
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA 5min K-lines (Ctrl+C to exit)...")
	for {
		select {
		case kl := <-ch:
			for _, bar := range kl.KLList {
				fmt.Printf("KL [5m]: %s O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
					*bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice, *bar.Volume)
			}
		case <-sig:
			return
		}
	}
}
