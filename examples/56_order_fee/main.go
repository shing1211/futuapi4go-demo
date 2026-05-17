package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
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
	_ = accID

	fees, err := client.GetOrderFee(context.Background(), mc.Client,
	accID, constant.TrdMarket_HK, []string{"00700"},
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

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(fees)
}
