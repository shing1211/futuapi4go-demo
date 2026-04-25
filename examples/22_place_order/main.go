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

	accounts, err := client.GetAccountList(context.Background(), cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	market := acc.TrdMarketAuthList[0]

	result, err := client.PlaceOrder(context.Background(), cli,
		accID,
		market,
		"00100",  // Tencent (from position list)
		int32(constant.TrdSide_Buy),
		int32(constant.OrderType_Normal),
		100.0, 100, 1, // qty=100 (1 lot for HK)
	)
	if err != nil {
		fmt.Printf("PlaceOrder: %v\n", err)
		return
	}
	fmt.Printf("Order placed: id=%d idEx=%s\n", result.OrderID, result.OrderIDEx)
}