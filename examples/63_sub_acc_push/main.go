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

	accounts, err := client.GetAccountList(cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	accIDs := make([]uint64, len(accounts))
	for i, acc := range accounts {
		accIDs[i] = acc.AccID
	}

	if err := client.SubAccPush(cli, accIDs); err != nil {
		log.Fatalf("SubAccPush failed: %v", err)
	}
	fmt.Printf("Subscribed to %d account push notifications.\n", len(accIDs))
}
