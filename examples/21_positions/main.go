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

	positions, err := client.GetPositionList(cli, accID)
	if err != nil {
		log.Fatalf("GetPositionList failed: %v", err)
	}
	if len(positions) == 0 {
		fmt.Println("(no positions)")
	}
	for _, p := range positions {
		fmt.Printf("POS: %s qty=%.0f cost=%.2f cur=%.2f pnl=%.2f (%.2f%%)\n",
			p.Code, p.Quantity, p.CostPrice, p.CurPrice, p.PnL, p.PnLRate)
	}
}