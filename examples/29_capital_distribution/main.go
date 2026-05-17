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

	dist, err := client.GetCapitalDistribution(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetCapitalDistribution failed: %v", err)
	}
	fmt.Printf("Main: %.2f  Large: %.2f  Medium: %.2f  Small: %.2f\n",
		dist.MainInflow, dist.BigInflow, dist.MidInflow, dist.SmallInflow)

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(dist)
}
