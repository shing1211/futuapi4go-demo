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
	orderType := constant.OrderType_Normal

	info, err := client.GetAccTradingInfo(context.Background(), mc.Client,
	accID,
	constant.TrdMarket_HK,
	"00100",
	orderType,
	100.0,
)
	if err != nil {
		log.Fatalf("GetAccTradingInfo failed: %v", err)
	}
	fmt.Printf("TRADING INFO: cashBuy=%.0f marginBuy=%.0f posSell=%.0f\n",
		info.MaxCashBuy, info.MaxCashAndMarginBuy, info.MaxPositionSell)

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(info)
}
