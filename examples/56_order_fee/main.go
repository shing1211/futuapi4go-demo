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
	market := constant.TrdMarket(acc.TrdMarketAuthList[0])

	fees, err := client.GetOrderFee(context.Background(), cli,
		accID,
		market,
		[]string{},
	)
	if err != nil {
		log.Fatalf("GetOrderFee failed: %v", err)
	}
	for _, f := range fees {
		fmt.Printf("Order Fee: orderIdEx=%s\n", f.OrderIDEx)
		for _, item := range f.FeeList {
			fmt.Printf("  %s: %.4f\n", item.Title, item.Value)
		}
	}
}
