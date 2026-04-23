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

	if err := client.Subscribe(cli, int32(constant.Market_US), "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	resp, err := client.GetSubInfo(cli)
	if err != nil {
		log.Fatalf("GetSubInfo failed: %v", err)
	}
	fmt.Printf("IsSub: %v  Detail: %s\n", resp.IsSub, resp.Security)
	for _, t := range resp.SubTypes {
		fmt.Printf("  Active SubType: %d\n", t)
	}
}
