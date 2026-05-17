package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/history"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Bulk Historical Data Pipeline Demo ===")
	fmt.Println("Downloads K-line data with pagination, retry, and progress tracking")
	fmt.Println()

	fmt.Println("--- Single Symbol Download ---")
	d := history.NewDownloader(mc.Client,
		history.WithProgress(func(p history.DownloadProgress) {
			fmt.Printf("\r  Progress: %d bars", p.Downloaded)
		}))

	err := d.DownloadKLine(ctx, history.KLineRequest{
		Code:      "NVDA",
		Market:    constant.Market_US,
		KLType:    constant.KLType_K_Day,
		StartDate: "2025-01-01",
		EndDate:   "2026-05-15",
	})
	if err != nil {
		fmt.Printf("  DownloadKLine: %v\n", err)
	} else {
		fmt.Println("  ✓ Single symbol download complete")
	}

	fmt.Println()
	fmt.Println("--- Download with Statistics ---")
	d2 := history.NewDownloader(mc.Client,
		history.WithProgress(func(p history.DownloadProgress) {
			if p.Total > 0 {
				pct := float64(p.Downloaded) / float64(p.Total) * 100
				fmt.Printf("\r  Progress: %.1f%% (%d/%d, %.0f bars/s, ETA %v)",
					pct, p.Downloaded, p.Total, p.Speed, p.ETA)
			}
		}))

	bars, stats, err := d2.DownloadWithStats(ctx, history.KLineRequest{
		Code:      "AAPL",
		Market:    constant.Market_US,
		KLType:    constant.KLType_K_Day,
		StartDate: "2025-01-01",
		EndDate:   "2026-05-15",
	})
	if err != nil {
		fmt.Printf("  DownloadWithStats: %v\n", err)
	} else {
		fmt.Printf("  Downloaded: %d bars\n", stats.TotalBars)
		fmt.Printf("  Requests:   %d\n", stats.Requests)
		fmt.Printf("  Errors:     %d\n", stats.Errors)
		fmt.Printf("  Duration:   %v\n", stats.DownloadTime)

		if len(bars) > 0 {
			fmt.Println("\n  Sample bars (first 5):")
			fmt.Printf("  %-14s %-10s %-10s %-10s %-10s %s\n",
				"Date", "Open", "High", "Low", "Close", "Volume")
			fmt.Println("  " + repeat("─", 65))
			limit := 5
			if len(bars) < limit {
				limit = len(bars)
			}
			for _, bar := range bars[:limit] {
				fmt.Printf("  %-14s %-10.2f %-10.2f %-10.2f %-10.2f %d\n",
					bar.Time, bar.OpenPrice, bar.HighPrice,
					bar.LowPrice, bar.ClosePrice, bar.Volume)
			}
		}
	}

	fmt.Println()
	fmt.Println("--- Concurrent Multi-Symbol Download ---")
	cd := history.NewConcurrentDownloader(mc.Client,
		history.WithWorkers(3))

	requests := []history.KLineRequest{
		{Code: "NVDA", Market: constant.Market_US, KLType: constant.KLType_K_Day,
			StartDate: "2026-01-01", EndDate: "2026-05-15"},
		{Code: "AAPL", Market: constant.Market_US, KLType: constant.KLType_K_Day,
			StartDate: "2026-01-01", EndDate: "2026-05-15"},
		{Code: "TSLA", Market: constant.Market_US, KLType: constant.KLType_K_Day,
			StartDate: "2026-01-01", EndDate: "2026-05-15"},
	}

	results, err := cd.DownloadMultiple(ctx, requests)
	if err != nil {
		fmt.Printf("  DownloadMultiple: %v\n", err)
	} else {
		fmt.Println("  ✓ Multi-symbol download complete")
		fmt.Println()
		fmt.Printf("%-8s %-12s %-12s %s\n", "Symbol", "Bars", "Duration", "Errors")
		fmt.Println("  " + repeat("─", 45))
		for _, r := range results {
			if r.Error != nil {
				fmt.Printf("%-8s %-12s %-12s %v\n",
					r.Request.Code, "FAILED", "-", r.Error)
			} else {
				fmt.Printf("%-8s %-12d %-12v %d\n",
					r.Request.Code, r.Stats.TotalBars,
					r.Stats.DownloadTime.Round(time.Millisecond), r.Stats.Errors)
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Pipeline Complete ===")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(results)
}

func repeat(s string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = s[0]
	}
	return string(b)
}
