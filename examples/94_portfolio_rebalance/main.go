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

type targetHolding struct {
	Code     string
	Name     string
	TargetPct float64
	Market   constant.Market
}

type positionStatus struct {
	Code       string
	Name       string
	Qty        float64
	Price      float64
	Value      float64
	ActualPct  float64
	TargetPct  float64
	DriftPct   float64
	Action     string
	TradeValue float64
}

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Multi-Asset Portfolio Rebalancer ===")
	fmt.Println("Target allocation → current drift → rebalance plan → execute")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, mc.Client)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			for _, auth := range acc.TrdMarketAuthList {
				if auth == constant.TrdMarket_US.Int32() {
					accID = acc.AccID
					break
				}
			}
		}
		if accID != 0 {
			break
		}
	}

	if accID == 0 {
		for _, acc := range accounts {
			if acc.TrdEnv == 0 {
				accID = acc.AccID
				break
			}
		}
	}
	if accID == 0 {
		accID = accounts[0].AccID
	}
	fmt.Printf("Using AccID=%d\n", accID)

	targets := []targetHolding{
		{"NVDA", "NVIDIA Corp", 50, constant.Market_US},
		{"AAPL", "Apple Inc", 30, constant.Market_US},
	}

	fmt.Println()
	fmt.Println("--- Step 1: Target Allocation ---")
	fmt.Printf("%-8s %-20s %12s\n", "Symbol", "Name", "Target")
	fmt.Println(repeat("─", 42))
	totalTarget := 0.0
	for _, t := range targets {
		fmt.Printf("%-8s %-20s %11.0f%%\n", t.Code, t.Name, t.TargetPct)
		totalTarget += t.TargetPct
	}
	var cashTargetPct float64
	if totalTarget < 100 {
		cashTargetPct = 100 - totalTarget
		fmt.Printf("%-8s %-20s %11.0f%%\n", "CASH", "Cash Reserve", cashTargetPct)
	}

	fmt.Println()
	fmt.Println("--- Step 2: Current Portfolio ---")
	funds, err := client.GetFunds(ctx, mc.Client, accID)
	if err != nil {
		log.Fatalf("GetFunds: %v", err)
	}
	fmt.Printf("Total Assets: $%.2f | Cash: $%.2f | Power: $%.2f\n",
		funds.TotalAssets, funds.Cash, funds.Power)

	positions, err := client.GetPositionList(ctx, mc.Client, accID)
	if err != nil {
		fmt.Printf("GetPositionList: %v (will use empty portfolio)\n", err)
		positions = nil
	}

	posMap := make(map[string]float64)
	for _, p := range positions {
		posMap[p.Code] = p.Quantity
	}

	var statuses []positionStatus
	portfolioValue := funds.Cash
	for _, t := range targets {
		qty := posMap[t.Code]
		prices := []float64{0}
		if t.Market == constant.Market_US {
			quote, err := client.GetQuote(ctx, mc.Client, t.Market, t.Code)
			if err == nil {
				prices[0] = quote.Price
			}
		}
		price := prices[0]
		value := qty * price
		portfolioValue += value

		statuses = append(statuses, positionStatus{
			Code: t.Code, Name: t.Name, Qty: qty,
			Price: price, Value: value, TargetPct: t.TargetPct,
		})
	}

	fmt.Println()
	fmt.Printf("%-8s %-8s %-10s %-12s %-12s %-10s\n",
		"Symbol", "Qty", "Price", "Value", "Actual%", "Target%")
	fmt.Println(repeat("─", 65))
	totalPct := 0.0
	for i := range statuses {
		s := &statuses[i]
		if portfolioValue > 0 {
			s.ActualPct = s.Value / portfolioValue * 100
		}
		totalPct += s.ActualPct
		fmt.Printf("%-8s %-8.0f $%-8.2f $%-10.2f %-10.1f%% %-10.0f%%\n",
			s.Code, s.Qty, s.Price, s.Value, s.ActualPct, s.TargetPct)
	}
	if cashTargetPct > 0 {
		cashActualPct := funds.Cash / portfolioValue * 100
		totalPct += cashActualPct
		fmt.Printf("%-8s %-8s %-10s $%-10.2f %-10.1f%% %-10.0f%%\n",
			"CASH", "—", "—", funds.Cash, cashActualPct, cashTargetPct)
	}

	fmt.Println()
	fmt.Println("--- Step 3: Drift Analysis ---")
	fmt.Printf("%-8s %-10s %-10s %-10s %s\n",
		"Symbol", "Target%", "Actual%", "Drift", "Action")
	fmt.Println(repeat("─", 50))

	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, mc.Client, pwd); err != nil {
			fmt.Printf("UnlockTrading: %v\n", err)
		}
	}

	hasAction := false
	trdAPI := mc.Client.Trade()

	for i := range statuses {
		s := &statuses[i]
		s.DriftPct = s.ActualPct - s.TargetPct

		if s.DriftPct > 2.0 {
			s.Action = "SELL"
			s.TradeValue = (s.DriftPct / 100) * portfolioValue
			s.TradeValue = minFloat(s.TradeValue, s.Value)

			secMarket := constant.TrdSecMarket_US
			if s.Code == "00700" {
				secMarket = constant.TrdSecMarket_HK
			}
			_ = secMarket

			fmt.Printf("%-8s %-8.1f%% %-8.1f%% %+8.1f%% 🔴 SELL $%.0f\n",
				s.Code, s.TargetPct, s.ActualPct, s.DriftPct, s.TradeValue)
			hasAction = true

			if pwd != "" {
				sellQty := s.TradeValue / s.Price
				if sellQty > 0 {
					req, err := trdAPI.NewOrder(accID, constant.TrdMarket_US).
						Sell("US."+s.Code, sellQty).
						Market().
						AutoDetectMarket().
						Build()
					if err == nil {
						resp, err := trdAPI.PlaceOrder(ctx, req)
						if err != nil {
							fmt.Printf("    → PlaceOrder: %v\n", err)
						} else {
							fmt.Printf("    → Sell order placed: OrderID=%d\n", resp.OrderID)
						}
					} else {
						fmt.Printf("    → Build: %v\n", err)
					}
				}
			}

		} else if s.DriftPct < -2.0 {
			s.Action = "BUY"
			s.TradeValue = (-s.DriftPct / 100) * portfolioValue
			s.TradeValue = minFloat(s.TradeValue, funds.Power)

			fmt.Printf("%-8s %-8.1f%% %-8.1f%% %+8.1f%% 🟢 BUY $%.0f\n",
				s.Code, s.TargetPct, s.ActualPct, s.DriftPct, s.TradeValue)
			hasAction = true

			if pwd != "" {
				buyQty := s.TradeValue / s.Price
				if buyQty > 0 {
					req, err := trdAPI.NewOrder(accID, constant.TrdMarket_US).
						Buy("US."+s.Code, buyQty).
						Market().
						AutoDetectMarket().
						Build()
					if err == nil {
						resp, err := trdAPI.PlaceOrder(ctx, req)
						if err != nil {
							fmt.Printf("    → PlaceOrder: %v\n", err)
						} else {
							fmt.Printf("    → Buy order placed: OrderID=%d\n", resp.OrderID)
						}
					}
				}
			}

		} else {
			fmt.Printf("%-8s %-8.1f%% %-8.1f%% %+8.1f%% ✅ In range\n",
				s.Code, s.TargetPct, s.ActualPct, s.DriftPct)
		}
	}

	if cashTargetPct > 0 {
		cashDrift := (funds.Cash / portfolioValue * 100) - cashTargetPct
		if cashDrift > 5 {
			fmt.Printf("%-8s %-8.0f%% %-8.1f%% %+8.1f%% 🔴 Deploy cash\n",
				"CASH", cashTargetPct, funds.Cash/portfolioValue*100, cashDrift)
		} else if cashDrift < -5 {
			fmt.Printf("%-8s %-8.0f%% %-8.1f%% %+8.1f%% 🟢 Raise cash\n",
				"CASH", cashTargetPct, funds.Cash/portfolioValue*100, cashDrift)
		} else {
			fmt.Printf("%-8s %-8.0f%% %-8.1f%% %+8.1f%% ✅ In range\n",
				"CASH", cashTargetPct, funds.Cash/portfolioValue*100, cashDrift)
		}
	}

	if funds.IsPDT {
		fmt.Println()
		fmt.Println("─── PDT Advisory ───")
		fmt.Printf("  PDT Account: Yes (remaining DT: %s)\n", funds.PDTSeq)
		fmt.Printf("  Remaining DTBP: $%.2f\n", funds.RemainingDTBP)
		if funds.RemainingDTBP <= 0 {
			fmt.Println("  ⚠️  Day trade buying power exhausted — limit day trades today")
		}
	}

	if !hasAction {
		fmt.Println("  Portfolio within tolerance — no rebalance needed")
	}

	fmt.Println()
	if pwd == "" {
		fmt.Println("Note: Set FUTU_TRADE_PWD to enable order execution.")
	}
	fmt.Println("=== Rebalance Complete ===")
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func repeat(s string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = s[0]
	}
	return string(b)
}
