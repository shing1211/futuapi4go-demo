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

	if err := client.Subscribe(context.Background(), cli, constant.Market_US, "NVDA", []constant.SubType{
		constant.SubType_Quote,
		constant.SubType_Ticker,
		constant.SubType_K_Day,
	}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}
	fmt.Println("Subscribed to NVDA quote, ticker, day K-line.")

	if err := client.UnsubscribeAll(context.Background(), cli); err != nil {
		log.Fatalf("UnsubscribeAll failed: %v", err)
	}
	fmt.Println("Unsubscribed from all.")
}
