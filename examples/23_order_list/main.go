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

	orders, err := client.GetOrderList(cli, accounts[0].AccID)
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
}
