package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := mc.Client.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	_ = accID // unused below

	accIDs := []uint64{accID}

	if err := client.SubAccPush(context.Background(), mc.Client, accIDs); err != nil {
		log.Fatalf("SubAccPush failed: %v", err)
	}
	fmt.Printf("Subscribed to %d account push notifications.\n", len(accIDs))
}
