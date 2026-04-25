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

	accID := accounts[0].AccID
	for _, acc := range accounts {
		if acc.TrdEnv == int32(constant.TrdEnv_Real) {
			accID = acc.AccID
			break
		}
	}

	// ReconfirmOrder is used to confirm an order that requires additional verification
	result, err := client.ReconfirmOrder(cli,
		accID,
		int32(constant.TrdMarket_US),
		0, // orderID: 0 means not confirming a specific order
		0, // reason: 0=Normal
	)
	if err != nil {
		log.Fatalf("ReconfirmOrder failed: %v", err)
	}
	fmt.Printf("ReconfirmOrder: accId=%d\n", result.AccID)
}
