package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

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

	fmt.Println("=== Momentum Scanner Demo (US Stocks) ===")
	fmt.Println("Scanning for momentum breakout candidates")
	fmt.Println()

	// Use StockFilter for momentum screening
	fmt.Println("--- Step 1: Screen with StockFilter ---")
	results, err := client.StockFilter(ctx, cli, constant.Market_US, 0, 30)
	if err != nil {
		fmt.Printf("StockFilter failed: %v\n", err)
		return
	}

	fmt.Printf("Found %d stocks from filter\n", len(results))

	// Get snapshots for enriched data
	fmt.Println("\n--- Step 2: Enrich with Snapshot Data ---")

	var symbols []string
	for _, r := range results {
		if r.StockName != "" {
			symbols = append(symbols, r.StockName)
		}
	}

	if len(symbols) > 10 {
		symbols = symbols[:10] // Limit to 10 for demo
	}

	snapshots, err := client.GetSecuritySnapshot(ctx, cli, securitiesFromNames(constant.Market_US, symbols))
	if err != nil {
		fmt.Printf("GetSecuritySnapshot failed: %v\n", err)
	} else {
		fmt.Printf("Got snapshots for %d stocks\n", len(snapshots))
	}

	// Calculate momentum scores
	fmt.Println("\n--- Step 3: Momentum Analysis ---")
	type MomentumStock struct {
		Code       string
		Name       string
		Price      float64
		ChangePct  float64
		Volume     int64
		High52Week float64
		Score      float64
	}

	var momentumStocks []MomentumStock

	for _, snap := range snapshots {
		if snap.LastPrice <= 0 || snap.High52Week <= 0 {
			continue
		}

		// Simple momentum score: (price vs 52w high) * volume factor
		proximityToHigh := snap.LastPrice / snap.High52Week
		volumeFactor := float64(snap.Volume) / 1000000 // Normalize

		// Score based on: approaching 52w high + reasonable volume
		score := proximityToHigh*50 + volumeFactor

		momentumStocks = append(momentumStocks, MomentumStock{
			Code:       snap.Code,
			Name:       snap.Name,
			Price:      snap.LastPrice,
			ChangePct:  (snap.LastPrice - snap.ClosePrevDay) / snap.ClosePrevDay * 100,
			Volume:     snap.Volume,
			High52Week: snap.High52Week,
			Score:      score,
		})
	}

	// Sort by momentum score
	sort.Slice(momentumStocks, func(i, j int) bool {
		return momentumStocks[i].Score > momentumStocks[j].Score
	})

	fmt.Printf("\n%-10s %-20s %10s %10s %15s %10s\n",
		"Symbol", "Name", "Price", "Change%", "Volume", "Score")
	fmt.Println("─────────────────────────────────────────────────────────────────")

	for i, stock := range momentumStocks {
		if i >= 10 {
			break
		}
		fmt.Printf("%-10s %-20s $%9.2f %+9.2f%% %15d %10.2f\n",
			stock.Code, stock.Name, stock.Price, stock.ChangePct, stock.Volume, stock.Score)
	}

	// Get K-lines for top candidate
	if len(momentumStocks) > 0 {
		topStock := momentumStocks[0].Code
		fmt.Printf("\n--- Step 4: K-Line Analysis for Top Candidate (%s) ---\n", topStock)

		klines, err := client.GetKLines(ctx, cli, constant.Market_US, topStock,
			constant.KLType_KDay, "", 10)
		if err != nil {
			fmt.Printf("GetKLines failed: %v\n", err)
		} else {
			fmt.Println("Recent 10 days:")
			for _, k := range klines {
				fmt.Printf("  %s: O=%.2f H=%.2f L=%.2f C=%.2f\n",
					k.Time, k.Open, k.High, k.Low, k.Close)
			}
		}
	}

	fmt.Println("\n=== Momentum Scanner Complete ===")
}

func securitiesFromNames(market constant.Market, names []string) []*constant.Security {
	var securities []*constant.Security
	for _, name := range names {
		securities = append(securities, &constant.Security{Market: market, Code: name})
	}
	return securities
}