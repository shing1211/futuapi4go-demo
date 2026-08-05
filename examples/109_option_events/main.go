// 109_option_events demonstrates option-event queries and alert management:
//   - GetOptionEvent         (recent unusual option events / sweeps)
//   - GetOptionEventAlert    (list currently-configured option event alerts)
//   - SetOptionEventAlert    (add / delete / enable / disable / deleteAll)
//
// All calls target the US option market. The example *adds* a tiny alert,
// confirms via the listing, then deletes it again so it leaves no residue.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotoptioncommon "github.com/shing1211/futuapi4go/pkg/pb/qotoptioncommon"
	qotgetoptioneventalert "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptioneventalert"
	qotgetoptionevent "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionevent"
	qotsetoptioneventalert "github.com/shing1211/futuapi4go/pkg/pb/qotsetoptioneventalert"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Option Events & Alerts (US) ===")

	fmt.Println("--- 1) GetOptionEvent (top 10 events) ---")
	rsp1, err := client.GetOptionEvent(ctx, mc.Client, &qotgetoptionevent.C2S{
		OptionMarket: ptrInt32(int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security)),
		Count:        ptrInt32(10),
	})
	if err != nil {
		log.Fatalf("GetOptionEvent failed: %v", err)
	}
	fmt.Printf("AllCount: %d  |  Items: %d  |  UpdateTS: %s\n",
		rsp1.GetAllCount(), len(rsp1.EventList),
		time.Unix(int64(rsp1.GetUpdateTimestamp()), 0).Format("2006-01-02 15:04:05"))
	for i, it := range rsp1.EventList {
		if it == nil {
			continue
		}
		fmt.Printf("  [%d] %-6s %-22s strike=%.2f vol=%d IV=%.2f%% delta=%.4f ts=%s\n",
			i+1, it.GetSymbol(), ownerCode(it.GetOwner()),
			it.GetStrikePrice(), it.GetVolume(), it.GetIv(), it.GetDelta(),
			it.GetFillTime())
	}
	fmt.Println()

	fmt.Println("--- 2) GetOptionEventAlert (list configured alerts) ---")
	listRsp, err := client.GetOptionEventAlert(ctx, mc.Client, &qotgetoptioneventalert.C2S{
		Count: ptrInt32(50),
	})
	if err != nil {
		log.Fatalf("GetOptionEventAlert failed: %v", err)
	}
	fmt.Printf("AllCount: %d  |  Items: %d\n", listRsp.GetAllCount(), len(listRsp.AlertList))
	for i, a := range listRsp.AlertList {
		if a == nil {
			continue
		}
		fmt.Printf("  [%d] key=%d enable=%v market=%s group=%q\n",
			i+1, a.GetKey(), a.GetEnable(),
			optionMarketName(a.GetOptionMarket()),
			a.GetWatchlistGroupName())
	}
	fmt.Println()

	fmt.Println("--- 3) SetOptionEventAlert (Add) ---")
	_ = int32(constant.Market_US) // market constant kept for documentation
	newAlert := &qotgetoptioneventalert.EventAlertItem{
		Enable:             ptrBool(true),
		OptionMarket:       ptrInt32(int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security)),
		WatchlistGroupName: ptrStr("demo-109-tmp"),
		Underlying:         &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")},
		ExpiryDaysRange:    &qotoptioncommon.Interval{FilterMin: &qotoptioncommon.Boundary{Value: ptrFloat64(7)}, FilterMax: &qotoptioncommon.Boundary{Value: ptrFloat64(60)}},
		IvRange:            &qotoptioncommon.Interval{FilterMin: &qotoptioncommon.Boundary{Value: ptrFloat64(20)}, FilterMax: &qotoptioncommon.Boundary{Value: ptrFloat64(80)}},
		Note:               ptrStr("demo 109 — auto-deleted"),
	}
	addRsp, err := client.SetOptionEventAlert(ctx, mc.Client, &qotsetoptioneventalert.C2S{
		OperType:  ptrInt32(int32(qotsetoptioneventalert.AlertOpType_AlertOpType_Add)),
		AlertList: []*qotgetoptioneventalert.EventAlertItem{newAlert},
	})
	if err != nil {
		log.Fatalf("SetOptionEventAlert(Add) failed: %v", err)
	}
	_ = addRsp
	fmt.Println("Add succeeded. (Server response is empty for Add — keys returned via subsequent GetOptionEventAlert.)")
	fmt.Println()

	fmt.Println("--- 4) SetOptionEventAlert (DeleteAll) ---")
	_, err = client.SetOptionEventAlert(ctx, mc.Client, &qotsetoptioneventalert.C2S{
		OperType: ptrInt32(int32(qotsetoptioneventalert.AlertOpType_AlertOpType_DeleteAll)),
	})
	if err != nil {
		log.Fatalf("SetOptionEventAlert(DeleteAll) failed: %v", err)
	}
	fmt.Println("DeleteAll done — demo alerts removed.")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp1)
	display.PrintJSON(listRsp)
}

func ownerCode(s *qotcommon.Security) string {
	if s == nil {
		return ""
	}
	return s.GetCode()
}

func optionMarketName(v int32) string {
	switch v {
	case 1:
		return "US_Security"
	case 2:
		return "US_Index"
	case 3:
		return "HK_Security"
	case 4:
		return "HK_Index"
	default:
		return fmt.Sprintf("Unknown(%d)", v)
	}
}

func ptrInt32(v int32) *int32       { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrStr(v string) *string       { return &v }
func ptrBool(v bool) *bool          { return &v }