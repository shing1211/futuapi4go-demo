package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ch := make(chan *push.UpdatePriceReminder, 100)
	stop, err := chanpkg.SubscribePriceReminder(context.Background(), mc.Client, ch)
	if err != nil {
		fmt.Printf("SubscribePriceReminder: %v\n", err)
		return
	}
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
