// 34_holding_change demonstrates GetHoldingChangeList for tracking
// major shareholder holdings over time.
//
// NOTE: This API was discontinued by Futu's upstream data provider
// after 2020-12-21. The server returns an error for all requests.
// This example demonstrates graceful handling of a discontinued API.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	changes, err := client.GetHoldingChangeList(context.Background(), mc.Client,
		constant.Market_US, "NVDA",
		1, // holderCategory: 1=Mutual Fund
		"2020-01-01", "2020-12-20", // historical range before discontinuation
	)
	if err != nil {
		fe, ok := constant.AsFutuError(err)
		if ok && fe.Code == -1 {
			fmt.Println("ℹ️  GetHoldingChangeList: This API was discontinued by the upstream")
			fmt.Println("   data provider after 2020-12-21 and is no longer available.")
			fmt.Println("   See: Futu API documentation for alternative data sources.")
			return
		}
		log.Fatalf("GetHoldingChangeList failed: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("No holding change data available for the requested period.")
		return
	}

	fmt.Printf("Found %d holding change records:\n", len(changes))
	for _, c := range changes {
		fmt.Printf("  %s: %.0f shares (%.2f%%) held by %s\n",
			c.Time, c.HoldingQty, c.HoldingRatio, c.HolderName)
	}
}