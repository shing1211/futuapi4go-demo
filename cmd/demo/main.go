package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetreference"
	"github.com/shing1211/futuapi4go/pkg/pb/qotstockfilter"
	"github.com/shing1211/futuapi4go/pkg/pb/trdcommon"
	"github.com/shing1211/futuapi4go/pkg/qot"
	"github.com/shing1211/futuapi4go/pkg/sys"
	"github.com/shing1211/futuapi4go/pkg/trd"
)

const defaultAddr = "127.0.0.1:11111"

const (
	MarketHK = client.Market_HK_Security
	MarketUS = client.Market_US_Security
)

func ptrStr(v string) *string   { return &v }
func ptrInt32(v int32) *int32   { return &v }
func ptrInt64(v int64) *int64  { return &v }
func ptrUint64(v uint64) *uint64 { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrBool(v bool) *bool { return &v }

func sec(market int32, code string) *qotcommon.Security {
	return &qotcommon.Security{Market: ptrInt32(market), Code: ptrStr(code)}
}

func green(v any)  { fmt.Printf("\033[32m%v\033[0m", v) }
func red(v any)    { fmt.Printf("\033[31m%v\033[0m", v) }
func yellow(v any) { fmt.Printf("\033[33m%v\033[0m", v) }
func bold(v any)   { fmt.Printf("\033[1m%v\033[0m", v) }

func section(n int, title string) {
	fmt.Printf("\n\033[1;36m═══════════════════════════════════════════════════════════════\033[0m\n")
	fmt.Printf("\033[1;36m  %d. %s\033[0m\n", n, title)
	fmt.Printf("\033[1;36m═══════════════════════════════════════════════════════════════\033[0m\n")
}

func formatVolume(v int64) string {
	if v >= 1_000_000_000 {
		return fmt.Sprintf("%.2fB", float64(v)/1_000_000_000)
	}
	if v >= 1_000_000 {
		return fmt.Sprintf("%.2fM", float64(v)/1_000_000)
	}
	if v >= 1_000 {
		return fmt.Sprintf("%.2fK", float64(v)/1_000)
	}
	return fmt.Sprintf("%d", v)
}

func formatMoney(v float64) string {
	if v >= 1_000_000_000 {
		return fmt.Sprintf("%.2fB", v/1_000_000_000)
	}
	if v >= 1_000_000 {
		return fmt.Sprintf("%.2fM", v/1_000_000)
	}
	if v >= 1_000 {
		return fmt.Sprintf("%.2fK", v/1_000)
	}
	return fmt.Sprintf("%.2f", v)
}

func trdSideName(s int32) string {
	switch s {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	case 3:
		return "SHORT"
	case 4:
		return "BUYBACK"
	default:
		return fmt.Sprintf("Side(%d)", s)
	}
}

// must checks an error. If fatal is true and error is non-nil, it exits.
// Otherwise it prints the error in red and returns false so the caller can return early.
func must(err error, fatal bool) bool {
	if err != nil {
		if fatal {
			red(fmt.Sprintf("  FATAL: %v\n", err))
			os.Exit(1)
		}
		red(fmt.Sprintf("  ERROR: %v\n", err))
		return false
	}
	return true
}

// ============================================================================
// 1. CONNECTION & SYSTEM
// ============================================================================

func demoConnection(cli *client.Client) {
	section(1, "Connection & System")

	state, err := sys.GetGlobalState(cli.Inner())
	if !must(err, false) { return }
	fmt.Printf("  ConnID:       %d\n", cli.GetConnID())
	fmt.Printf("  Server Ver:   %d\n", cli.GetServerVer())
	fmt.Printf("  Market HK:    %v\n", state.MarketHK)
	fmt.Printf("  Market US:    %v\n", state.MarketUS)
	fmt.Printf("  Market SH:    %v\n", state.MarketSH)
	fmt.Printf("  Market SZ:    %v\n", state.MarketSZ)
	fmt.Printf("  Qot Logined:  %v\n", state.QotLogined)
	fmt.Printf("  Trd Logined:  %v\n", state.TrdLogined)

	userInfo, err := sys.GetUserInfo(cli.Inner())
	if !must(err, false) { return }
	fmt.Printf("  User ID:      %d\n", userInfo.UserID)
	fmt.Printf("  Nickname:     %s\n", userInfo.NickName)
	fmt.Printf("  API Level:    %s\n", userInfo.ApiLevel)

	fmt.Println("  [GetDelayStatistics] skipped — known proto2 wire-format incompatibility with OpenD serverVer 1003")

	ms, err := qot.GetMarketState(cli.Inner(), &qot.GetMarketStateRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
	})
	if !must(err, false) { return }
	if len(ms.MarketInfoList) > 0 {
		fmt.Printf("  Market HK State: %v\n", ms.MarketInfoList[0].MarketState)
	}

	subInfo, err := qot.GetSubInfo(cli.Inner())
	if !must(err, false) { return }
	fmt.Printf("  TotalUsedQuota:  %d\n", subInfo.TotalUsedQuota)
	fmt.Printf("  RemainQuota:     %d\n", subInfo.RemainQuota)

	td, err := qot.RequestTradeDate(cli.Inner(), &qot.RequestTradeDateRequest{
		Market:    int32(qotcommon.QotMarket_QotMarket_HK_Security),
		BeginTime: "2026-01-01",
		EndTime:   "2026-12-31",
	})
	if !must(err, false) {
		return
	}
	fmt.Printf("  Trade dates (2026 HK): %d days\n", len(td.TradeDateList))
}

// ============================================================================
// 2. MARKET DATA
// ============================================================================

