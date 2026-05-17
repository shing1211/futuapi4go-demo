package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Smart Money Flow Tracker ===")
	fmt.Println("Institutional Accumulation/Distribution Signal from 4 data sources")
	fmt.Println()

	symbols := []string{"NVDA", "AAPL", "TSLA"}

	for _, sym := range symbols {
		fmt.Printf("─── %s ───\n", sym)

		score := 0.0
		maxScore := 0.0
		var signals []string

		fmt.Println()
		fmt.Println("  Source 1: Capital Flow (Main Force)")
		flows, err := client.GetCapitalFlow(ctx, mc.Client, constant.Market_US, sym)
		if err != nil {
			fmt.Printf("    GetCapitalFlow: %v\n", err)
		} else if len(flows.Items) > 0 {
			f := flows.Items[len(flows.Items)-1]
			fmt.Printf("    Main InFlow: $%.0f | Super: $%.0f | Big: $%.0f | Mid: $%.0f | Small: $%.0f\n",
				f.MainInFlow, f.SuperInFlow, f.BigInFlow, f.MidInFlow, f.SmlInFlow)

			totalFlow := f.MainInFlow + f.SuperInFlow + f.BigInFlow
			if totalFlow > 0 {
				flowScore := minFloat(totalFlow/1e6, 10)
				score += flowScore
				signals = append(signals, fmt.Sprintf("Main flow +$%.0fM", totalFlow/1e6))
			} else {
				flowScore := maxFloat(totalFlow/1e6, -10)
				score += flowScore
				signals = append(signals, fmt.Sprintf("Main flow $%.0fM (outflow)", totalFlow/1e6))
			}
			maxScore += 10
		}

		fmt.Println()
		fmt.Println("  Source 2: Capital Distribution (Trade Size)")

		dist, err := client.GetCapitalDistribution(ctx, mc.Client, constant.Market_US, sym)
		if err != nil {
			fmt.Printf("    GetCapitalDistribution: %v\n", err)
		} else {
			largeIn := dist.MainInflow + dist.BigInflow
			largeOut := dist.MainOutflow + dist.BigOutflow
			totalLarge := largeIn + largeOut

			fmt.Printf("    Main: +$%.0f/-$%.0f | Big: +$%.0f/-$%.0f\n",
				dist.MainInflow, dist.MainOutflow, dist.BigInflow, dist.BigOutflow)
			fmt.Printf("    Mid: +$%.0f/-$%.0f | Small: +$%.0f/-$%.0f\n",
				dist.MidInflow, dist.MidOutflow, dist.SmallInflow, dist.SmallOutflow)

			largeRatio := 0.0
			if totalLarge > 0 {
				largeRatio = (largeIn - largeOut) / totalLarge * 100
				score += largeRatio / 10
				maxScore += 10
				dir := "net buying"
				if largeRatio < 0 {
					dir = "net selling"
				}
				signals = append(signals, fmt.Sprintf("Large trades %.1f%% (%s)", largeRatio, dir))
			}
		}

		fmt.Println()
		fmt.Println("  Source 3: Broker Queue")
		bids, asks, err := client.GetBroker(ctx, mc.Client, constant.Market_US, sym, 10)
		if err != nil {
			fmt.Printf("    GetBroker: %v\n", err)
		} else {
			bidVol := int64(0)
			for _, b := range bids {
				bidVol += b.Volume
			}
			askVol := int64(0)
			for _, a := range asks {
				askVol += a.Volume
			}
			fmt.Printf("    Bid brokers: %d (vol=%d) | Ask brokers: %d (vol=%d)\n",
				len(bids), bidVol, len(asks), askVol)

			totalBrk := bidVol + askVol
			if totalBrk > 0 {
				brkImbalance := float64(bidVol-askVol) / float64(totalBrk) * 100
				score += brkImbalance / 10
				maxScore += 10
				dir := "more aggressive on bid"
				if brkImbalance < 0 {
					dir = "more aggressive on ask"
				}
				signals = append(signals, fmt.Sprintf("Broker imbalance %.1f%% (%s)", brkImbalance, dir))
			}
		}

		fmt.Println()
		fmt.Println("  Source 4: Order Book Depth")
		if err := client.Subscribe(ctx, mc.Client, constant.Market_US, sym,
			[]constant.SubType{constant.SubType_OrderBook}); err != nil {
			fmt.Printf("    Subscribe: %v\n", err)
		} else {
			book, err := client.GetOrderBook(ctx, mc.Client, constant.Market_US, sym, 5)
			if err != nil {
				fmt.Printf("    GetOrderBook: %v\n", err)
			} else {
				bidDepth := int64(0)
				for _, b := range book.Bids {
					bidDepth += b.Volume
				}
				askDepth := int64(0)
				for _, a := range book.Asks {
					askDepth += a.Volume
				}
				fmt.Printf("    Bid depth: %d | Ask depth: %d\n", bidDepth, askDepth)

				totalDepth := bidDepth + askDepth
				if totalDepth > 0 {
					depthRatio := float64(bidDepth-askDepth) / float64(totalDepth) * 100
					score += depthRatio / 10
					maxScore += 10
					dir := "bid-heavy"
					if depthRatio < 0 {
						dir = "ask-heavy"
					}
					signals = append(signals, fmt.Sprintf("Depth imbalance %.1f%% (%s)", depthRatio, dir))
				}
			}
		}

		fmt.Println()
		finalScore := 0.0
		if maxScore > 0 {
			finalScore = (score / maxScore) * 100
		}
		finalScore = maxFloat(0, minFloat(finalScore, 100))

		arrow := "➡️"
		label := "Neutral"
		if finalScore >= 70 {
			arrow = "🟢"
			label = "Strong Accumulation"
		} else if finalScore >= 55 {
			arrow = "↗️"
			label = "Mild Accumulation"
		} else if finalScore <= 30 {
			arrow = "🔴"
			label = "Strong Distribution"
		} else if finalScore <= 45 {
			arrow = "↘️"
			label = "Mild Distribution"
		}
		fmt.Printf("  %s Smart Money Score: %.0f/100 — %s\n", arrow, finalScore, label)
		for _, s := range signals {
			fmt.Printf("    • %s\n", s)
		}
		fmt.Println()
	}

	fmt.Println("=== Smart Money Analysis Complete ===")
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
