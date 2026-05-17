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

	fmt.Println("=== RequestHistoryKL (Working with current OpenD) ===")
	klines, err := client.RequestHistoryKL(
		context.Background(), mc.Client,
		constant.Market_US, "NVDA",
		constant.KLType_K_Day,
		"2026-01-01", "2026-04-24",
	)
	if err != nil {
		log.Fatalf("RequestHistoryKL failed: %v", err)
	}
	for _, bar := range klines {
		fmt.Printf("%s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
			bar.Time, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(klines)
}
