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

	accID := accounts[0].AccID
	for _, acc := range accounts {
		if acc.TrdEnv == int32(constant.TrdEnv_Real) {
			accID = acc.AccID
			break
		}
	}

	// ReconfirmOrder is used to confirm an order that requires additional verification
	result, err := client.ReconfirmOrder(context.Background(), mc.Client,
	accID,
	constant.TrdMarket_US,
	uint64(0),
	int32(0),
)
	if err != nil {
		log.Fatalf("ReconfirmOrder failed: %v", err)
	}
	fmt.Printf("ReconfirmOrder: accId=%d\n", result.AccID)

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(result)
}
