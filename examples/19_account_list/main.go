package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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
		fmt.Printf("Account %d: AccID=%d TrdEnv=%d AccType=%d Markets=%v\n",
			i, acc.AccID, acc.TrdEnv, acc.AccType, acc.TrdMarketAuthList)
	}

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	fmt.Printf("Using AccID=%d (TrdEnv=%d) for market=%d\n",
		acc.AccID, acc.TrdEnv, acc.TrdMarketAuthList[0])

	funds, err := client.GetAccountInfo(cli, acc.AccID, acc.TrdMarketAuthList[0])
	if err != nil {
		log.Fatalf("GetAccountInfo failed: %v", err)
	}
	fmt.Printf("Cash: %.2f  Power: %.2f  Frozen: %.2f\n",
		funds.Cash, funds.Power, funds.FrozenCash)
	fmt.Printf("Assets: %.2f  UnrealizedPL: %.2f  RealizedPL: %.2f\n",
		funds.TotalAssets, funds.UnrealizedPL, funds.RealizedPL)
}