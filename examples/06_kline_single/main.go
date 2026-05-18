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

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_K_Day}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	klineResult, err := client.GetKLines(context.Background(), mc.Client, constant.Market_US, "NVDA", constant.KLType_K_Day, 10)
	if err != nil {
		log.Fatalf("GetKLines failed: %v", err)
	}
	for _, bar := range klineResult.Items {
		fmt.Printf("%s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
			bar.Time, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
	}

	fmt.Println("\n── KLine Data (JSON) ─────────────────────")
	display.PrintJSON(klineResult.Items)
}