func demoMarketData(cli *client.Client) {
	section(2, "Market Data")

	fmt.Println("  [Subscribe] Subscribing to US.AAPL...")
	_, subErr := qot.Subscribe(cli.Inner(), &qot.SubscribeRequest{
		SecurityList: []*qotcommon.Security{sec(MarketUS, "AAPL")},
		SubTypeList:  []qot.SubType{qot.SubType_Basic, qot.SubType_KL_Day, qot.SubType_OrderBook, qot.SubType_Ticker, qot.SubType_RT},
		IsSubOrUnSub: true,
	})
	if subErr != nil {
		yellow(fmt.Sprintf("  Subscription warning: %v (some APIs may fail)\n", subErr))
	} else {
		yellow("  Subscribed successfully.\n")
	}

	fmt.Println("  [Subscribe] Subscribing to HK.00700 (Basic, KL, OrderBook, Ticker, RT)...")
	_, subErr2 := qot.Subscribe(cli.Inner(), &qot.SubscribeRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
		SubTypeList:  []qot.SubType{qot.SubType_Basic, qot.SubType_KL_Day, qot.SubType_OrderBook, qot.SubType_Ticker, qot.SubType_RT},
		IsSubOrUnSub: true,
	})
	if subErr2 != nil {
		yellow(fmt.Sprintf("  Subscription warning: %v\n", subErr2))
	} else {
		yellow("  HK subscription OK.\n")
	}

	fmt.Println("\n  [GetBasicQot] US market quote:")
	quotes, err := qot.GetBasicQot(context.Background(), cli.Inner(), []*qotcommon.Security{
		sec(MarketUS, "AAPL"),
	})
	if !must(err, false) { return }
	for _, q := range quotes {
		if q == nil {
			continue
		}
		code := q.Security.GetCode()
		name := q.Name
		if name == "" {
			name = code
		}
		chgRate := 0.0
		if q.LastClosePrice != 0 {
			chgRate = (q.CurPrice - q.LastClosePrice) / q.LastClosePrice * 100
		}
		fmt.Printf("  %-10s %-20s price=%.2f chg=%.2f%% vol=%s\n",
			code, name, q.CurPrice, chgRate, formatVolume(q.Volume))
	}

	fmt.Println("\n  [GetKL] Daily K-lines for Tencent (00700):")
	klResp, err := qot.GetKL(cli.Inner(), &qot.GetKLRequest{
		Security:  sec(MarketHK, "00700"),
		RehabType: int32(qotcommon.RehabType_RehabType_None),
		KLType:    int32(qotcommon.KLType_KLType_Day),
		ReqNum:    5,
	})
	if !must(err, false) { return }
	for _, kl := range klResp.KLList {
		fmt.Printf("  %s  O=%.2f H=%.2f L=%.2f C=%.2f  vol=%s\n",
			kl.Time, kl.OpenPrice, kl.HighPrice, kl.LowPrice, kl.ClosePrice,
			formatVolume(kl.Volume))
	}

	fmt.Println("\n  [GetOrderBook] Top 5 levels for Tencent:")
	obResp, err := qot.GetOrderBook(cli.Inner(), &qot.GetOrderBookRequest{
		Security: sec(MarketHK, "00700"),
		Num:      5,
	})
	if !must(err, false) { return }
	for i, ask := range obResp.OrderBookAskList {
		fmt.Printf("  A%02d %.2f × %s\n", i+1, ask.Price, formatVolume(ask.Volume))
	}
	fmt.Println("  ──────────────────────")
	for i, bid := range obResp.OrderBookBidList {
		fmt.Printf("  B%02d %.2f × %s\n", i+1, bid.Price, formatVolume(bid.Volume))
	}

	fmt.Println("\n  [GetTicker] Recent ticks for Tencent:")
	tickerResp, err := qot.GetTicker(cli.Inner(), &qot.GetTickerRequest{
		Security: sec(MarketHK, "00700"),
		Num:      5,
	})
	if !must(err, false) { return }
	for _, tk := range tickerResp.TickerList {
		side := "BUY"
		if tk.Dir == 2 {
			side = "SELL"
		}
		fmt.Printf("  %s  price=%.2f vol=%s [%s]\n",
			tk.Time, tk.Price, formatVolume(tk.Volume), side)
	}

	fmt.Println("\n  [GetRT] Intraday time-share for Tencent (last 10 points):")
	rtResp, err := qot.GetRT(cli.Inner(), &qot.GetRTRequest{
		Security: sec(MarketHK, "00700"),
	})
	if !must(err, false) { return }
	rtList := rtResp.RTList
	if len(rtList) > 10 {
		rtList = rtList[len(rtList)-10:]
	}
	for _, rt := range rtList {
		fmt.Printf("  %s  avg=%.2f vol=%s\n",
			rt.Time, rt.AvgPrice, formatVolume(rt.Volume))
	}

	fmt.Println("\n  [GetSecuritySnapshot] Full snapshot for Tencent:")
	snap, err := qot.GetSecuritySnapshot(cli.Inner(), &qot.GetSecuritySnapshotRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
	})
	if !must(err, false) { return }
	if len(snap.SnapshotList) > 0 {
		basic := snap.SnapshotList[0].GetBasic()
		if basic != nil {
			fmt.Printf("  Price:      %.2f\n", basic.GetCurPrice())
			fmt.Printf("  52W High:   %.2f  Low: %.2f\n", basic.GetHighest52WeeksPrice(), basic.GetLowest52WeeksPrice())
			fmt.Printf("  High today: %.2f  Low today: %.2f\n", basic.GetHighPrice(), basic.GetLowPrice())
			fmt.Printf("  Turnover:   %s  VolRatio: %.2f\n", formatVolume(int64(basic.GetTurnover())), basic.GetVolumeRatio())
		}
		if snap.SnapshotList[0].GetEquityExData() != nil {
			d := snap.SnapshotList[0].GetEquityExData()
			fmt.Printf("  PE TTM:     %.2f  PE Annual: %.2f\n", d.GetPeTTMRate(), d.GetPeRate())
			fmt.Printf("  PB:         %.2f\n", d.GetPbRate())
			fmt.Printf("  Market Cap: %s\n", formatMoney(d.GetIssuedMarketVal()))
		}
	}
}

// ============================================================================
// 3. MARKET ANALYSIS
// ============================================================================

