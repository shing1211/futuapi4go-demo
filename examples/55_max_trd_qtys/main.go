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
	_ = accID
	orderType := constant.OrderType_Normal

info, err := client.GetMaxTrdQtys(context.Background(), cli,
	accID,
	constant.TrdMarket_HK,
	"00100",
	orderType,
	100.0,
	constant.TrdSecMarket(1),
)
	if err != nil {
		log.Fatalf("GetMaxTrdQtys failed: %v", err)
	}
	fmt.Printf("Max Cash Buy: %.0f  Max Margin Buy: %.0f  Max Pos Sell: %.0f\n",
		info.MaxCashBuy, info.MaxCashAndMarginBuy, info.MaxPositionSell)
}
