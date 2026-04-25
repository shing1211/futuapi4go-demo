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
	market := acc.TrdMarketAuthList[0]

	info, err := client.GetAccTradingInfo(cli,
		accID,
		market,
		"00100", // Tencent (HK stock)
		int32(constant.OrderType_Normal),
		100.0,
	)
	if err != nil {
		log.Fatalf("GetAccTradingInfo failed: %v", err)
	}
	fmt.Printf("TRADING INFO: cashBuy=%.0f marginBuy=%.0f posSell=%.0f\n",
		info.MaxCashBuy, info.MaxCashAndMarginBuy, info.MaxPositionSell)
}
