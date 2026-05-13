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

	orders, err := client.GetOrderList(context.Background(), mc.Client, accID)
	if err != nil {
		log.Fatalf("GetOrderList failed: %v", err)
	}
	if len(orders) == 0 {
		fmt.Println("(no open orders to cancel)")
		return
	}

	order := orders[0]
	fmt.Printf("Cancelling order %d (%s)...\n", order.OrderID, order.Code)
	_, err = client.ModifyOrder(context.Background(), mc.Client,
		accID,
		market,
		order.OrderID,
		2, // ModifyOrderOp_Cancel
		0, 0,
	)
	if err != nil {
		log.Fatalf("CancelOrder failed: %v", err)
	}
	fmt.Println("Order cancelled.")
}