func demoMarketAnalysis(cli *client.Client) {
	section(3, "Market Analysis")

	fmt.Println("\n  [GetPlateSet] HK industry plates:")
	plateSet, err := qot.GetPlateSet(cli.Inner(), &qot.GetPlateSetRequest{
		Market:       int32(qotcommon.QotMarket_QotMarket_HK_Security),
		PlateSetType: int32(qotcommon.PlateSetType_PlateSetType_Industry),
	})
	if !must(err, false) { return }
	for i, p := range plateSet.PlateSetList {
		if i >= 5 {
			fmt.Printf("  ... and %d more plates\n", len(plateSet.PlateSetList)-5)
			break
		}
		fmt.Printf("  %-10s %s\n", p.Plate.GetCode(), p.Name)
	}

	fmt.Println("\n  [GetPlateSecurity] Stocks in BK1094 (IT plate):")
	psResp, err := qot.GetPlateSecurity(cli.Inner(), &qot.GetPlateSecurityRequest{
		Plate: sec(MarketHK, "BK1094"),
	})
	if !must(err, false) { return }
	fmt.Printf("  Found %d stocks in BK1094\n", len(psResp.StaticInfoList))
	for i, s := range psResp.StaticInfoList {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(psResp.StaticInfoList)-5)
			break
		}
		basic := s.GetBasic()
		fmt.Printf("  %s  lot=%d\n", basic.GetSecurity().GetCode(), basic.GetLotSize())
	}

	fmt.Println("\n  [GetCapitalFlow] Capital flow for Tencent:")
	cfResp, err := qot.GetCapitalFlow(cli.Inner(), &qot.GetCapitalFlowRequest{
		Security:   sec(MarketHK, "00700"),
		PeriodType: 1,
	})
	if !must(err, false) { return }
	if len(cfResp.FlowItemList) > 0 {
		cf := cfResp.FlowItemList[len(cfResp.FlowItemList)-1]
		fmt.Printf("  Main:    %s\n", formatMoney(cf.MainInFlow))
		fmt.Printf("  Big:     %s\n", formatMoney(cf.BigInFlow))
		fmt.Printf("  Mid:     %s\n", formatMoney(cf.MidInFlow))
		fmt.Printf("  Small:   %s\n", formatMoney(cf.SmlInFlow))
	}

	fmt.Println("\n  [GetCapitalDistribution] Capital distribution for Tencent:")
	cdResp, err := qot.GetCapitalDistribution(cli.Inner(), &qotcommon.Security{
		Market: ptrInt32(MarketHK),
		Code:   ptrStr("00700"),
	})
	if !must(err, false) { return }
	if cdResp.CapitalDistribution != nil {
		cd := cdResp.CapitalDistribution
		fmt.Printf("  Super in:     %s\n", formatMoney(cd.CapitalInSuper))
		fmt.Printf("  Super out:    %s\n", formatMoney(cd.CapitalOutSuper))
		fmt.Printf("  Big in:       %s\n", formatMoney(cd.CapitalInBig))
		fmt.Printf("  Big out:      %s\n", formatMoney(cd.CapitalOutBig))
		fmt.Printf("  Mid in:       %s\n", formatMoney(cd.CapitalInMid))
		fmt.Printf("  Mid out:      %s\n", formatMoney(cd.CapitalOutMid))
		fmt.Printf("  Small in:     %s\n", formatMoney(cd.CapitalInSmall))
		fmt.Printf("  Small out:    %s\n", formatMoney(cd.CapitalOutSmall))
	}

	fmt.Println("\n  [GetOwnerPlate] Plates for Tencent:")
	opResp, err := qot.GetOwnerPlate(cli.Inner(), &qot.GetOwnerPlateRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
	})
	if !must(err, false) { return }
	for _, p := range opResp.OwnerPlateList {
		for _, plate := range p.GetPlateInfoList() {
			fmt.Printf("  %s [%d]\n", plate.GetName(), plate.GetPlateType())
		}
	}

	fmt.Println("\n  [GetReference] Related securities for HSI Futures:")
	refResp, err := qot.GetReference(cli.Inner(), &qot.GetReferenceRequest{
		Security:      sec(MarketHK, "HSImain"),
		ReferenceType: int32(qotgetreference.ReferenceType_ReferenceType_Future),
	})
	if !must(err, false) { return }
	for _, r := range refResp.StaticInfoList {
		basic := r.GetBasic()
		fmt.Printf("  %s  lot=%d\n", basic.GetSecurity().GetCode(), basic.GetLotSize())
	}

	fmt.Println("\n  [GetStaticInfo] Static info for AAPL and 600519:")
	staticResp, err := qot.GetStaticInfo(cli.Inner(), &qot.GetStaticInfoRequest{
		Market: int32(qotcommon.QotMarket_QotMarket_US_Security),
		SecType: int32(qotcommon.SecurityType_SecurityType_Warrant),
		SecurityList: []*qotcommon.Security{
			sec(MarketUS, "AAPL"),
			sec(MarketHK, "00700"),
		},
	})
	if !must(err, false) { return }
	for _, s := range staticResp.StaticInfoList {
		basic := s.GetBasic()
		fmt.Printf("  %s  type=%d  lot=%d  list=%s\n",
			basic.GetSecurity().GetCode(), basic.GetSecType(), basic.GetLotSize(), basic.GetListTime())
	}

	fmt.Println("\n  [GetFutureInfo] HSI Futures contract info:")
	fiResp, err := qot.GetFutureInfo(cli.Inner(), &qot.GetFutureInfoRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "HSImain")},
	})
	if !must(err, false) { return }
	for _, fi := range fiResp.FutureInfoList {
		fmt.Printf("  Name:         %s\n", fi.Name)
		fmt.Printf("  Contract:     %g %s\n", fi.ContractSize, fi.ContractSizeUnit)
		fmt.Printf("  Quote Unit:   %s\n", fi.QuoteUnit)
		fmt.Printf("  Min Var:      %s\n", fi.MinVarUnit)
	}
}

// ============================================================================
// 4. STOCK SCREENING
// ============================================================================

