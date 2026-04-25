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

	dates, err := client.GetTradeDate(context.Background(), cli, constant.Market_HK, "2026-01-01", "2026-04-24")
	if err != nil {
		log.Fatalf("GetTradeDate failed: %v", err)
	}
	for _, d := range dates {
		fmt.Println("TRADE DATE:", d)
	}
}
