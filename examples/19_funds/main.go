package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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

	funds, err := client.GetFunds(cli, 0)
	if err != nil {
		log.Fatalf("GetFunds failed: %v", err)
	}
	fmt.Printf("Cash: %.2f  Power: %.2f  Frozen: %.2f\n",
		funds.Cash, funds.Power, funds.FrozenCash)
	fmt.Printf("Assets: %.2f  UnrealizedPL: %.2f  RealizedPL: %.2f\n",
		funds.TotalAssets, funds.UnrealizedPL, funds.RealizedPL)
}
