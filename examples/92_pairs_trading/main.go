package main

import (
	"context"
	"fmt"
	"math"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

type pairDef struct {
	Name   string
	CodeA  string
	CodeB  string
	Market constant.Market
}

type pairResult struct {
	Pair        pairDef
	Correlation float64
	SpreadZ     float64
	SpreadMean  float64
	SpreadStd   float64
	CurSpread   float64
	PriceA      float64
	PriceB      float64
	Signal      string
	Direction   string
}

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Pairs Trading / Statistical Arbitrage Scanner ===")
	fmt.Println("Detects divergent correlated pairs for mean-reversion trades")
	fmt.Println()

	pairs := []pairDef{
		{"NVDA/AMD", "NVDA", "AMD", constant.Market_US},
		{"AAPL/MSFT", "AAPL", "MSFT", constant.Market_US},
		{"TSLA/AMZN", "TSLA", "AMZN", constant.Market_US},
	}

	fmt.Printf("%-14s %-12s %-12s %-12s %-10s %s\n",
		"Pair", "Correlation", "Spread Z", "Cur Spread", "Signal", "Direction")
	fmt.Println(repeat("─", 80))

	var results []pairResult
	for _, p := range pairs {
		r := analyzePair(ctx, mc.Client, p)
		results = append(results, r)

		signal := "WAIT"
		dir := "—"
		if r.SpreadZ > 2.0 {
			signal = "SHORT"
			dir = fmt.Sprintf("Short %s / Long %s", p.CodeA, p.CodeB)
		} else if r.SpreadZ < -2.0 {
			signal = "LONG"
			dir = fmt.Sprintf("Long %s / Short %s", p.CodeA, p.CodeB)
		} else if r.SpreadZ > 1.5 {
			signal = "WATCH SHORT"
			dir = fmt.Sprintf("Monitor %s/%s", p.CodeA, p.CodeB)
		} else if r.SpreadZ < -1.5 {
			signal = "WATCH LONG"
			dir = fmt.Sprintf("Monitor %s/%s", p.CodeB, p.CodeA)
		}
		r.Signal = signal
		r.Direction = dir

		fmt.Printf("%-14s %-12.4f %-+12.2f %-12.4f %-10s %s\n",
			p.Name, r.Correlation, r.SpreadZ, r.CurSpread, signal, dir)
	}

	fmt.Println()
	for _, r := range results {
		if r.Signal == "SHORT" || r.Signal == "LONG" {
			fmt.Printf("--- %s Strategy ---\n", r.Pair.Name)
			fmt.Printf("  Correlation: %.4f (over 60 days)\n", r.Correlation)
			fmt.Printf("  Spread mean: $%.4f  std: $%.4f\n", r.SpreadMean, r.SpreadStd)
			fmt.Printf("  Current spread: $%.4f (z-score: %.2fσ)\n", r.CurSpread, r.SpreadZ)
			fmt.Printf("  Signal: %s — %s\n", r.Signal, r.Direction)
			fmt.Printf("  %s: $%.2f | %s: $%.2f\n",
				r.Pair.CodeA, r.PriceA, r.Pair.CodeB, r.PriceB)

			fmt.Printf("  Entry at spread $%.4f, target reversion to $%.4f\n", r.CurSpread, r.SpreadMean)
			fmt.Printf("  Stop if z-score exceeds ±3.0σ ($%.4f)\n", r.SpreadMean+3*r.SpreadStd)

			if r.SpreadZ > 2.0 {
				fmt.Printf("  → SHORT %s (overvalued relative to %s)\n", r.Pair.CodeA, r.Pair.CodeB)
				fmt.Printf("  → LONG %s (undervalued relative to %s)\n", r.Pair.CodeB, r.Pair.CodeA)
			} else {
				fmt.Printf("  → LONG %s (undervalued relative to %s)\n", r.Pair.CodeA, r.Pair.CodeB)
				fmt.Printf("  → SHORT %s (overvalued relative to %s)\n", r.Pair.CodeB, r.Pair.CodeA)
			}
			fmt.Println()
		}
	}

	fmt.Println("=== Pairs Analysis Complete ===")
	fmt.Println("Note: Pairs trading requires real trading environment.")
	fmt.Println("  Set FUTU_TRADE_PWD and use WithTradeEnv(TrdEnv_Real).")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(results)
}

func analyzePair(ctx context.Context, c *client.Client, p pairDef) pairResult {
	r := pairResult{Pair: p}

	klA, err := client.GetKLines(ctx, c, p.Market, p.CodeA, constant.KLType_K_Day, 60)
	if err != nil {
		fmt.Printf("  %s GetKLines(%s): %v\n", p.Name, p.CodeA, err)
		return r
	}
	klB, err := client.GetKLines(ctx, c, p.Market, p.CodeB, constant.KLType_K_Day, 60)
	if err != nil {
		fmt.Printf("  %s GetKLines(%s): %v\n", p.Name, p.CodeB, err)
		return r
	}

	n := min(len(klA), len(klB))
	if n < 5 {
		fmt.Printf("  %s: insufficient data (%d bars)\n", p.Name, n)
		return r
	}
	klA, klB = klA[:n], klB[:n]

	pricesA := make([]float64, n)
	pricesB := make([]float64, n)
	spreads := make([]float64, n)
	for i := 0; i < n; i++ {
		pricesA[i] = klA[i].Close
		pricesB[i] = klB[i].Close
		spreads[i] = pricesA[i] - pricesB[i]
	}

	r.Correlation = pearson(pricesA, pricesB)
	r.SpreadMean = mean(spreads)
	r.SpreadStd = std(spreads, r.SpreadMean)
	r.CurSpread = spreads[n-1]
	r.PriceA = pricesA[n-1]
	r.PriceB = pricesB[n-1]

	if r.SpreadStd > 0 {
		r.SpreadZ = (r.CurSpread - r.SpreadMean) / r.SpreadStd
	}

	return r
}

func pearson(x, y []float64) float64 {
	n := len(x)
	if n < 2 {
		return 0
	}
	mx, my := mean(x), mean(y)
	var num, dx2, dy2 float64
	for i := 0; i < n; i++ {
		dx := x[i] - mx
		dy := y[i] - my
		num += dx * dy
		dx2 += dx * dx
		dy2 += dy * dy
	}
	denom := math.Sqrt(dx2 * dy2)
	if denom == 0 {
		return 0
	}
	return num / denom
}

func mean(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func std(v []float64, m float64) float64 {
	s := 0.0
	for _, x := range v {
		d := x - m
		s += d * d
	}
	return math.Sqrt(s / float64(len(v)))
}

func repeat(s string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = s[0]
	}
	return string(b)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
