package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
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
	market := constant.TrdMarket(acc.TrdMarketAuthList[0])

	result, err := client.PlaceOrder(context.Background(), mc.Client,
		accID,
		market,
		"00100",
		constant.TrdSide_Buy,
		constant.OrderType_Normal,
		100.0, 100, 1,
	)
	if err != nil {
		fmt.Printf("PlaceOrder: %v\n", err)
		return
	}
	fmt.Printf("Order placed: id=%d idEx=%s\n", result.OrderID, result.OrderIDEx)
}