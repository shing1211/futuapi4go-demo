package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	alerts, err := client.GetPriceReminder(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetPriceReminder failed: %v", err)
	}
	if len(alerts) == 0 {
		fmt.Println("(no price reminders set)")
	}
	for _, pr := range alerts {
		fmt.Printf("PRICE REMINDER: %s\n", pr.Name)
		for _, item := range pr.ItemList {
			fmt.Printf("  type=%d value=%.2f enable=%v note=%s\n",
				item.Type, item.Value, item.IsEnable, item.Note)
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(alerts)
}
