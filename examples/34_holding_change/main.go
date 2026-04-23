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

	changes, err := client.GetHoldingChangeList(cli,
		int32(constant.Market_US), "NVDA",
		1, // holderCategory: 1=Mutual Fund
		"2026-01-01", "2026-04-24",
	)
	if err != nil {
		log.Fatalf("GetHoldingChangeList failed: %v", err)
	}
	for _, c := range changes {
		fmt.Printf("HOLDING CHANGE: holder=%s qty=%.0f ratio=%.2f%%\n",
			c.HolderName, c.HoldingQty, c.HoldingRatio)
	}
}
