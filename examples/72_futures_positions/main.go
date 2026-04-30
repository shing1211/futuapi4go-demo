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

	fmt.Println("=== Stock Positions (client.GetPositionList) ===")
	accounts, err := client.GetAccountList(ctx, cli)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	for _, acc := range accounts {
		if acc.TrdEnv != 0 {
			continue // skip real accounts
		}

		positions, err := client.GetPositionList(ctx, cli, acc.AccID)
		if err != nil {
			fmt.Printf("AccID %d: GetPositionList failed: %v\n", acc.AccID, err)
			continue
		}

		fmt.Printf("Stock AccID=%d: %d positions\n", acc.AccID, len(positions))
		if len(positions) == 0 {
			fmt.Println("  (no positions)")
			continue
		}
		for _, p := range positions {
			fmt.Printf("  %s %s: Qty=%.0f Cost=%.2f Cur=%.2f P/L=%.2f (%.2f%%)\n",
				p.Code, p.Name, p.Quantity, p.CostPrice, p.CurPrice, p.PnL, p.PnLRate*100)
		}
	}

	fmt.Println("\n=== Futures Positions (TradeAPI.GetPositionList) ===")
	resp, err := cli.Trade().GetAccList(ctx, constant.TrdCategory_Future, true)
	if err != nil {
		log.Fatalf("GetAccList(TrdCategory_Future) failed: %v", err)
	}

	// Markets to try for futures positions
	markets := []struct {
		name   string
		market constant.TrdMarket
	}{
		{"Futures", constant.TrdMarket_Futures},
		{"FuturesSimulateHK", constant.TrdMarket_FuturesSimulateHK},
		{"FuturesSimulateUS", constant.TrdMarket_FuturesSimulateUS},
		{"FuturesSimulateSG", constant.TrdMarket_FuturesSimulateSG},
		{"FuturesSimulateJP", constant.TrdMarket_FuturesSimulateJP},
	}

	for _, acc := range resp.AccList {
		env := "Real"
		if acc.TrdEnv == 0 {
			env = "Simulated"
		}

		fmt.Printf("\nFutures AccID=%d (TrdEnv=%s, Auth=%v):\n",
			acc.AccID, env, acc.TrdMarketAuthList)

		for _, m := range markets {
			req := &trd.GetPositionListRequest{
				AccID:     acc.AccID,
				TrdMarket: m.market,
				TrdEnv:    constant.TrdEnv(acc.TrdEnv),
			}

			posResp, err := cli.Trade().GetPositionList(ctx, req)
			if err != nil {
				continue
			}

			if len(posResp.PositionList) > 0 {
				fmt.Printf("  [%s] %d positions:\n", m.name, len(posResp.PositionList))
				for _, p := range posResp.PositionList {
					fmt.Printf("    %s %s: Qty=%.0f Cost=%.2f P/L=%.2f (%.2f%%)\n",
						p.Code, p.Name, p.Qty, p.CostPrice, p.PlVal, p.PlRatio*100)
				}
			}
		}
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("Stock positions: client.GetPositionList(accID)")
	fmt.Println("Futures positions: cli.Trade().GetPositionList with TrdMarket_Futures")
}