func demoStockFilter(cli *client.Client) {
	section(4, "Stock Screening (StockFilter)")

	fmt.Println("\n  [StockFilter] HK stocks: VolumeRatio > 1.5")
	filterResp, err := qot.StockFilter(cli.Inner(), &qot.StockFilterRequest{
		Market: int32(qotcommon.QotMarket_QotMarket_HK_Security),
		Begin:  0,
		Num:    20,
		BaseFilterList: []*qotstockfilter.BaseFilter{{
			FieldName:  ptrInt32(8),
			FilterMin:  ptrFloat64(1.5),
			FilterMax:  ptrFloat64(100.0),
			IsNoFilter: ptrBool(false),
		}},
	})
	if !must(err, false) { return }
	fmt.Printf("  Matches: %d\n", filterResp.AllCount)
	for i, d := range filterResp.DataList {
		if i >= 10 {
			fmt.Printf("  ... and %d more\n", len(filterResp.DataList)-10)
			break
		}
		name := d.Name
		if name == "" {
			name = d.Security.GetCode()
		}
		vr := 0.0
		for _, bd := range d.BaseDataList {
			if bd.GetFieldName() == 8 {
				vr = bd.GetValue()
				break
			}
		}
		fmt.Printf("  %-10s %-20s  VR=%.1f\n",
			d.Security.GetCode(), name, vr)
	}
}

// ============================================================================
// 5. OPTIONS & WARRANTS
// ============================================================================

func demoOptionsWarrants(cli *client.Client) {
	section(5, "Options & Warrants")

	fmt.Println("\n  [GetOptionExpirationDate] AAPL option expiry dates:")
	expResp, err := qot.GetOptionExpirationDate(cli.Inner(), &qot.GetOptionExpirationDateRequest{
		Owner:           sec(MarketUS, "AAPL"),
		IndexOptionType: 0,
	})
	if !must(err, false) { return }
	for i, d := range expResp.DateList {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(expResp.DateList)-5)
			break
		}
		fmt.Printf("  %s  dist=%d\n", d.StrikeTime, d.OptionExpiryDateDistance)
	}

	fmt.Println("\n  [GetOptionChain] AAPL option chain (first expiry):")
	if len(expResp.DateList) > 0 {
		ocResp, err := qot.GetOptionChain(cli.Inner(), &qot.GetOptionChainRequest{
			Owner:           sec(MarketUS, "AAPL"),
			IndexOptionType: 0,
			Type:            1,
			BeginTime:       expResp.DateList[0].StrikeTime,
			EndTime:         expResp.DateList[0].StrikeTime,
		})
		if !must(err, false) { return }
		count := 0
		for _, oc := range ocResp.OptionChain {
			for _, opt := range oc.Option {
				if opt.Call != nil {
					basic := opt.Call.GetBasic()
					fmt.Printf("  CALL %s  name=%s\n", basic.GetSecurity().GetCode(), basic.GetName())
					count++
				}
				if count >= 5 {
					break
				}
			}
			if count >= 5 {
				break
			}
		}
		fmt.Printf("  Total expiry groups: %d\n", len(ocResp.OptionChain))
	}

	fmt.Println("\n  [GetWarrant] HK warrants on Tencent (00700):")
	warResp, err := qot.GetWarrant(cli.Inner(), &qot.GetWarrantRequest{
		Begin:     0,
		Num:       10,
		SortField: 11,
		Ascend:    false,
		Owner:     sec(MarketHK, "00700"),
		TypeList:  []int32{1},
		Status:    1,
	})
	if !must(err, false) { return }
	fmt.Printf("  Total warrants found: %d\n", warResp.AllCount)
	for i, w := range warResp.WarrantDataList {
		if i >= 3 {
			break
		}
		fmt.Printf("  %s  price=%.3f  IV=%.2f%%  recovery=%.2f%%\n",
			w.Name, w.CurPrice, w.ImpliedVolatility*100, w.PriceRecoveryRatio*100)
	}
}

// ============================================================================
// 6. HISTORICAL DATA
// ============================================================================

