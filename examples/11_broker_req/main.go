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

	if err := client.Subscribe(context.Background(), cli, int32(constant.Market_US), "NVDA", []constant.SubType{constant.SubType_Broker}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	bids, asks, err := client.GetBroker(context.Background(), cli, int32(constant.Market_US), "NVDA", 10)
	if err != nil {
		log.Fatalf("GetBroker failed: %v", err)
	}
	for _, b := range bids {
		fmt.Printf("BID BROKER: name=%s pos=%d vol=%d\n", b.Name, b.Pos, b.Volume)
	}
	for _, a := range asks {
		fmt.Printf("ASK BROKER: name=%s pos=%d vol=%d\n", a.Name, a.Pos, a.Volume)
	}
}
