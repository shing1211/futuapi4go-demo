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

	accounts, err := client.GetAccountList(cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	accID := accounts[0].AccID
	for _, acc := range accounts {
		if acc.TrdEnv == int32(constant.TrdEnv_Real) {
			accID = acc.AccID
			break
		}
	}

	accIDs := []uint64{accID}

	if err := client.SubAccPush(cli, accIDs); err != nil {
		log.Fatalf("SubAccPush failed: %v", err)
	}
	fmt.Printf("Subscribed to %d account push notifications.\n", len(accIDs))
}
