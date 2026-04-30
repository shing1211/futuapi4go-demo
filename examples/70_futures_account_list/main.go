package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/trd"
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

	ctx := context.Background()

	fmt.Println("=== Stock Accounts (GetAccountList) ===")
	accounts, err := client.GetAccountList(ctx, cli)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	for i, acc := range accounts {
		env := "Real"
		if acc.TrdEnv == 0 {
			env = "Simulated"
		}
		fmt.Printf("Account %d: AccID=%d TrdEnv=%s AccType=%d Markets=%v\n",
			i, acc.AccID, env, acc.AccType, acc.TrdMarketAuthList)
	}

	fmt.Println("\n=== Futures Accounts (GetAccList with TrdCategory_Future) ===")
	resp, err := cli.Trade().GetAccList(ctx, constant.TrdCategory_Future, true)
	if err != nil {
		log.Fatalf("GetAccList(TrdCategory_Future) failed: %v", err)
	}
	if len(resp.AccList) == 0 {
		fmt.Println("(no futures accounts)")
		return
	}
	for i, acc := range resp.AccList {
		env := "Real"
		if acc.TrdEnv == 0 {
			env = "Simulated"
		}
		authStr := "Unknown"
		if len(acc.TrdMarketAuthList) == 1 {
			switch acc.TrdMarketAuthList[0] {
			case 5:
				authStr = "Futures"
			case 10:
				authStr = "FuturesSimulateHK"
			case 11:
				authStr = "FuturesSimulateUS"
			case 12:
				authStr = "FuturesSimulateSG"
			case 13:
				authStr = "FuturesSimulateJP"
			}
		}
		fmt.Printf("Futures Account %d: AccID=%d TrdEnv=%s Auth=%v (%s)\n",
			i, acc.AccID, env, acc.TrdMarketAuthList, authStr)
	}

	fmt.Println("\n=== Auto-Select Simulated Futures Account ===")
	var simAcc *trd.Acc
	for _, acc := range resp.AccList {
		if acc.TrdEnv == 0 {
			simAcc = acc
			fmt.Printf("Selected simulated futures AccID=%d (Auth=%v)\n",
				acc.AccID, acc.TrdMarketAuthList)
			break
		}
	}
	if simAcc == nil {
		fmt.Println("No simulated futures account found")
	}
}