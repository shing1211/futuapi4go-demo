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

	result, err := client.PlaceOrder(cli,
		accID,
		int32(constant.TrdMarket_US),
		"NVDA",
		int32(constant.TrdSide_Buy),
		int32(constant.OrderType_Normal),
		100.0, 1,
	)
	if err != nil {
		log.Fatalf("PlaceOrder failed: %v", err)
	}
	fmt.Printf("Order placed: id=%d idEx=%s\n", result.OrderID, result.OrderIDEx)
}