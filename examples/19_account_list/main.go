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
	for i, acc := range accounts {
		fmt.Printf("Account %d: AccID=%d TrdEnv=%d AccType=%d\n", i, acc.AccID, acc.TrdEnv, acc.AccType)
	}

	accID := accounts[0].AccID
	trdMkt := accounts[0].TrdMarketAuthList[0]
	for _, acc := range accounts {
		if acc.TrdEnv == int32(constant.TrdEnv_Real) && len(acc.TrdMarketAuthList) > 0 {
			accID = acc.AccID
			trdMkt = acc.TrdMarketAuthList[0]
			break
		}
	}

	funds, err := client.GetAccountInfo(cli, accID, trdMkt)
	if err != nil {
		log.Fatalf("GetAccountInfo failed: %v", err)
	}
	fmt.Printf("Cash: %.2f  Power: %.2f  Frozen: %.2f\n",
		funds.Cash, funds.Power, funds.FrozenCash)
	fmt.Printf("Assets: %.2f  UnrealizedPL: %.2f  RealizedPL: %.2f\n",
		funds.TotalAssets, funds.UnrealizedPL, funds.RealizedPL)
}