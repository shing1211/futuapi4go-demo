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

func main() {
	fmt.Println("═══════════════════════════════════════════")
	fmt.Println("         Quote Snapshot: US.NVDA")
	fmt.Println("═══════════════════════════════════════════")
	fmt.Printf("Timestamp: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	quote, err := client.GetQuote(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetQuote failed: %v", err)
	}

	fmt.Println("── Quote Data ────────────────────────────")
	fmt.Printf("  Symbol:      %s\n", quote.Symbol)
	fmt.Printf("  Name:        %s\n", quote.Name)
	fmt.Printf("  Price:       %.3f\n", quote.Price)
	fmt.Printf("  Open:        %.3f\n", quote.Open)
	fmt.Printf("  High:        %.3f\n", quote.High)
	fmt.Printf("  Low:         %.3f\n", quote.Low)
	fmt.Printf("  LastClose:   %.3f\n", quote.LastClose)
	if quote.LastClose > 0 {
		change := quote.Price - quote.LastClose
		pct := change / quote.LastClose * 100
		fmt.Printf("  Change:      %+.3f (%+.2f%%)\n", change, pct)
	}
	fmt.Printf("  Volume:      %d\n", quote.Volume)
	fmt.Printf("  Turnover:    %.0f\n", quote.Turnover)
	fmt.Printf("  TurnoverRate: %.4f%%\n", quote.TurnoverRate)
	fmt.Printf("  Amplitude:   %.2f%%\n", quote.Amplitude)
	fmt.Printf("  PriceSpread: %.4f\n", quote.PriceSpread)
	fmt.Printf("  Suspended:   %v\n", quote.IsSuspended)
	fmt.Printf("  DarkStatus:  %d\n", quote.DarkStatus)
	fmt.Printf("  SecStatus:   %d\n", quote.SecStatus)
	fmt.Printf("  ListTime:    %s\n", quote.ListTime)
	fmt.Printf("  UpdateTime:  %s\n", quote.Timestamp)

	if quote.PreMarket != nil {
		fmt.Println()
		fmt.Println("── Pre-Market Data ────────────────────────")
		pm := quote.PreMarket
		fmt.Printf("  Price:      %.3f\n", pm.Price)
		fmt.Printf("  High:       %.3f\n", pm.HighPrice)
		fmt.Printf("  Low:        %.3f\n", pm.LowPrice)
		fmt.Printf("  Volume:     %d\n", pm.Volume)
		fmt.Printf("  Turnover:   %.0f\n", pm.Turnover)
		fmt.Printf("  Change:     %+.3f\n", pm.ChangeVal)
		fmt.Printf("  Change%%:    %+.2f%%\n", pm.ChangeRate)
		fmt.Printf("  Amplitude:  %.2f%%\n", pm.Amplitude)
	}
	if quote.AfterMarket != nil {
		fmt.Println()
		fmt.Println("── After-Market Data ──────────────────────")
		am := quote.AfterMarket
		fmt.Printf("  Price:      %.3f\n", am.Price)
		fmt.Printf("  High:       %.3f\n", am.HighPrice)
		fmt.Printf("  Low:        %.3f\n", am.LowPrice)
		fmt.Printf("  Volume:     %d\n", am.Volume)
		fmt.Printf("  Turnover:   %.0f\n", am.Turnover)
		fmt.Printf("  Change:     %+.3f\n", am.ChangeVal)
		fmt.Printf("  Change%%:    %+.2f%%\n", am.ChangeRate)
		fmt.Printf("  Amplitude:  %.2f%%\n", am.Amplitude)
	}
	if quote.Overnight != nil {
		fmt.Println()
		fmt.Println("── Overnight Data ─────────────────────────")
		on := quote.Overnight
		fmt.Printf("  Price:      %.3f\n", on.Price)
		fmt.Printf("  Volume:     %d\n", on.Volume)
		fmt.Printf("  Turnover:   %.0f\n", on.Turnover)
		fmt.Printf("  Change:     %+.3f\n", on.ChangeVal)
		fmt.Printf("  Change%%:    %+.2f%%\n", on.ChangeRate)
	}
	if quote.OptionExData != nil {
		fmt.Println("── Option Extended Data ─────────────────────")
		o := quote.OptionExData
		fmt.Printf("  StrikePrice:  %.2f\n", o.GetStrikePrice())
		fmt.Printf("  ContractSize: %d\n", o.GetContractSize())
		fmt.Printf("  ImpliedVol:   %.4f\n", o.GetImpliedVolatility())
		fmt.Printf("  Delta:        %.4f\n", o.GetDelta())
		fmt.Printf("  Gamma:        %.4f\n", o.GetGamma())
		fmt.Printf("  Vega:         %.4f\n", o.GetVega())
		fmt.Printf("  Theta:        %.4f\n", o.GetTheta())
		fmt.Printf("  Rho:          %.4f\n", o.GetRho())
	}
	if quote.FutureExData != nil {
		fmt.Println("── Future Extended Data ─────────────────────")
		f := quote.FutureExData
		fmt.Printf("  LastSettlePrice:    %.2f\n", f.GetLastSettlePrice())
		fmt.Printf("  Position:           %d\n", f.GetPosition())
		fmt.Printf("  PositionChange:     %d\n", f.GetPositionChange())
		fmt.Printf("  ExpiryDateDistance: %d\n", f.GetExpiryDateDistance())
	}
	if quote.WarrantExData != nil {
		fmt.Println("── Warrant Extended Data ────────────────────")
		w := quote.WarrantExData
		fmt.Printf("  Delta:          %.4f\n", w.GetDelta())
		fmt.Printf("  ImpliedVol:     %.4f\n", w.GetImpliedVolatility())
		fmt.Printf("  Premium:        %.4f\n", w.GetPremium())
	}
	fmt.Println()

	fmt.Println("── Quote Data (JSON) ─────────────────────")
	display.PrintJSON(quote)
}
