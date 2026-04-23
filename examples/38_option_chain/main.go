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

	// Get US equity option chain for NVDA
	chains, err := client.GetOptionChain(cli,
		int32(constant.Market_US), "NVDA",
		1, // indexOptionType: 1=US Equity
		0, // optType: 0=All
		0, // condition: 0=All
		"",  // beginTime: no filter
		"",  // endTime: no filter
	)
	if err != nil {
		log.Fatalf("GetOptionChain failed: %v", err)
	}
	for _, chain := range chains {
		fmt.Printf("OPTION CHAIN: strikeTime=%s\n", chain.StrikeTime)
		for _, opt := range chain.Option {
			if opt.Call != nil {
				fmt.Printf("  CALL: code=%s name=%s\n",
					opt.Call.GetBasic().GetSecurity().GetCode(), opt.Call.GetBasic().GetName())
			}
			if opt.Put != nil {
				fmt.Printf("  PUT: code=%s name=%s\n",
					opt.Put.GetBasic().GetSecurity().GetCode(), opt.Put.GetBasic().GetName())
			}
		}
	}
}
