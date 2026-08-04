package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	predcommon "github.com/shing1211/futuapi4go/pkg/pb/common"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotkline "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractkline"
	qotorder "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractorderbook"
	qotsnap "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractsnapshot"
	qotticker "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractticker"
	qothistkl "github.com/shing1211/futuapi4go/pkg/pb/qotrequesthistoryeventcontractkl"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

// Real-time and historical market data for Event Contracts.
// Each EC contract has a YES and NO side with independent bids/asks.
func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client

	// A concrete EC contract security (market=101). Replace with an active
	// contract code discovered via example 102.
	sec := client.NewECSecurity("EC.SUPERBOWL")

	fmt.Println("=== Event Contract Market Data ===")

	// 1. Snapshot — latest YES/NO bid/ask, last price, cumulative volume.
	fmt.Println("\n-- GetEventContractSnapshot --")
	snapRsp, err := client.GetEventContractSnapshot(context.Background(), cli,
		&qotsnap.C2S{SecurityList: []*qotcommon.Security{sec}})
	if err != nil {
		fmt.Printf("  GetEventContractSnapshot: %v\n", err)
	} else {
		for _, s := range snapRsp.SnapshotList {
			fmt.Printf("  %s  status=%d price=%.2f vol=%.0f\n",
				s.GetCode().GetCode(), s.GetStatus(), s.GetPrice(), s.GetCumulativeVolume())
			fmt.Printf("    YES bid=%.2f/%.0f ask=%.2f/%.0f\n",
				s.GetYesBid(), s.GetYesBidSize(), s.GetYesAsk(), s.GetYesAskSize())
			fmt.Printf("    NO  bid=%.2f/%.0f ask=%.2f/%.0f\n",
				s.GetNoBid(), s.GetNoBidSize(), s.GetNoAsk(), s.GetNoAskSize())
		}
	}

	// 2. Order book — YES and NO bid/ask levels.
	fmt.Println("\n-- GetEventContractOrderBook (5 levels) --")
	obRsp, err := client.GetEventContractOrderBook(context.Background(), cli,
		&qotorder.C2S{Security: sec, Num: ptrInt32(5)})
	if err != nil {
		fmt.Printf("  GetEventContractOrderBook: %v\n", err)
	} else {
		for _, ob := range obRsp.OrderBookList {
			fmt.Printf("  %s\n", ob.GetCode().GetCode())
			fmt.Printf("    YES bids=%d asks=%d  NO bids=%d asks=%d\n",
				len(ob.GetYesBids()), len(ob.GetYesAsks()), len(ob.GetNoBids()), len(ob.GetNoAsks()))
		}
	}

	// 3. Ticker — recent trades with PredSide (Yes/No).
	fmt.Println("\n-- GetEventContractTicker (recent 5) --")
	tkRsp, err := client.GetEventContractTicker(context.Background(), cli,
		&qotticker.C2S{Security: sec, Count: ptrInt32(5)})
	if err != nil {
		fmt.Printf("  GetEventContractTicker: %v\n", err)
	} else {
		for _, t := range tkRsp.TickerList {
			for _, p := range t.GetTickerList() {
				side := int32(predcommon.PredSide_PredSide_Unknown)
				if p.Side != nil {
					side = int32(*p.Side)
				}
				fmt.Printf("  time=%s yes=%.2f no=%.2f vol=%.0f side=%d\n",
					p.GetTime(), p.GetYesPrice(), p.GetNoPrice(), p.GetVolume(), side)
			}
		}
	}

	// 4. Real-time Kline (contract-level YES/NO via preSide).
	fmt.Println("\n-- GetEventContractKline (Day, predSide=YES) --")
	klRsp, err := client.GetEventContractKline(context.Background(), cli,
		&qotkline.C2S{
			Security: sec,
			PreSide:  ptrPred(predcommon.PredSide_PredSide_Yes),
			Ktype:    ptrKLType(qotcommon.KLType_KLType_Day),
			MaxCount: ptrInt32(30),
		})
	if err != nil {
		fmt.Printf("  GetEventContractKline: %v\n", err)
	} else {
		for _, k := range klRsp.KlineList {
			for _, bar := range k.GetKlineList() {
				fmt.Printf("  %s O=%.2f H=%.2f L=%.2f C=%.2f V=%.0f\n",
					bar.GetTimeKey(), bar.GetOpen(), bar.GetHigh(), bar.GetLow(), bar.GetClose(), bar.GetVolume())
			}
		}
	}

	// 5. Historical K-line over a time range.
	fmt.Println("\n-- RequestHistoryEventContractKL --")
	histRsp, err := client.RequestHistoryEventContractKL(context.Background(), cli,
		&qothistkl.C2S{
			Security: sec,
			PreSide:  ptrPred(predcommon.PredSide_PredSide_Yes),
			KlType:   ptrKLType(qotcommon.KLType_KLType_Day),
			BeginTime: ptrStr("2026-01-01 00:00:00"),
			EndTime:  ptrStr("2026-07-01 00:00:00"),
		})
	if err != nil {
		fmt.Printf("  RequestHistoryEventContractKL: %v\n", err)
	} else {
		for _, k := range histRsp.KlineList {
			for _, bar := range k.GetKlineList() {
				fmt.Printf("  %s O=%.2f H=%.2f L=%.2f C=%.2f V=%.0f\n",
					bar.GetTimeKey(), bar.GetOpen(), bar.GetHigh(), bar.GetLow(), bar.GetClose(), bar.GetVolume())
			}
		}
	}

	fmt.Println("\n── Full responses (JSON) ────────────────────────")
	if snapRsp != nil {
		display.PrintJSON(snapRsp)
	}
	if obRsp != nil {
		display.PrintJSON(obRsp)
	}
}

func ptrInt32(v int32) *int32                        { return &v }
func ptrStr(v string) *string                        { return &v }
func ptrPred(v predcommon.PredSide) *predcommon.PredSide { return &v }
func ptrKLType(v qotcommon.KLType) *qotcommon.KLType     { return &v }
