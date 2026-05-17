package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Account Risk & Margin Analyzer ===")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, mc.Client)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	isUS := false
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			for _, auth := range acc.TrdMarketAuthList {
				if auth == constant.TrdMarket_US.Int32() {
					accID = acc.AccID
					isUS = true
					break
				}
			}
		}
		if accID != 0 {
			break
		}
	}
	if !isUS && accID == 0 {
		for _, acc := range accounts {
			if acc.TrdEnv == 0 {
				for _, auth := range acc.TrdMarketAuthList {
					if auth == constant.TrdMarket_HK.Int32() {
						accID = acc.AccID
						break
					}
				}
			}
			if accID != 0 {
				break
			}
		}
	}
	if accID == 0 {
		accID = accounts[0].AccID
	}
	fmt.Printf("Using AccID=%d\n\n", accID)

	funds, err := client.GetFunds(ctx, mc.Client, accID)
	if err != nil {
		log.Fatalf("GetFunds failed: %v", err)
	}

	fmt.Println("=== Account Risk Profile ===")
	marginPct := 0.0
	if funds.TotalAssets > 0 {
		marginPct = funds.MaintenanceMargin / funds.TotalAssets * 100
	}

	fmt.Printf("%-25s $%13.2f\n", "Total Assets", funds.TotalAssets)
	fmt.Printf("%-25s $%13.2f\n", "Cash", funds.Cash)
	fmt.Printf("%-25s $%13.2f\n", "Market Value", funds.MarketVal)
	fmt.Printf("%-25s $%13.2f\n", "Buying Power (Long)", funds.Power)
	fmt.Printf("%-25s $%13.2f\n", "Buying Power (Short)", funds.MaxPowerShort)
	fmt.Printf("%-25s $%13.2f\n", "Initial Margin", funds.InitialMargin)
	fmt.Printf("%-25s $%13.2f\n", "Maintenance Margin", funds.MaintenanceMargin)
	fmt.Printf("%-25s $%13.2f (%.1f%% of assets)\n", "Margin Utilization", funds.MaintenanceMargin, marginPct)
	fmt.Printf("%-25s $%13.2f\n", "Margin Call Amount", funds.MarginCallMargin)
	fmt.Printf("%-25s $%13.2f\n", "Frozen Cash", funds.FrozenCash)
	fmt.Printf("%-25s $%13.2f\n", "Debt Cash", funds.DebtCash)
	fmt.Printf("%-25s $%13.2f\n", "Unrealized P/L", funds.UnrealizedPL)
	fmt.Printf("%-25s $%13.2f\n", "Avail. Withdrawal", funds.AvlWithdrawalCash)

	riskStr := "Safe"
	switch funds.RiskLevel {
	case 1:
		riskStr = "⚠️ Warning"
	case 2:
		riskStr = "🔴 Danger"
	case 3:
		riskStr = "✅ Absolute Safe"
	case 4:
		riskStr = "🔴 Danger (Options)"
	}
	fmt.Printf("%-25s %s (Level %d)\n", "Risk Level", riskStr, funds.RiskLevel)

	fmt.Println()
	fmt.Println("Risk Status Breakdown:")
	fmt.Printf("  Margin Utilization: %.1f%%", marginPct)
	if marginPct > 80 {
		fmt.Print(" 🔴 CRITICAL — margin call risk")
	} else if marginPct > 50 {
		fmt.Print(" ⚠️ Elevated")
	} else {
		fmt.Print(" ✅ Safe")
	}
	fmt.Println()

	if funds.MarginCallMargin > 0 {
		fmt.Printf("  Margin Call: $%.2f 🔴 Pending\n", funds.MarginCallMargin)
	} else {
		fmt.Println("  Margin Call: $0.00 ✅ None")
	}

	if isUS {
		fmt.Println()
		fmt.Println("─── PDT Status ───")
		if funds.IsPDT {
			fmt.Printf("%-25s Yes\n", "PDT Account")
			fmt.Printf("%-25s %s\n", "Remaining Day Trades", funds.PDTSeq)
			fmt.Printf("%-25s $%.2f\n", "Initial DTBP", funds.BeginningDTBP)
			fmt.Printf("%-25s $%.2f\n", "Remaining DTBP", funds.RemainingDTBP)
			fmt.Printf("%-25s $%.2f\n", "DT Call Amount", funds.DtCallAmount)

			dtStr := "Unknown"
			switch funds.DtStatus {
			case 1:
				dtStr = "Unlimited"
			case 2:
				dtStr = "EM Call (needs $25k+)"
			case 3:
				dtStr = "DT Call pending"
			}
			fmt.Printf("%-25s %s\n", "DT Status", dtStr)

			if funds.DtCallAmount > 0 {
				fmt.Println(" 🔴 DT Call pending — deposit funds immediately")
			}
			if funds.RemainingDTBP == 0 {
				fmt.Println(" 🔴 Day trade limit exhausted — cannot day trade")
			}
		} else {
			fmt.Printf("%-25s No\n", "PDT Account")
			fmt.Println(" (Account < $25k or non-US broker)")
		}
	} else {
		fmt.Println()
		fmt.Println("─── PDT Status ───")
		fmt.Println("PDT: Not applicable (HK market)")
	}

	fmt.Println()
	fmt.Println("─── Margin Details (Sample: NVDA) ───")
	sec := &qotcommon.Security{
		Market: ptrInt32(int32(constant.Market_US)),
		Code:   ptrStr("NVDA"),
	}
	marginRatios, err := client.GetMarginRatio(ctx, mc.Client, accID, constant.TrdMarket_US, []*qotcommon.Security{sec})
	if err != nil {
		fmt.Printf("GetMarginRatio: %v\n", err)
	} else {
		for _, mr := range marginRatios {
			fmt.Printf("%-25s %.1f%%\n", "IM Long Ratio", mr.ImLongRatio*100)
			fmt.Printf("%-25s %.1f%%\n", "IM Short Ratio", mr.ImShortRatio*100)
			_ = mr.IsLongPermit
			_ = mr.IsShortPermit
			fmt.Printf("%-25s %.4f%%\n", "Short Fee Rate", mr.ShortFeeRate*100)
		}
	}

	fmt.Println()
	fmt.Println("─── Max Trade Quantities (NVDA @ $180) ───")
	maxQtys, err := client.GetMaxTrdQtys(ctx, mc.Client, accID, constant.TrdMarket_US,
		"US.NVDA", constant.OrderType_Normal, 180.0, constant.TrdSecMarket_US)
	if err != nil {
		fmt.Printf("GetMaxTrdQtys: %v\n", err)
	} else {
		fmt.Printf("%-30s %.0f\n", "Max Cash Buy", maxQtys.MaxCashBuy)
		fmt.Printf("%-30s %.0f\n", "Max Margin Buy", maxQtys.MaxCashAndMarginBuy)
		fmt.Printf("%-30s %.0f\n", "Max Position Sell", maxQtys.MaxPositionSell)
		fmt.Printf("%-30s %.0f\n", "Max Short Sell", maxQtys.MaxSellShort)
	}

	fmt.Println()
	fmt.Println("=== Analysis Complete ===")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(maxQtys)
}

func ptrInt32(v int32) *int32  { return &v }
func ptrStr(v string) *string { return &v }
