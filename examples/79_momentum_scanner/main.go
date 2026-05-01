package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
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

	fmt.Println("--- Step 1: Get snapshots for top US stocks ---")

	topStocks := []string{"AAPL", "TSLA", "NVDA", "MSFT", "AMZN", "GOOGL", "META", "AMAT", "MU", "INTC"}
	var secs []*qotcommon.Security
	for _, code := range topStocks {
		marketPtr := int32(constant.Market_US)
		sec := &qotcommon.Security{
			Market: &marketPtr,
			Code:   &code,
		}
		secs = append(secs, sec)
	}

	snapshots, err := client.GetSecuritySnapshot(ctx, cli, secs)
	if err != nil {
		fmt.Printf("GetSecuritySnapshot failed: %v\n", err)
		return
	}

	fmt.Printf("Got snapshots for %d stocks\n", len(snapshots))

	fmt.Println("\n--- Step 2: Momentum Analysis ---")
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
		if snap.CurPrice <= 0 || snap.Highest52WeeksPrice <= 0 {
			continue
		}

		proximityToHigh := snap.CurPrice / snap.Highest52WeeksPrice
		volumeFactor := float64(snap.Volume) / 1000000

		score := proximityToHigh*50 + volumeFactor

		momentumStocks = append(momentumStocks, MomentumStock{
			Code:       snap.Security.GetCode(),
			Name:       snap.Name,
			Price:      snap.CurPrice,
			ChangePct:  snap.ChangeVal / snap.Highest52WeeksPrice * 100,
			Volume:     snap.Volume,
			High52Week: snap.Highest52WeeksPrice,
			Score:      score,
		})
	}

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

	if len(momentumStocks) > 0 {
		topStock := momentumStocks[0].Code
		fmt.Printf("\n--- Step 3: K-Line Analysis for Top Candidate (%s) ---\n", topStock)

		if err := client.Subscribe(ctx, cli, constant.Market_US, topStock,
			[]constant.SubType{constant.SubType_K_Day}); err != nil {
			fmt.Printf("Subscribe failed: %v\n", err)
		} else {
			klines, err := client.GetKLines(ctx, cli, constant.Market_US, topStock,
				constant.KLType_K_Day, 10)
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
	}

	fmt.Println("\n=== Momentum Scanner Complete ===")
}