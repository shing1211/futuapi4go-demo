package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
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

	alerts, err := client.GetPriceReminder(context.Background(), cli, constant.Market_US, "NVDA")
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
}
