package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/futuapi4go/client"
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

	ch := make(chan *push.UpdatePriceReminder, 100)
	stop := chanpkg.SubscribePriceReminder(cli, ch)
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for price reminders (Ctrl+C to exit)...")
	for {
		select {
		case pr := <-ch:
			fmt.Printf("PRICE REMINDER: code=%s name=%s price=%.2f setVal=%.2f curVal=%.2f\n",
				pr.Security.GetCode(), pr.Name, pr.Price, pr.SetValue, pr.CurValue)
		case <-sig:
			return
		}
	}
}
