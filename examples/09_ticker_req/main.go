package main

import (
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

	if err := client.Subscribe(cli, int32(constant.Market_US), "NVDA", []constant.SubType{constant.SubType_Ticker}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	tickers, err := client.GetTicker(cli, int32(constant.Market_US), "NVDA", 20)
	if err != nil {
		log.Fatalf("GetTicker failed: %v", err)
	}
	for _, t := range tickers {
		fmt.Printf("TICKER: time=%s price=%.2f vol=%d dir=%s\n",
			t.Time, t.Price, t.Volume, t.Direction)
	}
}
