// 113_top_brokers demonstrates the top-10 buy/sell brokers (Tier 2):
//   - GetTopTenBuySellBrokers  (combined buy+sell top-10 broker rankings)
//
// The BuySellType distinguishes buy-side (0) vs sell-side (1) rows.
// Use a liquid HK stock (00700.Tencent) or US mega cap for reliable data.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	market := constant.Market(constant.Market_HK)
	code := "00700" // Tencent

	fmt.Println("=== Top Brokers (HK.00700) ===")

	fmt.Println("--- GetTopTenBuySellBrokers (latest day) ---")
	rsp, err := client.GetTopTenBuySellBrokers(ctx, mc.Client, market, code, 0 /*daysBefore*/)
	if err != nil {
		log.Panicf("GetTopTenBuySellBrokers failed: %v", err)
	}

	realTime := "archive"
	if rsp.IsRealTime {
		realTime = "real-time"
	}
	fmt.Printf("  isRealTime=%s  dataTime=%s  brokers=%d\n",
		realTime, rsp.DataTimeStr, len(rsp.BrokerList))

	for _, b := range rsp.BrokerList {
		side := "BUY "
		if b.BuySellType == 1 {
			side = "SELL"
		}
		fmt.Printf("  [%s] %-30s netVol=%-10d avgPrice=%8.2f\n",
			side, b.BrokerName, b.NetVol, b.AvgPrice)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp)
}
