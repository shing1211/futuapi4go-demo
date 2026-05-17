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

	resp, err := client.GetCapitalFlow(context.Background(), mc.Client, constant.Market_HK, "00700")
	if err != nil {
		log.Fatalf("GetCapitalFlow failed: %v", err)
	}
	for _, f := range resp.Items {
		fmt.Printf("FLOW: time=%s in=%.2f main=%.2f\n",
			f.Time, f.InFlow, f.MainInFlow)
	}
	fmt.Printf("LastValidTime=%s LastValidTimestamp=%.0f\n", resp.LastValidTime, resp.LastValidTimestamp)
}
