package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go/pkg/pb/qotsubeventcontract"
	"github.com/shing1211/futuapi4go/pkg/push"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

// Stream real-time Event Contract data (order book / kline / ticker) via
// SubEventContract + push handlers. Mirrors the general-subscription flow but
// for EC instruments (market=101).
func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client
	ctx := context.Background()

	fmt.Println("=== Event Contract Subscription (orderBook/ticker) ===")

	sec := client.NewECSecurity("EC.SUPERBOWL")

	// Register push handlers for the EC update protoIDs before subscribing.
	cli.RegisterHandler(constant.ProtoID_Qot_UpdateEventContractOrderBook, func(pid uint32, body []byte) {
		ob, err := push.ParseUpdateEventContractOrderBook(body)
		if err != nil || ob == nil {
			fmt.Printf("[EC-OrderBook] parse failed: %v\n", err)
			return
		}
		for _, item := range ob.OrderBookList {
			fmt.Printf("[EC-OrderBook] %s yesBids=%d yesAsks=%d noBids=%d noAsks=%d\n",
				item.GetCode().GetCode(), len(item.GetYesBids()), len(item.GetYesAsks()),
				len(item.GetNoBids()), len(item.GetNoAsks()))
		}
	})
	cli.RegisterHandler(constant.ProtoID_Qot_UpdateEventContractKline, func(pid uint32, body []byte) {
		kl, err := push.ParseUpdateEventContractKline(body)
		if err != nil || kl == nil {
			fmt.Printf("[EC-Kline] parse failed: %v\n", err)
			return
		}
		for _, item := range kl.KlineList {
			fmt.Printf("[EC-Kline] %s name=%s klines=%d\n",
				item.GetCode().GetCode(), item.GetName(), len(item.GetKlineList()))
		}
	})
	cli.RegisterHandler(constant.ProtoID_Qot_UpdateEventContractTicker, func(pid uint32, body []byte) {
		tk, err := push.ParseUpdateEventContractTicker(body)
		if err != nil || tk == nil {
			fmt.Printf("[EC-Ticker] parse failed: %v\n", err)
			return
		}
		for _, item := range tk.TickerList {
			for _, p := range item.GetTickerList() {
				fmt.Printf("[EC-Ticker] %s time=%s yes=%.2f vol=%.0f\n",
					item.GetCode().GetCode(), p.GetTime(), p.GetYesPrice(), p.GetVolume())
			}
		}
	})

	// Subscribe to order book + ticker pushes on the EC contract.
	subReq := &qotsubeventcontract.C2S{
		SecurityList:     []*qotcommon.Security{sec},
		SubTypeList:      []int32{int32(constant.SubType_OrderBook), int32(constant.SubType_Ticker)},
		IsSubOrUnSub:     ptrBool(true),
		IsRegOrUnRegPush: ptrBool(true),
		IsFirstPush:      ptrBool(true),
	}
	if err := client.SubEventContract(ctx, cli, subReq); err != nil {
		fmt.Printf("SubEventContract: %v\n", err)
	} else {
		fmt.Printf("Subscribed EC %s to orderbook+ticker pushes\n", sec.GetCode())
	}

	// Wait briefly for push events, then unsubscribe.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("\nListening for EC pushes (5s)...")
	select {
	case <-time.After(5 * time.Second):
	case <-sig:
	}

	unsubReq := &qotsubeventcontract.C2S{
		SecurityList:     []*qotcommon.Security{sec},
		SubTypeList:      []int32{int32(constant.SubType_OrderBook), int32(constant.SubType_Ticker)},
		IsSubOrUnSub:     ptrBool(false),
		IsRegOrUnRegPush: ptrBool(true),
	}
	if err := client.SubEventContract(ctx, cli, unsubReq); err != nil {
		fmt.Printf("Unsubscribe: %v\n", err)
	} else {
		fmt.Println("Unsubscribed EC pushes.")
	}
}

func ptrBool(v bool) *bool { return &v }
