package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Ticker}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	result, err := client.GetTicker(context.Background(), mc.Client, constant.Market_US, "NVDA", 20)
	if err != nil {
		log.Fatalf("GetTicker failed: %v", err)
	}
	for _, t := range result.Items {
		fmt.Printf("TICKER: time=%s price=%.2f vol=%d dir=%s\n",
			t.Time, t.Price, t.Volume, t.Direction)
	}

	fmt.Println("\n── Ticker Data (JSON) ────────────────────")
	display.PrintJSON(result)
}
