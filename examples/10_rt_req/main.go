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

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_RT}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	rt, err := client.GetRT(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetRT failed: %v", err)
	}
	for _, r := range rt {
		fmt.Printf("RT: time=%s minute=%d isBlank=%v price=%.2f lastClose=%.2f avgPrice=%.2f vol=%d turnover=%.2f ts=%.0f\n",
			r.Time, r.Minute, r.IsBlank, r.Price, r.LastClose, r.AvgPrice, r.Volume, r.Turnover, r.Timestamp)
	}

	fmt.Println("\n── RT Data (JSON) ────────────────────────")
	display.PrintJSON(rt)
}
