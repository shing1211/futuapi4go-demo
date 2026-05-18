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

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Broker}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	result, err := client.GetBroker(context.Background(), mc.Client, constant.Market_US, "NVDA", 10)
	if err != nil {
		log.Fatalf("GetBroker failed: %v", err)
	}
	for _, b := range result.Bids {
		fmt.Printf("BID BROKER: name=%s pos=%d vol=%d\n", b.Name, b.Pos, b.Volume)
	}
	for _, a := range result.Asks {
		fmt.Printf("ASK BROKER: name=%s pos=%d vol=%d\n", a.Name, a.Pos, a.Volume)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(result)
}
