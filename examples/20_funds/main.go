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
	fmt.Printf("Power=%.2f Cash=%.2f Assets=%.2f\n",
		funds.Power, funds.Cash, funds.TotalAssets)
}