package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
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
	_ = accID // unused in this example

	orders, err := client.GetOrderList(context.Background(), mc.Client, accID)
	if err != nil {
		log.Fatalf("GetOrderList failed: %v", err)
	}
	if len(orders) == 0 {
		fmt.Println("(no open orders)")
	}
	for _, o := range orders {
		fmt.Printf("ORDER: id=%d %s %s qty=%.0f price=%.2f status=%d\n",
			o.OrderID, o.Code, o.Name, o.Qty, o.Price, o.OrderStatus)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(orders)
}
