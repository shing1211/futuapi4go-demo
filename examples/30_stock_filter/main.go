package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	results, err := client.StockFilter(context.Background(), mc.Client, constant.Market_US, 0, 10)
	if err != nil {
		log.Fatalf("StockFilter failed: %v", err)
	}
	for _, r := range results {
		fmt.Printf("STOCK: code=%s name=%s price=%.2f vol=%d\n",
			r.Security.GetCode(), r.Name, r.CurPrice, r.Volume)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(results)
}
