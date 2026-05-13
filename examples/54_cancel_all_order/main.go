// 54_cancel_all_order demonstrates unlocking trading and cancelling all open orders.
//
// WARNING: This cancels ALL open orders on the selected account/market.
// Trading is unlocked via the MD5 hash of your Futu trading password.
// This example uses the SIMULATE trading environment by default.
//
// Prerequisites:
//   - Set FUTU_TRADE_PWD environment variable (32-char MD5 hex of your trading password)
//   - For simulate trading: account must have simulated trading enabled
//   - For real trading: account must have real trading enabled and correct password
//
// Generate MD5 of your password:
//   echo -n "your_password" | md5sum | cut -d' ' -f1
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	pwdMD5 := os.Getenv("FUTU_TRADE_PWD")
	if pwdMD5 == "" {
		fmt.Println()
		fmt.Println("╔════════════════════════════════════════════════════════════════╗")
		fmt.Println("║  FUTU_TRADE_PWD environment variable is not set.              ║")
		fmt.Println("║                                                                ║")
		fmt.Println("║  This example requires a trading password to unlock trading.   ║")
		fmt.Println("║                                                                ║")
		fmt.Println("║  For Simulate Trading:                                         ║")
		fmt.Println("║    1. Enable simulated trading in Futu OpenD settings          ║")
		fmt.Println("║    2. Set FUTU_TRADE_PWD to MD5 of any password (e.g. '123')   ║")
		fmt.Println("║       $ echo -n \"123\" | md5sum | cut -d' ' -f1                   ║")
		fmt.Println("║       18192ef7582ee4e8e561df90c5e94c6d                         ║")
		fmt.Println("║    3. $ export FUTU_TRADE_PWD=18192ef7582ee4e8e561df90c5e94c6d ║")
		fmt.Println("║                                                                ║")
		fmt.Println("║  For Real Trading:                                             ║")
		fmt.Println("║    1. Use your actual Futu trading password                    ║")
		fmt.Println("║    2. Generate MD5 hash as shown above                         ║")
		fmt.Println("║    3. Export the hash as FUTU_TRADE_PWD                       ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("Alternatively, use example 23_order_list to view orders without")
		fmt.Println("cancelling them, or example 54_cancel_all_order for a dry-run.")
		log.Fatal("FUTU_TRADE_PWD not set")
	}

	if len(pwdMD5) != 32 {
		log.Fatalf("FUTU_TRADE_PWD must be 32 hex characters, got %d", len(pwdMD5))
	}

	if err := client.UnlockTrading(context.Background(), mc.Client, pwdMD5); err != nil {
		log.Fatalf("UnlockTrading failed: %v", err)
	}
	fmt.Println("Trading unlocked successfully.")

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v (is account set up for trading?)", err)
	}

	acc := mc.Client.FindAccount(accounts)
	if acc == nil {
		log.Fatal("No account with trading permissions found.")
	}

	accID := acc.AccID
	market := constant.TrdMarket(acc.TrdMarketAuthList[0])
	trdEnv := constant.TrdEnv(acc.TrdEnv)

	trdEnvStr := "Simulate"
	if trdEnv == constant.TrdEnv_Real {
		trdEnvStr = "Real"
	}

	fmt.Printf("Account: %d | Market: %s | Env: %s\n", accID, market.Prefix(), trdEnvStr)
	fmt.Println("Cancelling all open orders...")

	if err := client.CancelAllOrder(context.Background(), mc.Client, accID, market, trdEnv); err != nil {
		log.Fatalf("CancelAllOrder failed: %v", err)
	}
	fmt.Println("✓ All open orders cancelled.")
}