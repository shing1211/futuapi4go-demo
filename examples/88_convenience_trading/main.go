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

	ctx := context.Background()

	fmt.Println("=== Convenience Trading API Demo ===")
	fmt.Println("Simplest possible one-liner patterns for common trading operations")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, mc.Client)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	var mkt constant.TrdMarket
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			for _, auth := range acc.TrdMarketAuthList {
				if auth == constant.TrdMarket_US.Int32() {
					accID = acc.AccID
					mkt = constant.TrdMarket_US
					break
				}
			}
		}
		if accID != 0 {
			break
		}
	}
	if accID == 0 {
		accID = accounts[0].AccID
		mkt = constant.TrdMarket(accounts[0].TrdMarketAuthList[0])
	}
	fmt.Printf("Using AccID=%d market=%v\n", accID, mkt)

	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, mc.Client, pwd); err != nil {
			fmt.Printf("UnlockTrading: %v\n", err)
		} else {
			fmt.Println("Trading unlocked")
		}
	}

	fmt.Println()
	fmt.Println("--- Pattern 1: One-Liner Account Funds ---")
	funds, err := client.GetFunds(ctx, mc.Client, accID)
	if err != nil {
		fmt.Printf("GetFunds: %v\n", err)
	} else {
		fmt.Printf("Cash=$%.2f Power=$%.2f Assets=$%.2f\n",
			funds.Cash, funds.Power, funds.TotalAssets)
	}

	fmt.Println()
	fmt.Println("--- Pattern 2: One-Liner Position Check ---")
	positions, err := client.GetPositionList(ctx, mc.Client, accID)
	if err != nil {
		fmt.Printf("GetPositionList: %v\n", err)
	} else if len(positions) == 0 {
		fmt.Println("No open positions")
	} else {
		for _, p := range positions {
			fmt.Printf("%s: Qty=%.0f Cost=$%.2f Cur=$%.2f P/L=$%.2f\n",
				p.Code, p.Quantity, p.CostPrice, p.CurPrice, p.PnL)
		}
	}

	fmt.Println()
	fmt.Println("--- Pattern 3: One-Liner Stock Quote ---")
	for _, sym := range []string{"NVDA", "AAPL"} {
		quote, err := client.GetQuote(ctx, mc.Client, constant.Market_US, sym)
		if err != nil {
			fmt.Printf("%s: %v\n", sym, err)
		} else {
			fmt.Printf("%-5s $%.2f (open=$%.2f high=$%.2f low=$%.2f vol=%d)\n",
				sym, quote.Price, quote.Open, quote.High, quote.Low, quote.Volume)
		}
	}

	fmt.Println()
	fmt.Println("--- Pattern 4: One-Liner Limit Buy (Display Only) ---")
	stock, price, qty, secMarket := "US.NVDA", 120.0, 1.0, constant.TrdSecMarket_US
	fmt.Printf("client.PlaceOrder(ctx, client, %d, %v, %q, Buy, Normal, $%.2f, %.0f, %v)\n",
		accID, mkt, stock, price, qty, secMarket)
	result, err := client.PlaceOrder(ctx, mc.Client, accID, mkt,
		stock, constant.TrdSide_Buy, constant.OrderType_Normal, price, qty, secMarket)
	if err != nil {
		fmt.Printf("  → %v (expected in simulate mode)\n", err)
	} else {
		fmt.Printf("  → OrderID=%d\n", result.OrderID)
	}

	fmt.Println()
	fmt.Println("--- Pattern 5: One-Liner Open Orders ---")
	orders, err := client.GetOrderList(ctx, mc.Client, accID)
	if err != nil {
		fmt.Printf("GetOrderList: %v\n", err)
	} else {
		fmt.Printf("%d open order(s)\n", len(orders))
		for _, o := range orders {
			fmt.Printf("  Order %d: %s side=%d qty=%.0f price=$%.2f status=%d\n",
				o.OrderID, o.Code, o.TrdSide, o.Qty, o.Price, o.OrderStatus)
		}
	}

	fmt.Println()
	fmt.Println("--- Pattern 6: Fluent TradeAPI (Alternative) ---")
	trdAPI := mc.Client.Trade()

	simpleReq, err := trdAPI.NewOrder(accID, mkt).
		Buy("US.AAPL", 10).
		At(175.0).
		AutoDetectMarket().
		Build()
	if err != nil {
		fmt.Printf("Build: %v\n", err)
	} else {
		resp, err := trdAPI.PlaceOrder(ctx, simpleReq)
		if err != nil {
			fmt.Printf("PlaceOrder: %v\n", err)
		} else {
			fmt.Printf("OrderID=%d\n", resp.OrderID)
		}
	}

	fmt.Println()
	fmt.Println("--- Quick Reference ---")
	fmt.Println("  client.PlaceOrder(ctx, c, accID, mkt, code, side, type, price, qty, secMkt)")
	fmt.Println("  client.GetFunds(ctx, c, accID)")
	fmt.Println("  client.GetPositionList(ctx, c, accID)")
	fmt.Println("  client.GetOrderList(ctx, c, accID)")
	fmt.Println("  client.GetQuote(ctx, c, market, code)")
	fmt.Println()
	fmt.Println("  trdAPI := c.Trade()")
	fmt.Println("  trdAPI.NewOrder(accID, mkt).Buy(code, qty).At(price).Build()")
	fmt.Println()
	fmt.Println("Note: Requires FUTU_TRADE_PWD for real execution.")

	fmt.Println("\n=== Convenience API Demo Complete ===")
}
