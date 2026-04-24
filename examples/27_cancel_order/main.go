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

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	market := acc.TrdMarketAuthList[0]

	// Get open orders, cancel the first one
	orders, err := client.GetOrderList(cli, accID)
	if err != nil {
		log.Fatalf("GetOrderList failed: %v", err)
	}
	if len(orders) == 0 {
		fmt.Println("(no open orders to cancel)")
		return
	}

	order := orders[0]
	fmt.Printf("Cancelling order %d (%s)...\n", order.OrderID, order.Code)
	_, err = client.ModifyOrder(cli,
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
