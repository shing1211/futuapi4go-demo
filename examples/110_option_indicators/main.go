// 110_option_indicators demonstrates the technical-indicator engine:
//   - GetIndicatorList     (search for indicators by keyword — MyLang or Python)
//   - RequestIndicatorCalc (compute an indicator over a security's K-line)
//
// Workflow:
//   1. Search the MyLang indicator catalog for "MACD".
//   2. Take the resulting short-name, build an IndicatorCalcData request
//      (security + K-line + indicator inputs).
//   3. Print the computed series.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotgetindicatorlist "github.com/shing1211/futuapi4go/pkg/pb/qotgetindicatorlist"
	qotrequestindicatorcalc "github.com/shing1211/futuapi4go/pkg/pb/qotrequestindicatorcalc"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Option Indicators ===")

	fmt.Println("--- 1) GetIndicatorList (search MyLang for 'MACD') ---")
	listRsp, err := client.GetIndicatorList(ctx, mc.Client, &qotgetindicatorlist.C2S{
		SearchKey: ptrStr("MACD"),
		LangType:  ptrInt32(int32(qotcommon.IndicatorLangType_IndicatorLangType_MyLang)),
		SearchMode: ptrInt32(0),
	})
	if err != nil {
		log.Fatalf("GetIndicatorList failed: %v", err)
	}
	fmt.Printf("Indicators matched: %d\n", len(listRsp.IndicatorList))
	shortName := ""
	for i, ind := range listRsp.IndicatorList {
		if ind == nil {
			continue
		}
		info := ind.GetMyLang()
		name := ""
		if info != nil {
			name = info.GetShortName()
		}
		fmt.Printf("  [%d] MyLang=%s\n", i+1, name)
		if shortName == "" {
			shortName = name
		}
	}
	fmt.Println()

	fmt.Println("--- 2) RequestIndicatorCalc (compute first indicator on NVDA) ---")
	market := int32(constant.Market_US)
	calcRsp, err := client.RequestIndicatorCalc(ctx, mc.Client, &qotrequestindicatorcalc.C2S{
		ShortName: ptrStr(shortName),
		LangType:  langPtr(qotcommon.IndicatorLangType_IndicatorLangType_MyLang),
		Data: &qotrequestindicatorcalc.IndicatorCalcData{
			Security: &qotcommon.Security{Market: &market, Code: ptrStr("NVDA")},
			KlType:   klPtr(qotcommon.KLType_KLType_Day),
			KLine: []*qotcommon.KLine{
				{
					Time:       ptrStr("2024-12-01"),
					OpenPrice:  ptrFloat64(140.0),
					HighPrice:  ptrFloat64(145.0),
					LowPrice:   ptrFloat64(139.0),
					ClosePrice: ptrFloat64(144.5),
					Volume:     ptrInt64(100000),
				},
				{
					Time:       ptrStr("2024-12-02"),
					OpenPrice:  ptrFloat64(144.0),
					HighPrice:  ptrFloat64(148.0),
					LowPrice:   ptrFloat64(143.0),
					ClosePrice: ptrFloat64(147.0),
					Volume:     ptrInt64(120000),
				},
				{
					Time:       ptrStr("2024-12-03"),
					OpenPrice:  ptrFloat64(147.0),
					HighPrice:  ptrFloat64(149.5),
					LowPrice:   ptrFloat64(146.0),
					ClosePrice: ptrFloat64(149.0),
					Volume:     ptrInt64(110000),
				},
			},
		},
		Num:    ptrInt32(3),
		Inputs: nil,
	})
	if err != nil {
		log.Printf("RequestIndicatorCalc returned error (likely needs richer K-line history): %v\n", err)
	} else {
		fmt.Printf("Calc response received. Inspect JSON for computed series.\n")
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(listRsp)
	if calcRsp != nil {
		display.PrintJSON(calcRsp)
	}
}

func ptrInt32(v int32) *int32    { return &v }
func ptrInt64(v int64) *int64    { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrStr(v string) *string    { return &v }
func langPtr(v qotcommon.IndicatorLangType) *qotcommon.IndicatorLangType {
	return &v
}
func klPtr(v qotcommon.KLType) *qotcommon.KLType {
	return &v
}