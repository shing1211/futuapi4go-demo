package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	infos, err := client.GetStaticInfo(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetStaticInfo failed: %v", err)
	}
	for _, info := range infos {
		fmt.Printf("STATIC: code=%s name=%s type=%d lotSize=%d listTime=%s\n",
			info.Code, info.Name, info.Type, info.LotSize, info.ListTime)
		if info.WarrantExData != nil {
			fmt.Printf("  Warrant: type=%d\n", info.WarrantExData.GetType())
		}
		if info.OptionExData != nil {
			fmt.Printf("  Option: type=%d strike=%.2f strikeTime=%s\n",
				info.OptionExData.GetType(), info.OptionExData.GetStrikePrice(),
				info.OptionExData.GetStrikeTime())
		}
		if info.FutureExData != nil {
			fmt.Printf("  Future: lastTrade=%s isMain=%v\n",
				info.FutureExData.GetLastTradeTime(), info.FutureExData.GetIsMainContract())
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(infos)
}
