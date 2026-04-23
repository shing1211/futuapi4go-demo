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

	err := client.UnlockTrading(cli, "futu123456")
	if err != nil {
		log.Fatalf("UnlockTrading failed: %v (hint: set password env var or use empty for simulator)", err)
	}
	fmt.Println("Trading unlocked.")
}
