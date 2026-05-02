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

	fmt.Println("=== GetHistoryKL (v0.5.5) ===")
	klines, err := client.GetHistoryKL(
		context.Background(), cli,
		constant.Market_US, "NVDA",
		constant.KLType_K_Day, constant.RehabType_None,
		"2026-01-01", "2026-04-24", 100,
	)
	if err != nil {
		log.Fatalf("GetHistoryKL failed: %v", err)
	}
	for _, bar := range klines {
		fmt.Printf("%s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
			bar.Time, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
	}
}