func demoHistoricalData(cli *client.Client) {
	section(6, "Historical Data")

	fmt.Println("\n  [RequestHistoryKL] Last 30 daily bars for Tencent:")
	hkResp, err := qot.RequestHistoryKL(cli.Inner(), &qot.RequestHistoryKLRequest{
		RehabType:   int32(qotcommon.RehabType_RehabType_None),
		KlType:      int32(qotcommon.KLType_KLType_Day),
		Security:    sec(MarketHK, "00700"),
		BeginTime:   "2026-01-01",
		EndTime:     time.Now().Format("2006-01-02"),
		MaxAckKLNum: 30,
	})
	if !must(err, false) { return }
	for i, kl := range hkResp.KLList {
		if i >= 5 {
			fmt.Printf("  ... and %d more bars\n", len(hkResp.KLList)-5)
			break
		}
		fmt.Printf("  %s  C=%.2f  vol=%s\n",
			kl.Time, kl.ClosePrice, formatVolume(kl.Volume))
	}
	fmt.Printf("  Total returned: %d bars\n", len(hkResp.KLList))

	fmt.Println("\n  [RequestHistoryKL] Explicit page 1 (3 bars per page):")
	pageResp, err := qot.RequestHistoryKL(cli.Inner(), &qot.RequestHistoryKLRequest{
		RehabType:   int32(qotcommon.RehabType_RehabType_None),
		KlType:      int32(qotcommon.KLType_KLType_Day),
		Security:    sec(MarketUS, "AAPL"),
		BeginTime:   "2026-03-01",
		EndTime:     "2026-04-22",
		MaxAckKLNum: 3,
		NextReqKey:  nil,
	})
	if !must(err, false) { return }
	for _, kl := range pageResp.KLList {
		fmt.Printf("  %s  O=%.2f H=%.2f L=%.2f C=%.2f\n",
			kl.Time, kl.OpenPrice, kl.HighPrice, kl.LowPrice, kl.ClosePrice)
	}
	if len(pageResp.NextReqKey) > 0 {
		fmt.Printf("  Next page key: %s\n", pageResp.NextReqKey)
	}

	fmt.Println("\n  [RequestHistoryKL] 5-min K-lines for Tencent (March 2026):")
	hlResp, err := qot.RequestHistoryKL(cli.Inner(), &qot.RequestHistoryKLRequest{
		RehabType:   int32(qotcommon.RehabType_RehabType_None),
		KlType:      int32(qotcommon.KLType_KLType_5Min),
		Security:    sec(MarketHK, "00700"),
		BeginTime:   "2026-03-01 09:30:00",
		EndTime:     time.Now().Format("2006-01-02"),
		MaxAckKLNum: 5,
	})
	if err != nil {
		if strings.Contains(err.Error(), "未知的协议ID") || strings.Contains(err.Error(), "unknown protocol") {
			yellow("  [RequestHistoryKL] not supported by this OpenD server version\n")
		} else {
			fmt.Printf("\033[31m  ERROR: RequestHistoryKL failed: %v\033[0m\n", err)
		}
	} else {
		for _, kl := range hlResp.KLList {
			if kl.IsBlank {
				continue
			}
			fmt.Printf("  %s  C=%.2f  vol=%s\n",
				kl.Time, kl.ClosePrice, formatVolume(kl.Volume))
		}
		if len(hlResp.NextReqKey) > 0 {
			fmt.Println("  ... more data available (pagination)")
		}
	}

	fmt.Println("\n  [RequestHistoryKLQuota] API quota:")
	quotaResp, err := qot.RequestHistoryKLQuota(cli.Inner(), &qot.RequestHistoryKLQuotaRequest{
		GetDetail: false,
	})
	if !must(err, false) { return }
	fmt.Printf("  Used quota:   %d\n", quotaResp.UsedQuota)
	fmt.Printf("  Remain quota: %d\n", quotaResp.RemainQuota)

	fmt.Println("\n  [GetRehab] Adjustment factors for Tencent:")
	rehabResp, err := qot.GetRehab(cli.Inner(), &qot.GetRehabRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
	})
	if !must(err, false) { return }
	if len(rehabResp.SecurityRehabList) > 0 {
		r := rehabResp.SecurityRehabList[0]
		if len(r.GetRehabList()) > 0 {
			rehab := r.GetRehabList()[0]
			fmt.Printf("  Time:       %s\n", rehab.GetTime())
			fmt.Printf("  FwdFactor: %.6f / %.6f\n", rehab.GetFwdFactorA(), rehab.GetFwdFactorB())
			fmt.Printf("  BwdFactor: %.6f / %.6f\n", rehab.GetBwdFactorA(), rehab.GetBwdFactorB())
			fmt.Printf("  Dividend:  %.6f\n", rehab.GetDividend())
		}
	}
}

// ============================================================================
// 7. CORPORATE ACTIONS
// ============================================================================

func demoCorporateActions(cli *client.Client) {
	section(7, "Corporate Actions")

	fmt.Println("\n  [GetIpoList] Recent HK IPOs:")
	ipoResp, err := qot.GetIpoList(cli.Inner(), &qot.GetIpoListRequest{
		Market: int32(qotcommon.QotMarket_QotMarket_HK_Security),
	})
	if !must(err, false) { return }
	for i, ipo := range ipoResp.IpoList {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(ipoResp.IpoList)-5)
			break
		}
		if ipo.Basic != nil {
			fmt.Printf("  %-10s %-20s  list=%s\n",
				ipo.Basic.Security.GetCode(), ipo.Basic.Name, ipo.Basic.ListTime)
		}
	}

	fmt.Println("\n  [GetCodeChange] Code changes for Tencent:")
	ccResp, err := qot.GetCodeChange(cli.Inner(), &qot.GetCodeChangeRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
	})
	if err != nil {
		if strings.Contains(err.Error(), "未知的协议ID") || strings.Contains(err.Error(), "unknown protocol") || strings.Contains(err.Error(), "Unknown proto") {
			yellow("  [GetCodeChange] not supported by this OpenD server version\n")
		} else {
			if !must(err, false) { return }
		}
	} else if len(ccResp.CodeChangeList) == 0 {
		fmt.Println("  No recent code changes found.")
	} else {
		for _, c := range ccResp.CodeChangeList {
			fmt.Printf("  Code change: type=%d  effective=%s\n", c.Type, c.EffectiveTime)
		}
	}

	fmt.Println("\n  [GetSuspend] Suspension info for Tencent (2026):")
	suspResp, err := qot.GetSuspend(cli.Inner(), &qot.GetSuspendRequest{
		SecurityList: []*qotcommon.Security{sec(MarketHK, "00700")},
		BeginTime:    "2026-01-01",
		EndTime:      time.Now().Format("2006-01-02"),
	})
	if err != nil {
		if strings.Contains(err.Error(), "未知的协议ID") || strings.Contains(err.Error(), "unknown protocol") {
			yellow("  [GetSuspend] not supported by this OpenD server version\n")
		} else {
			fmt.Printf("\033[31m  ERROR: GetSuspend failed: %v\033[0m\n", err)
		}
	} else {
		found := false
		for _, ssl := range suspResp.SecuritySuspendList {
			for _, s := range ssl.SuspendList {
				if s != nil {
					fmt.Printf("  Suspended: %s\n", s.Time)
					found = true
				}
			}
		}
		if !found {
			fmt.Println("  Not suspended in 2026")
		}
	}

	fmt.Println("\n  [GetHoldingChangeList] Major holder changes for Tencent:")
	holdResp, err := qot.GetHoldingChangeList(cli.Inner(), &qot.GetHoldingChangeListRequest{
		Security:       sec(MarketHK, "00700"),
		HolderCategory: 1,
		BeginTime:      "2025-01-01",
		EndTime:        time.Now().Format("2006-01-02"),
	})
	if err != nil {
		if strings.Contains(err.Error(), "未知的协议ID") || strings.Contains(err.Error(), "unknown protocol") {
			yellow("  [GetHoldingChangeList] not supported by this OpenD server version\n")
		} else {
			fmt.Printf("\033[31m  ERROR: GetHoldingChangeList failed: %v\033[0m\n", err)
		}
	} else if len(holdResp.HoldingChangeList) == 0 {
		fmt.Println("  No holder changes found in this period.")
	} else {
		for i, h := range holdResp.HoldingChangeList {
			if i >= 3 {
				fmt.Printf("  ... and %d more\n", len(holdResp.HoldingChangeList)-3)
				break
			}
			fmt.Printf("  %s  holder=%s  holding=%.2f%%\n",
				h.GetTime(), h.GetHolderName(), h.GetHoldingRatio())
		}
	}
}

