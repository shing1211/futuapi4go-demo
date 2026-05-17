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

	rehabs, err := client.RequestRehab(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("RequestRehab failed: %v", err)
	}
	for _, r := range rehabs {
		fmt.Printf("REHAB: time=%s fwdA=%.4f fwdB=%.4f bwdA=%.4f bwdB=%.4f\n",
			r.Time, r.FwdFactorA, r.FwdFactorB, r.BwdFactorA, r.BwdFactorB)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rehabs)
}
