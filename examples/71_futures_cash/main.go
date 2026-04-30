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

	ctx := context.Background()

	fmt.Println("=== Stock Account Cash (GetFunds) ===")
	accounts, err := client.GetAccountList(ctx, cli)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	for _, acc := range accounts {
		if acc.TrdEnv != 0 {
			continue // skip real accounts
		}

		funds, err := client.GetFunds(ctx, cli, acc.AccID)
		if err != nil {
			fmt.Printf("AccID %d: GetFunds failed: %v\n", acc.AccID, err)
			continue
		}
		fmt.Printf("Stock AccID=%d:\n", acc.AccID)
		fmt.Printf("  Cash: %.2f  Power: %.2f  FrozenCash: %.2f\n",
			funds.Cash, funds.Power, funds.FrozenCash)
		fmt.Printf("  TotalAssets: %.2f  UnrealizedPL: %.2f  RealizedPL: %.2f\n",
			funds.TotalAssets, funds.UnrealizedPL, funds.RealizedPL)
	}

	fmt.Println("\n=== Futures Account Cash (GetAccTradingInfo) ===")
	resp, err := cli.Trade().GetAccList(ctx, constant.TrdCategory_Future, true)
	if err != nil {
		log.Fatalf("GetAccList(TrdCategory_Future) failed: %v", err)
	}

	for _, acc := range resp.AccList {
		if acc.TrdEnv != 0 {
			continue // skip real accounts
		}

		fmt.Printf("Futures AccID=%d:\n", acc.AccID)

		// Try a sample futures contract to get trading info
		// Use a common symbol - WTI crude oil futures
		sampleCode := "US.CL.0" // WTI Crude Oil Futures

		info, err := client.GetAccTradingInfo(ctx, cli, acc.AccID,
			constant.TrdMarket_FuturesSimulateUS, sampleCode,
			constant.OrderType_Normal, 75.0) // ~$75 price
		if err != nil {
			fmt.Printf("  GetAccTradingInfo failed: %v\n", err)
			continue
		}

		fmt.Printf("  MaxCashBuy: %.2f\n", info.MaxCashBuy)
		fmt.Printf("  MaxCashAndMarginBuy: %.2f\n", info.MaxCashAndMarginBuy)
		fmt.Printf("  MaxPositionSell: %.2f\n", info.MaxPositionSell)
		fmt.Printf("  MaxSellShort: %.2f\n", info.MaxSellShort)
		fmt.Printf("  MaxBuyBack: %.2f\n", info.MaxBuyBack)
		fmt.Printf("  LongRequiredIM: %.2f\n", info.LongRequiredIM)
		fmt.Printf("  ShortRequiredIM: %.2f\n", info.ShortRequiredIM)
	}

	fmt.Println("\n=== Note ===")
	fmt.Println("Futures accounts use GetAccTradingInfo for margin/cash info.")
	fmt.Println("GetFunds() only works for stock accounts, not futures accounts.")
}