// ============================================================================
// 8. TRADING OPERATIONS
// ============================================================================

func demoTrading(cli *client.Client) {
	section(8, "Trading Operations")

	fmt.Println("\n  [GetAccList] Trading accounts:")
	accResp, err := trd.GetAccList(cli.Inner(), int32(trdcommon.TrdCategory_TrdCategory_Security), false)
	if !must(err, false) { return }
	var realSecAccID uint64
	var realSecTrdEnv int32
	var realSecTrdMkt int32
	for _, acc := range accResp.AccList {
		env := "SIMULATE"
		if acc.TrdEnv == 1 {
			env = "REAL"
		}
		fmt.Printf("  AccID=%d  env=%s  card=%s  firm=%d\n",
			acc.AccID, env, acc.CardNum, acc.SecurityFirm)
		if acc.TrdEnv == 1 && acc.AccType != 0 && realSecAccID == 0 {
			realSecAccID = acc.AccID
			realSecTrdEnv = acc.TrdEnv
			realSecTrdMkt = int32(qotcommon.QotMarket_QotMarket_HK_Security)
		}
	}

	if realSecAccID == 0 {
		yellow("  No real securities trading account found — skipping trading API calls.\n")
		return
	}

	trdEnv := realSecTrdEnv
	trdMkt := realSecTrdMkt

	fmt.Printf("\n  [GetFunds] Account %d funds:\n", realSecAccID)
	fundsResp, err := trd.GetFunds(cli.Inner(), &trd.GetFundsRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
	})
	if !must(err, false) { return }
	if fundsResp.Funds != nil {
		f := fundsResp.Funds
		fmt.Printf("  Currency:      %d\n", f.Currency)
		fmt.Printf("  Total Assets: %s\n", formatMoney(f.TotalAssets))
		fmt.Printf("  Cash:          %s\n", formatMoney(f.Cash))
		fmt.Printf("  Market Value:  %s\n", formatMoney(f.MarketVal))
		fmt.Printf("  Available:    %s\n", formatMoney(f.AvailableFunds))
		fmt.Printf("  BP:            %s\n", formatMoney(f.Power))
	}

	fmt.Printf("\n  [GetPositionList] Positions for AccID %d:\n", realSecAccID)
	posResp, err := trd.GetPositionList(cli.Inner(), &trd.GetPositionListRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
	})
	if !must(err, false) { return }
	fmt.Printf("  Open positions: %d\n", len(posResp.PositionList))
	for i, p := range posResp.PositionList {
		if i >= 3 {
			fmt.Printf("  ... and %d more\n", len(posResp.PositionList)-3)
			break
		}
		fmt.Printf("  %s  qty=%.0f  cost=%.2f  P&L=%.2f (%.2f%%)\n",
			p.Code, p.Qty, p.CostPrice, p.PlVal, p.PlRatio*100)
	}

	fmt.Println("\n  [GetMaxTrdQtys] Max qty for Tencent (limit order):")
	maxResp, err := trd.GetMaxTrdQtys(cli.Inner(), &trd.GetMaxTrdQtysRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
		OrderType: int32(trdcommon.OrderType_OrderType_Normal),
		Code:      "00700",
		Price:     400.0,
	})
	if !must(err, false) { return }
	if maxResp.MaxTrdQtys != nil {
		m := maxResp.MaxTrdQtys
		fmt.Printf("  CashBuy:    %.0f  Margin: %.0f  Short: %.0f  BuyBack: %.0f\n",
			m.MaxCashBuy, m.MaxCashAndMarginBuy, m.MaxSellShort, m.MaxBuyBack)
	}

	fmt.Println("\n  [GetOrderList] Active orders:")
	ordResp, err := trd.GetOrderList(cli.Inner(), &trd.GetOrderListRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
	})
	if !must(err, false) { return }
	if len(ordResp.OrderList) == 0 {
		fmt.Println("  No active orders.")
	}
	for _, o := range ordResp.OrderList {
		side := "BUY"
		if o.TrdSide == 2 {
			side = "SELL"
		}
		fmt.Printf("  %s %s %.0f@%.2f  status=%d\n",
			o.Code, side, o.Qty, o.Price, o.OrderStatus)
	}

	fmt.Println("\n  [GetOrderFillList] Recent fills:")
	fillResp, err := trd.GetOrderFillList(cli.Inner(), &trd.GetOrderFillListRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
	})
	if !must(err, false) { return }
	if len(fillResp.OrderFillList) == 0 {
		fmt.Println("  No recent fills.")
	}
	for i, f := range fillResp.OrderFillList {
		if i >= 3 {
			break
		}
		side := "BUY"
		if f.TrdSide == 2 {
			side = "SELL"
		}
		fmt.Printf("  %s %s qty=%.0f@%.2f\n",
			f.Code, side, f.Qty, f.Price)
	}

	fmt.Println("\n  [GetHistoryOrderList] Last 7 days:")
	beginTime := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	endTime := time.Now().Format("2006-01-02")
	histOrdResp, err := trd.GetHistoryOrderList(cli.Inner(), &trd.GetHistoryOrderListRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
		FilterConditions: &trdcommon.TrdFilterConditions{
			BeginTime: ptrStr(beginTime),
			EndTime:   ptrStr(endTime),
		},
	})
	if !must(err, false) { return }
	fmt.Printf("  Historical orders: %d\n", len(histOrdResp.OrderList))
	for i, o := range histOrdResp.OrderList {
		if i >= 3 {
			fmt.Printf("  ... and %d more\n", len(histOrdResp.OrderList)-3)
			break
		}
		fmt.Printf("  %s %s %s  status=%d  qty=%.0f fill=%.0f\n",
			o.GetCreateTime(), o.GetCode(), trdSideName(o.GetTrdSide()),
			o.GetOrderStatus(), o.GetQty(), o.GetFillQty())
	}

	fmt.Println("\n  [PlaceOrder] Placing DEMO paper order (no real trade):")
	placeResp, err := trd.PlaceOrder(cli.Inner(), &trd.PlaceOrderRequest{
		AccID:     realSecAccID,
		TrdMarket: trdMkt,
		TrdEnv:    trdEnv,
		TrdSide:   int32(trdcommon.TrdSide_TrdSide_Buy),
		OrderType: int32(trdcommon.OrderType_OrderType_Normal),
		Code:      "00700",
		Qty:       100,
		Price:     350.0,
	})
	if err != nil {
		fmt.Printf("  Expected error (no unlock): %v\n", err)
	} else {
		fmt.Printf("  OrderID=%d  OrderIDEx=%s\n", placeResp.OrderID, placeResp.OrderIDEx)
	}

	fmt.Println("\n  [GetFlowSummary] Fund flow summary (today):")
	clearDate := time.Now().Format("2006-01-02")
	fsResp, err := trd.GetFlowSummary(cli.Inner(), &trd.GetFlowSummaryRequest{
		Header:       &trdcommon.TrdHeader{AccID: ptrUint64(realSecAccID), TrdEnv: ptrInt32(trdEnv)},
		ClearingDate: clearDate,
	})
	if err != nil {
		fmt.Printf("  Skipped: %v\n", err)
	} else if len(fsResp.FlowSummaryList) > 0 {
		for _, fs := range fsResp.FlowSummaryList {
			fmt.Printf("  Date: %s  Amount: %s\n",
				fs.GetClearingDate(), formatMoney(fs.GetCashFlowAmount()))
		}
	}
}

