package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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
	if err := client.UnlockTrading(cli, pwdMD5); err != nil {
		log.Fatalf("UnlockTrading failed: %v", err)
	}

	accounts, err := client.GetAccountList(cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	market := acc.TrdMarketAuthList[0]

	if err := client.CancelAllOrder(cli, accID, market, cli.GetTradeEnv()); err != nil {
		log.Fatalf("CancelAllOrder failed: %v", err)
	}
	fmt.Println("All orders cancelled.")
}
