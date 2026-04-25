package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	rehabs, err := client.RequestRehab(context.Background(), cli, int32(constant.Market_US), "NVDA")
	if err != nil {
		log.Fatalf("RequestRehab failed: %v", err)
	}
	for _, r := range rehabs {
		fmt.Printf("REHAB: time=%s fwdA=%.4f fwdB=%.4f bwdA=%.4f bwdB=%.4f\n",
			r.Time, r.FwdFactorA, r.FwdFactorB, r.BwdFactorA, r.BwdFactorB)
	}
}