// ============================================================================
// 9. USER SECURITY GROUPS & PRICE ALERTS
// ============================================================================

func demoUserGroupsAlerts(cli *client.Client) {
	section(9, "User Security Groups & Price Alerts")

	fmt.Println("\n  [GetUserSecurityGroup] User security groups:")
	grpResp, err := qot.GetUserSecurityGroup(cli.Inner(), &qot.GetUserSecurityGroupRequest{
		GroupType: 1,
	})
	if !must(err, false) { return }
	for _, g := range grpResp.GroupList {
		fmt.Printf("  Group: %s  (type=%d)\n", g.GroupName, g.GroupType)
	}

	if len(grpResp.GroupList) > 0 {
		grpName := grpResp.GroupList[0].GroupName
		fmt.Printf("\n  [GetUserSecurity] Securities in '%s':\n", grpName)
		secResp, err := qot.GetUserSecurity(cli.Inner(), grpName)
		if !must(err, false) { return }
		for i, s := range secResp.StaticInfoList {
			if i >= 5 {
				break
			}
			basic := s.GetBasic()
			fmt.Printf("  %s  name=%s\n", basic.GetSecurity().GetCode(), basic.GetName())
		}
	}

	fmt.Println("\n  [ModifyUserSecurity] Skipped in demo (requires existing group).")
	fmt.Println("  Add: secList=[HK.00700], opType=1 (add)")
	fmt.Println("  Remove: secList=[HK.00700], opType=2 (remove)")

	fmt.Println("\n  [SetPriceReminder] Setting alert for Tencent @ 400.00:")
	alertResp, err := qot.SetPriceReminder(cli.Inner(), &qot.SetPriceReminderRequest{
		Security: sec(MarketHK, "00700"),
		Op:       1,
		Type:     1,
		Value:    400.00,
		Note:     "Demo alert from futuapi4go-demo",
	})
	if err != nil {
		fmt.Printf("  Note: %v\n", err)
	} else {
		fmt.Printf("  Alert Key: %d\n", alertResp.Key)
	}

	fmt.Println("\n  [GetPriceReminder] Active alerts for Tencent:")
	prResp, err := qot.GetPriceReminder(cli.Inner(), sec(MarketHK, "00700"), MarketHK)
	if !must(err, false) { return }
	if len(prResp.PriceReminderList) == 0 {
		fmt.Println("  No active alerts.")
	}
	for i, a := range prResp.PriceReminderList {
		if i >= 5 {
			break
		}
		fmt.Printf("  Security: %s  Name: %s\n", a.Security.GetCode(), a.Name)
		for j, item := range a.ItemList {
			if j >= 3 {
				break
			}
			fmt.Printf("    Alert type=%d  value=%.2f  note=%s\n",
				item.Type, item.Value, item.Note)
		}
	}
}

// ============================================================================
// 10. REAL-TIME PUSH SUBSCRIPTIONS
// ============================================================================

