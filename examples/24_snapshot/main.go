package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	sec1 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")}
	sec2 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("AAPL")}
	sec3 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("TSLA")}

	snapshots, err := client.GetSecuritySnapshot(context.Background(), mc.Client, []*qotcommon.Security{sec1, sec2, sec3})
	if err != nil {
		log.Fatalf("GetSecuritySnapshot failed: %v", err)
	}
	for _, s := range snapshots {
		fmt.Printf("SNAP: %s %s price=%.2f open=%.2f high=%.2f low=%.2f vol=%d\n",
			s.Security.GetCode(), s.Name, s.CurPrice, s.OpenPrice, s.HighPrice, s.LowPrice, s.Volume)

		// Extended equity data
		if s.EquityExData != nil {
			fmt.Printf("  Equity:   issuedShares=%d marketVal=%.2f netAsset=%.2f netProfit=%.2f eps=%.4f\n",
				s.EquityExData.GetIssuedShares(), s.EquityExData.GetIssuedMarketVal(),
				s.EquityExData.GetNetAsset(), s.EquityExData.GetNetProfit(),
				s.EquityExData.GetEarningsPershare())
		}
		// Extended warrant data
		if s.WarrantExData != nil {
			fmt.Printf("  Warrant:  conversionRate=%.2f strike=%.2f maturity=%s impliedVol=%.4f delta=%.4f premium=%.4f\n",
				s.WarrantExData.GetConversionRate(), s.WarrantExData.GetStrikePrice(),
				s.WarrantExData.GetMaturityTime(), s.WarrantExData.GetImpliedVolatility(),
				s.WarrantExData.GetDelta(), s.WarrantExData.GetPremium())
		}
		// Extended option data
		if s.OptionExData != nil {
			fmt.Printf("  Option:   strike=%.2f contractSize=%d openInterest=%d impliedVol=%.4f delta=%.4f gamma=%.4f vega=%.4f\n",
				s.OptionExData.GetStrikePrice(), s.OptionExData.GetContractSize(),
				s.OptionExData.GetOpenInterest(), s.OptionExData.GetImpliedVolatility(),
				s.OptionExData.GetDelta(), s.OptionExData.GetGamma(), s.OptionExData.GetVega())
		}
		// Index data
		if s.IndexExData != nil {
			fmt.Printf("  Index:    raise=%d fall=%d equal=%d\n",
				s.IndexExData.GetRaiseCount(), s.IndexExData.GetFallCount(), s.IndexExData.GetEqualCount())
		}
		// Plate data
		if s.PlateExData != nil {
			fmt.Printf("  Plate:    raise=%d fall=%d equal=%d\n",
				s.PlateExData.GetRaiseCount(), s.PlateExData.GetFallCount(), s.PlateExData.GetEqualCount())
		}
		// Futures data
		if s.FutureExData != nil {
			fmt.Printf("  Future:   lastSettle=%.2f position=%d positionChange=%d\n",
				s.FutureExData.GetLastSettlePrice(), s.FutureExData.GetPosition(),
				s.FutureExData.GetPositionChange())
		}
		// Trust data
		if s.TrustExData != nil {
			fmt.Printf("  Trust:    dividendYield=%.4f aum=%.2f nav=%.2f premium=%.4f\n",
				s.TrustExData.GetDividendYield(), s.TrustExData.GetAum(),
				s.TrustExData.GetNetAssetValue(), s.TrustExData.GetPremium())
		}
	}
}

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string { return &v }
