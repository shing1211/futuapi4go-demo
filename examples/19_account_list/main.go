package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	for i, acc := range accounts {
		fmt.Printf("Account %d: AccID=%d TrdEnv=%d AccType=%d Markets=%v\n",
			i, acc.AccID, acc.TrdEnv, acc.AccType, acc.TrdMarketAuthList)
	}

	acc := mc.Client.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	fmt.Printf("Using AccID=%d (TrdEnv=%d) for market=%d\n",
		acc.AccID, acc.TrdEnv, acc.TrdMarketAuthList[0])

	funds, err := client.GetAccountInfo(context.Background(), mc.Client, acc.AccID, constant.TrdMarket(acc.TrdMarketAuthList[0]))
	if err != nil {
		log.Fatalf("GetAccountInfo failed: %v", err)
	}
	fmt.Printf("Cash: %.2f  Power: %.2f  Frozen: %.2f\n",
		funds.Cash, funds.Power, funds.FrozenCash)
	fmt.Printf("Assets: %.2f  UnrealizedPL: %.2f  RealizedPL: %.2f\n",
		funds.TotalAssets, funds.UnrealizedPL, funds.RealizedPL)
}