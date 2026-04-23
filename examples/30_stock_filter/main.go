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

	results, err := client.StockFilter(cli, int32(constant.Market_US), 0, 10)
	if err != nil {
		log.Fatalf("StockFilter failed: %v", err)
	}
	for _, r := range results {
		fmt.Printf("STOCK: code=%s name=%s price=%.2f vol=%d\n",
			r.Security.GetCode(), r.Name, r.CurPrice, r.Volume)
	}
}
