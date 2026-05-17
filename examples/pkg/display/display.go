package display

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go/pkg/sys"
)

type AllInfo struct {
	Connection  *connect.ConnectionInfo
	GlobalState *sys.GetGlobalStateResponse
	UserInfo    *sys.GetUserInfoResponse
	DelayStats  *sys.GetDelayStatisticsResponse
	UsedQuota   *sys.GetUsedQuotaResponse
	Duration    time.Duration
}

func marketState(s int32) string {
	switch s {
	case 0:
		return "N/A"
	case 1:
		return "PreMarketBegin"
	case 2:
		return "Morning"
	case 3:
		return "Rest"
	case 4:
		return "Afternoon"
	case 5:
		return "Close"
	case 6:
		return "AfterHours"
	case 7:
		return "PreMarketEnd"
	case 8:
		return "PreMarketMarket"
	case 9:
		return "PreMarketRest"
	case 10:
		return "PreMarketAfter"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

func check(v bool) string {
	if v {
		return "✓"
	}
	return "✗"
}

func protoName(id int32) string {
	names := map[int32]string{
		1002: "GetGlobalState",
		1005: "GetUserInfo",
		1007: "GetDelayStatistics",
		1010: "UsedQuota",
		3001: "GetQuote",
		3003: "GetSubInfo",
		3004: "GetKL",
		3005: "GetRT",
		3010: "GetTicker",
		3012: "GetOrderBook",
		3014: "GetBroker",
	}
	if n, ok := names[id]; ok {
		return n
	}
	return fmt.Sprintf("ProtoID_%d", id)
}

func CollectAllInfo(ctx context.Context, mc *connect.ManagedConnection) *AllInfo {
	info := &AllInfo{
		Connection: mc.Info,
	}

	cl := mc.Client
	if cl == nil {
		return info
	}
	inner := cl.Inner()

	t0 := time.Now()

	if gs, err := sys.GetGlobalState(ctx, inner); err == nil {
		info.GlobalState = gs
		if mc.Info == nil {
			mc.Info = &connect.ConnectionInfo{}
		}
		mc.Info.ConnID = gs.ConnID
		mc.Info.ServerVer = fmt.Sprintf("%d", gs.ServerVer)
		mc.Info.ServerBuildNo = gs.ServerBuildNo
		mc.Info.ServerTime = gs.Time
		mc.Info.QuoteLogin = gs.QotLogined
		mc.Info.TradeLogin = gs.TrdLogined
		mc.Info.QotSvrIpAddr = gs.QotSvrIpAddr
		mc.Info.TrdSvrIpAddr = gs.TrdSvrIpAddr
		info.Connection = mc.Info
	}

	if ui, err := sys.GetUserInfo(ctx, inner, nil); err == nil {
		info.UserInfo = ui
	}

	if q, err := sys.GetUsedQuota(ctx, inner); err == nil {
		info.UsedQuota = q
	}

	if ds, err := sys.GetDelayStatistics(ctx, inner, &sys.GetDelayStatisticsRequest{
		TypeList: []int32{1, 2, 3},
	}); err == nil {
		info.DelayStats = ds
	}

	info.Duration = time.Since(t0)
	return info
}

func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func PrintHeader() {
	fmt.Println("═══════════════════════════════════════════")
	fmt.Println("       FutuOpenD Connection Report")
	fmt.Println("═══════════════════════════════════════════")
	fmt.Printf("Timestamp: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
}

func PrintConnection(info *connect.ConnectionInfo) {
	rsaFlag := "plain"
	if info.RSAUsed {
		rsaFlag = "RSA"
	}
	fmt.Println("── Connection ─────────────────────────────")
	fmt.Printf("  Host:       %s:%d (%s)\n", info.Host, info.Port, rsaFlag)
	if info.TCPMs > 0 {
		fmt.Printf("  TCP Probe:  %.2f ms\n", info.TCPMs)
	}
	fmt.Println()
}

func PrintServer(gs *sys.GetGlobalStateResponse) {
	if gs == nil {
		fmt.Println("── Server ─────────────────────────────────")
		fmt.Println("  (failed to retrieve)")
		fmt.Println()
		return
	}
	fmt.Println("── Server ─────────────────────────────────")
	fmt.Printf("  Version:    %d\n", gs.ServerVer)
	fmt.Printf("  Build No:   %d\n", gs.ServerBuildNo)
	fmt.Printf("  ConnID:     %d\n", gs.ConnID)
	fmt.Printf("  Time:       %s\n", time.Unix(gs.Time, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  LocalTime:  %s\n", time.Unix(int64(gs.LocalTime), 0).Format("2006-01-02 15:04:05"))
	if gs.QotSvrIpAddr != "" {
		fmt.Printf("  Quote Svr:  %s\n", gs.QotSvrIpAddr)
	}
	if gs.TrdSvrIpAddr != "" {
		fmt.Printf("  Trade Svr:  %s\n", gs.TrdSvrIpAddr)
	}
	fmt.Println()
}

func PrintLogin(gs *sys.GetGlobalStateResponse) {
	if gs == nil {
		return
	}
	fmt.Println("── Login ──────────────────────────────────")
	fmt.Printf("  Quote:      %s\n", check(gs.QotLogined))
	fmt.Printf("  Trade:      %s\n", check(gs.TrdLogined))
	fmt.Println()
}

func PrintMarkets(gs *sys.GetGlobalStateResponse) {
	if gs == nil {
		return
	}
	fmt.Println("── Markets ────────────────────────────────")
	fmt.Printf("  HK:         %s\n", marketState(gs.MarketHK))
	fmt.Printf("  US:         %s\n", marketState(gs.MarketUS))
	fmt.Printf("  SH:         %s\n", marketState(gs.MarketSH))
	fmt.Printf("  SZ:         %s\n", marketState(gs.MarketSZ))
	fmt.Printf("  HK Future:  %s\n", marketState(gs.MarketHKFuture))
	fmt.Printf("  US Future:  %s\n", marketState(gs.MarketUSFuture))
	fmt.Printf("  SG Future:  %s\n", marketState(gs.MarketSGFuture))
	fmt.Printf("  JP Future:  %s\n", marketState(gs.MarketJPFuture))
	if gs.ProgramStatus != nil {
		fmt.Printf("  Program:    type=%d, desc=%s\n",
			gs.ProgramStatus.GetType(),
			gs.ProgramStatus.GetStrExtDesc(),
		)
	}
	fmt.Println()
}

func PrintUser(ui *sys.GetUserInfoResponse) {
	if ui == nil {
		return
	}
	fmt.Println("── User ───────────────────────────────────")
	fmt.Printf("  UserID:     %d\n", ui.UserID)
	if ui.NickName != "" {
		fmt.Printf("  Nickname:   %s\n", ui.NickName)
	}
	if ui.ApiLevel != "" {
		fmt.Printf("  API Level:  %s\n", ui.ApiLevel)
	}
	if ui.IsAppNNOrMM {
		fmt.Println("  Client:     NN")
	} else {
		fmt.Println("  Client:     MM")
	}
	fmt.Println("  Quote Rights:")
	fmt.Printf("    HK:       %d\n", ui.HkQotRight)
	fmt.Printf("    US:       %d\n", ui.UsQotRight)
	fmt.Printf("    CN:       %d\n", ui.CnQotRight)
	fmt.Printf("    SH:       %d\n", ui.ShQotRight)
	fmt.Printf("    SZ:       %d\n", ui.SzQotRight)
	fmt.Printf("    HK Opt:   %d\n", ui.HkOptionQotRight)
	fmt.Printf("    US Opt:   %v\n", ui.HasUSOptionQotRight)
	fmt.Printf("    US OptR:  %d\n", ui.UsOptionQotRight)
	fmt.Printf("    HK Fut:   %d\n", ui.HkFutureQotRight)
	fmt.Printf("    US Fut:   %d\n", ui.UsFutureQotRight)
	fmt.Printf("    US Idx:   %d\n", ui.UsIndexQotRight)
	fmt.Printf("    US OTC:   %d\n", ui.UsOtcQotRight)
	fmt.Printf("    US CME:   %d\n", ui.UsCMEFutureQotRight)
	fmt.Printf("    US CBOT:  %d\n", ui.UsCBOTFutureQotRight)
	fmt.Printf("    US NYMEX: %d\n", ui.UsNYMEXFutureQotRight)
	fmt.Printf("    US COMEX: %d\n", ui.UsCOMEXFutureQotRight)
	fmt.Printf("    US CBOE:  %d\n", ui.UsCBOEFutureQotRight)
	fmt.Printf("    SG Fut:   %d\n", ui.SgFutureQotRight)
	fmt.Printf("    JP Fut:   %d\n", ui.JpFutureQotRight)
	if ui.UpdateType > 0 {
		fmt.Printf("  Update:     type=%d\n", ui.UpdateType)
	}
	fmt.Println()
}

func PrintQuota(q *sys.GetUsedQuotaResponse, ui *sys.GetUserInfoResponse) {
	fmt.Println("── Quota ──────────────────────────────────")
	if q != nil {
		totalSub := int32(0)
		totalKL := int32(0)
		if ui != nil {
			totalSub = ui.SubQuota
			totalKL = ui.HistoryKLQuota
		}
		fmt.Printf("  Subscriptions:    %d / %d\n", q.UsedSubQuota, totalSub)
		fmt.Printf("  History K-Line:   %d / %d\n", q.UsedKLineQuota, totalKL)
	} else if ui != nil {
		fmt.Printf("  Sub Quota:   %d\n", ui.SubQuota)
		fmt.Printf("  KL Quota:    %d\n", ui.HistoryKLQuota)
	} else {
		fmt.Println("  (not available)")
	}
	fmt.Println()
}

func PrintDelayStats(ds *sys.GetDelayStatisticsResponse) {
	if ds == nil {
		return
	}

	if len(ds.ReqReplyStatisticsList) > 0 {
		fmt.Println("── Delay: Req-Reply ───────────────────────")
		fmt.Printf("  %-20s %6s %9s %9s %9s %5s\n",
			"Proto", "Count", "TotalAvg", "OpenDAvg", "NetDelay", "Local")
		for _, s := range ds.ReqReplyStatisticsList {
			local := "✗"
			if s.IsLocalReply {
				local = "✓"
			}
			fmt.Printf("  %-20s %6d %8.2fms %8.2fms %8.2fms  %s\n",
				protoName(s.ProtoID),
				s.Count,
				float64(s.TotalCostAvg),
				float64(s.OpenDCostAvg),
				float64(s.NetDelayAvg),
				local,
			)
		}
		fmt.Println()
	}

	if len(ds.QotPushStatisticsList) > 0 {
		fmt.Println("── Delay: Quote Push ─────────────────────")
		for _, s := range ds.QotPushStatisticsList {
			fmt.Printf("  Type=%d, count=%d, avg=%.2fms\n",
				s.QotPushType, s.Count, float32(s.DelayAvg))
		}
		fmt.Println()
	}

	if len(ds.PlaceOrderStatisticsList) > 0 {
		fmt.Println("── Delay: Place Order ────────────────────")
		for _, s := range ds.PlaceOrderStatisticsList {
			fmt.Printf("  Order=%s total=%.2fms opend=%.2fms net=%.2fms update=%.2fms\n",
				s.OrderID, s.TotalCost, s.OpenDCost, s.NetDelay, s.UpdateCost)
		}
		fmt.Println()
	}
}

func PrintDuration(d time.Duration) {
	infoCalls := 4
	fmt.Printf("  (collected %d APIs in %v)\n", infoCalls, d.Round(time.Millisecond))
	fmt.Println()
}

func PrintAll(ctx context.Context, mc *connect.ManagedConnection) *AllInfo {
	info := CollectAllInfo(ctx, mc)
	PrintHeader()
	PrintConnection(info.Connection)
	PrintServer(info.GlobalState)
	PrintLogin(info.GlobalState)
	if info.GlobalState != nil && (info.GlobalState.MarketHK != 5 || info.GlobalState.MarketUS != 5) {
		PrintMarkets(info.GlobalState)
	}
	PrintUser(info.UserInfo)
	PrintQuota(info.UsedQuota, info.UserInfo)
	PrintDelayStats(info.DelayStats)
	PrintDuration(info.Duration)

	fmt.Println("── Full Data (JSON) ──────────────────────")
	PrintJSON(info)
	return info
}
