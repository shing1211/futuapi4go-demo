package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	quote, err := client.GetQuote(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetQuote failed: %v", err)
	}
	fmt.Printf("NVDA: price=%.2f open=%.2f high=%.2f low=%.2f vol=%d isSuspended=%v secStatus=%d\n",
		quote.Price, quote.Open, quote.High, quote.Low, quote.Volume, quote.IsSuspended, quote.SecStatus)
}
