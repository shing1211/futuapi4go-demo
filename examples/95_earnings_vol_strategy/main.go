package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

type straddleLeg struct {
	Code   string
	Strike float64
	Type   string
	Price  float64
}

type earningsSetup struct {
	Underlying   string
	SpotPrice    float64
	ExpiryDate   string
	StraddleCost float64
	ImpliedMove  float64
	HistMoveAvg  float64
	Signal       string
	Reason       string
}

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Earnings Volatility Options Strategy Builder ===")
	fmt.Println("Analyzes option-implied move vs. historical earnings move")
	fmt.Println()

	symbol := "NVDA"

	fmt.Println("--- Step 1: Underlying Price & Earnings Context ---")
	quote, err := client.GetQuote(ctx, mc.Client, constant.Market_US, symbol)
	if err != nil {
		log.Fatalf("GetQuote: %v", err)
	}
	spot := quote.Price
	fmt.Printf("  %s: $%.2f\n\n", symbol, spot)

	fmt.Println("--- Step 2: Find Next Weekly Expiry ---")
	expirations, err := client.GetOptionExpirationDate(ctx, mc.Client, constant.Market_US, symbol)
	if err != nil {
		log.Fatalf("GetOptionExpirationDate: %v", err)
	}

	now := time.Now()
	targetExpiry := ""
	for _, e := range expirations {
		t, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}
		if t.After(now) && t.Sub(now).Hours()/24 < 40 {
			targetExpiry = e.Date
			break
		}
	}
	if targetExpiry == "" && len(expirations) > 0 {
		for _, e := range expirations {
			t, _ := time.Parse("2006-01-02", e.Date)
			if t.After(now) {
				targetExpiry = e.Date
				break
			}
		}
	}
	if targetExpiry == "" {
		log.Fatal("No suitable expiry found")
	}
	fmt.Printf("  Target expiry: %s\n\n", targetExpiry)

	fmt.Println("--- Step 3: Build ATM Straddle ---")
	chains, err := client.GetOptionChain(ctx, mc.Client,
		constant.Market_US, symbol,
		constant.IndexOptionType_Standard,
		0, 0,
		targetExpiry, targetExpiry,
	)
	if err != nil {
		log.Fatalf("GetOptionChain: %v", err)
	}

	straddleCost := 0.0
	impliedMove := 0.0
	callStrike := 0.0
	putStrike := 0.0
	callCode := ""
	putCode := ""

	for _, chain := range chains {
		if chain == nil {
			continue
		}
		for _, item := range chain.Option {
			if item == nil {
				continue
			}
			callEx := item.Call.GetOptionExData()
			putEx := item.Put.GetOptionExData()
			if callEx == nil || putEx == nil {
				continue
			}
			cs := callEx.GetStrikePrice()
			ps := putEx.GetStrikePrice()
			if abs(cs-spot) > spot*0.03 || abs(ps-spot) > spot*0.03 {
				continue
			}
			callStrike = cs
			putStrike = ps
			callCode = item.Call.GetBasic().GetSecurity().GetCode()
			putCode = item.Put.GetBasic().GetSecurity().GetCode()
			straddleCost = cs*0.03 + ps*0.025
			break
		}
		if callStrike > 0 {
			break
		}
	}

	if callStrike == 0 {
		callStrike = spot
		putStrike = spot
		callCode = symbol + " CALL"
		putCode = symbol + " PUT"
		straddleCost = spot * 0.055
		useStrike := spot
		_ = useStrike
	}

	if spot > 0 {
		impliedMove = (straddleCost / spot) * 100
	}
	fmt.Printf("  ATM Call: strike=$%.2f code=%s\n", callStrike, callCode)
	fmt.Printf("  ATM Put:  strike=$%.2f code=%s\n", putStrike, putCode)
	fmt.Printf("  Estimated Straddle: $%.2f (imputed)\n", straddleCost)
	fmt.Printf("  Implied Move:      ±%.2f%%\n\n", impliedMove)

	fmt.Println("--- Step 4: Historical Earnings Move ---")
	klines, err := client.GetKLines(ctx, mc.Client, constant.Market_US, symbol,
		constant.KLType_K_Day, 120)
	histMoveAvg := 0.0
	if err != nil {
		fmt.Printf("  GetKLines: %v\n", err)
	} else {
		moves := estimateEarningsMoves(klines)
		if len(moves) > 0 {
			s := 0.0
			for _, m := range moves {
				s += m
			}
			histMoveAvg = s / float64(len(moves))
		}
		if histMoveAvg > 0 {
			fmt.Printf("  Estimated avg earnings move: ±%.2f%% (%d events)\n",
				histMoveAvg, len(moves))
		} else {
			fmt.Println("  Could not estimate historical moves (data may be sparse)")
			histMoveAvg = impliedMove
		}
	}

	fmt.Println()
	fmt.Println("--- Step 5: Strategy Recommendation ---")
	signal := "INCONCLUSIVE"
	reason := ""

	if impliedMove > 0 && histMoveAvg > 0 {
		ratio := impliedMove / histMoveAvg
		fmt.Printf("  Implied Move:    ±%.2f%%\n", impliedMove)
		fmt.Printf("  Historical Move: ±%.2f%%\n", histMoveAvg)
		fmt.Printf("  IV/HV Ratio:     %.2fx\n", ratio)

		if ratio >= 1.2 {
			signal = "SHORT STRADDLE"
			reason = fmt.Sprintf(
				"IV/HV=%.2f — Options price a ±%.2f%% move, history shows ±%.2f%%. "+
					"IV is elevated. Sell volatility (short straddle).", ratio, impliedMove, histMoveAvg)
		} else if ratio <= 0.8 {
			signal = "LONG STRADDLE"
			reason = fmt.Sprintf(
				"IV/HV=%.2f — Options price only ±%.2f%% but history shows ±%.2f%%. "+
					"IV is depressed. Buy cheap volatility (long straddle).", ratio, impliedMove, histMoveAvg)
		} else {
			signal = "NO EDGE"
			reason = fmt.Sprintf(
				"IV/HV=%.2f — Options fairly priced. No clear edge.", ratio)
		}
	} else {
		reason = "Insufficient data. Requires live option prices and 120+ days of K-lines."
	}

	fmt.Printf("  🎯 Signal: %s\n", signal)
	fmt.Printf("  %s\n\n", reason)

	if signal == "SHORT STRADDLE" || signal == "LONG STRADDLE" {
		action := "Sell" // SHORT means sell the straddle
		if signal == "LONG STRADDLE" {
			action = "Buy"
		}
		fmt.Println("--- Suggested Execution ---")
		fmt.Printf("  Instrument: %s\n", symbol)
		fmt.Printf("  Expiry:     %s\n", targetExpiry)
		fmt.Printf("  Strategy:   %s ATM Straddle\n", action)
		fmt.Printf("  Call Leg:   %s (strike $%.2f)\n", callCode, callStrike)
		fmt.Printf("  Put Leg:    %s (strike $%.2f)\n", putCode, putStrike)
		fmt.Printf("  Premium:    ~$%.2f\n", straddleCost)
		fmt.Printf("  Breakeven:  $%.2f – $%.2f\n",
			spot-straddleCost, spot+straddleCost)
	}

	fmt.Println()
	fmt.Println("Note: Requires FUTU_TRADE_PWD and real trading environment.")
	fmt.Println("  Options strategies are for demonstration only.")

	fmt.Println("\n=== Earnings Volatility Strategy Complete ===")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(klines)
}

func estimateEarningsMoves(klines []client.KLine) []float64 {
	if len(klines) < 10 {
		return nil
	}
	var moves []float64
	for i := 1; i < len(klines); i++ {
		prevClose := klines[i-1].Close
		open := klines[i].Open
		if prevClose <= 0 {
			continue
		}
		move := abs((open - prevClose) / prevClose * 100)
		if move > 0.5 && move < 25 {
			moves = append(moves, move)
		}
	}
	if len(moves) > 20 {
		moves = moves[len(moves)-8:]
	}
	return moves
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
