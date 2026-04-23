package main

import (
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

	flows, err := client.GetCapitalFlow(cli, int32(constant.Market_US), "NVDA")
	if err != nil {
		log.Fatalf("GetCapitalFlow failed: %v", err)
	}
	for _, f := range flows {
		fmt.Printf("FLOW: time=%s in=%.2f main=%.2f\n",
			f.Time, f.InFlow, f.MainInFlow)
	}
}
