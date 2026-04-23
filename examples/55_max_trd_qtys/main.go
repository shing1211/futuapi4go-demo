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

	info, err := client.GetMaxTrdQtys(cli,
		accounts[0].AccID,
		int32(constant.Market_US), "NVDA",
		int32(constant.OrderType_Normal), 100.0,
	)
	if err != nil {
		log.Fatalf("GetMaxTrdQtys failed: %v", err)
	}
	fmt.Printf("Max Cash Buy: %.0f  Max Margin Buy: %.0f  Max Pos Sell: %.0f\n",
		info.MaxCashBuy, info.MaxCashAndMarginBuy, info.MaxPositionSell)
}
