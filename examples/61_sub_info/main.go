package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	if err := client.Subscribe(context.Background(), mc.Client, constant.Market_US, "NVDA", []constant.SubType{constant.SubType_Quote}); err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	info, err := client.GetSubInfo(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("GetSubInfo failed: %v", err)
	}
	fmt.Printf("IsSub: %v  Detail: %s\n", info.IsSub, info.Security)
	for _, t := range info.SubTypes {
		fmt.Printf("  Active SubType: %d\n", t)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(info)
}
