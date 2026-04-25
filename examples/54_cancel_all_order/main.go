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
	cli := client.New().WithTradeEnv(1) // Real trading
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	pwdMD5 := os.Getenv("FUTU_TRADE_PWD")
	if pwdMD5 == "" {
		log.Fatal("FUTU_TRADE_PWD environment variable not set")
	}
	if len(pwdMD5) != 32 {
		log.Fatal("FUTU_TRADE_PWD must be a 32-char MD5 hex string")
	}
	if err := client.UnlockTrading(context.Background(), cli, pwdMD5); err != nil {
		log.Fatalf("UnlockTrading failed: %v", err)
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
	market := constant.TrdMarket(acc.TrdMarketAuthList[0])
	trdEnv := constant.TrdEnv(acc.TrdEnv)

	if err := client.CancelAllOrder(context.Background(), cli, accID, market, trdEnv); err != nil {
		log.Fatalf("CancelAllOrder failed: %v", err)
	}
	fmt.Println("All orders cancelled.")
}