func demoPushSubscriptions(cli *client.Client) {
	section(10, "Real-time Push Subscriptions")

	fmt.Println("\n  Subscribing to US.AAPL...")
	fmt.Println("  Watching for BasicQot, KL (1min), OrderBook, Ticker updates...")
	fmt.Println("  Press Ctrl+C to stop.")

	cli.RegisterHandler(client.ProtoID_Qot_UpdateBasicQot, func(protoID uint32, body []byte) {
		q, err := client.ParsePushQuote(body)
		if err != nil || q == nil {
			return
		}
		fmt.Printf("\033[35m[QUOTE]   %s  %.2f  vol=%s\033[0m\n",
			q.Code, q.CurPrice, formatVolume(q.Volume))
	})

	cli.RegisterHandler(client.ProtoID_Qot_UpdateKL, func(protoID uint32, body []byte) {
		kl, err := client.ParsePushKLine(body)
		if err != nil || kl == nil {
			return
		}
		klType := "K"
		switch kl.KLType {
		case 1:
			klType = "1m"
		case 2:
			klType = "1d"
		}
		fmt.Printf("\033[34m[KL]      %s %s  C=%.2f  vol=%s\033[0m\n",
			kl.Code, klType, kl.Close, formatVolume(kl.Volume))
	})

	cli.RegisterHandler(client.ProtoID_Qot_UpdateOrderBook, func(protoID uint32, body []byte) {
		ob, err := client.ParsePushOrderBook(body)
		if err != nil || ob == nil || len(ob.Asks) == 0 || len(ob.Bids) == 0 {
			return
		}
		fmt.Printf("\033[33m[BOOK]    %s  ask=%.2f  bid=%.2f  depth=%d\033[0m\n",
			ob.Code, ob.Asks[0].Price, ob.Bids[0].Price,
			len(ob.Asks)+len(ob.Bids))
	})

	cli.RegisterHandler(client.ProtoID_Qot_UpdateTicker, func(protoID uint32, body []byte) {
		tk, err := client.ParsePushTicker(body)
		if err != nil || tk == nil {
			return
		}
		side := "B"
		if tk.Side == 2 {
			side = "S"
		}
		fmt.Printf("\033[36m[TICK]    %s  %.2f × %s [%s]\033[0m\n",
			tk.Code, tk.Price, formatVolume(tk.Volume), side)
	})

	_, err := qot.Subscribe(cli.Inner(), &qot.SubscribeRequest{
		SecurityList: []*qotcommon.Security{
			sec(MarketUS, "AAPL"),
		},
		SubTypeList:      []qot.SubType{qot.SubType_Basic, qot.SubType_KL_1Min, qot.SubType_OrderBook, qot.SubType_Ticker},
		IsSubOrUnSub:     true,
		IsRegOrUnRegPush: true,
	})
	if !must(err, false) { return }
	fmt.Println("  Subscribed successfully.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	client.UnsubscribeAll(cli)
	fmt.Println("\n  Unsubscribed. Done.")
}

// ============================================================================
// MAIN MENU
// ============================================================================

func main() {
	addr := defaultAddr
	if a := os.Getenv("FUTU_ADDR"); a != "" {
		addr = a
	}

	bold("\n╔══════════════════════════════════════════════════════════════╗\n")
	bold("║              futuapi4go - SDK Demo  v0.1.0                   ║\n")
	bold("║     https://github.com/shing1211/futuapi4go                 ║\n")
	bold("╠══════════════════════════════════════════════════════════════╣\n")
	bold("║  Futu OpenD address: ")
	fmt.Print(addr)
	bold("                                      ║\n")
	bold("╚══════════════════════════════════════════════════════════════╝\n")

	cli := client.New(client.WithLogLevel(2))
	defer cli.Close()

	if err := cli.Connect(addr); err != nil {
		red(fmt.Sprintf("\n  Connection failed: %v\n", err))
		red("  Make sure Futu OpenD is running on " + addr + "\n")
		red("  Or start the mock simulator:\n")
		red("    go run github.com/shing1211/futuapi4go/cmd/simulator\n")
		os.Exit(1)
	}

	green("  Connected\n")
	fmt.Printf("  ConnID: %d  ServerVer: %d\n", cli.GetConnID(), cli.GetServerVer())

	for {
		fmt.Println()
		bold("╔══════════════════════════════════════════════╗\n")
		bold("║         SDK Demo Menu                        ║\n")
		bold("╠══════════════════════════════════════════════╣\n")
		fmt.Printf("║  1.  Connection & System                    ║\n")
		fmt.Printf("║  2.  Market Data (Quotes, K-Lines, Book)   ║\n")
		fmt.Printf("║  3.  Market Analysis (Plates, Capital)     ║\n")
		fmt.Printf("║  4.  Stock Screening (StockFilter)         ║\n")
		fmt.Printf("║  5.  Options & Warrants                     ║\n")
		fmt.Printf("║  6.  Historical Data (paginated K-lines)   ║\n")
		fmt.Printf("║  7.  Corporate Actions (IPO, Splits)       ║\n")
		fmt.Printf("║  8.  Trading Operations                    ║\n")
		fmt.Printf("║  9.  User Groups & Price Alerts            ║\n")
		fmt.Printf("║  10. Real-time Push Subscriptions          ║\n")
		fmt.Printf("║  0.  Run All (comprehensive demo)          ║\n")
		fmt.Printf("║  q.  Quit                                  ║\n")
		bold("╚══════════════════════════════════════════════╝\n")
		fmt.Print("\nChoice: ")

		var choice string
		if _, err := fmt.Scan(&choice); err != nil {
			break
		}
		choice = strings.TrimSpace(strings.ToLower(choice))

		switch choice {
		case "1":
			demoConnection(cli)
		case "2":
			demoMarketData(cli)
		case "3":
			demoMarketAnalysis(cli)
		case "4":
			demoStockFilter(cli)
		case "5":
			demoOptionsWarrants(cli)
		case "6":
			demoHistoricalData(cli)
		case "7":
			demoCorporateActions(cli)
		case "8":
			demoTrading(cli)
		case "9":
			demoUserGroupsAlerts(cli)
		case "10":
			demoPushSubscriptions(cli)
		case "0":
			runAll(cli)
		case "q", "quit", "exit":
			fmt.Println("Bye!")
			return
		default:
			red("Invalid choice.\n")
		}
	}
}

func runAll(cli *client.Client) {
	bold("\n███████████████████████████████████████████████████████████████████\n")
	bold("█               RUNNING ALL DEMOS                                  █\n")
	bold("███████████████████████████████████████████████████████████████████\n")
	demoConnection(cli)
	demoMarketData(cli)
	demoMarketAnalysis(cli)
	demoStockFilter(cli)
	demoOptionsWarrants(cli)
	demoHistoricalData(cli)
	demoCorporateActions(cli)
	demoTrading(cli)
	demoUserGroupsAlerts(cli)
	fmt.Println()
	bold("███████████████████████████████████████████████████████████████████\n")
	bold("█               ALL DEMOS COMPLETE                                  █\n")
	bold("███████████████████████████████████████████████████████████████████\n")
	fmt.Println("\nFor live push demo, run option 10 separately.")
}
