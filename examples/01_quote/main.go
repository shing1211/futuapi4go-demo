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

	if err := client.Subscribe(context.Background(), cli, int32(constant.Market_US), "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	quote, err := client.GetQuote(context.Background(), cli, int32(constant.Market_US), "NVDA")
	if err != nil {
		log.Fatalf("GetQuote failed: %v", err)
	}
	fmt.Printf("NVDA: price=%.2f open=%.2f high=%.2f low=%.2f vol=%d\n",
		quote.Price, quote.Open, quote.High, quote.Low, quote.Volume)
}
