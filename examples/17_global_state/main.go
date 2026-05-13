package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	state, err := client.GetGlobalState(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("GetGlobalState failed: %v", err)
	}
	fmt.Printf("ServerVer: %d  BuildNo: %d\n", state.ServerVer, state.ServerBuildNo)
	fmt.Printf("QotLogined: %v  TrdLogined: %v\n", state.QotLogined, state.TrdLogined)
	fmt.Printf("Market HK=%d US=%d SH=%d SZ=%d\n",
		state.MarketHK, state.MarketUS, state.MarketSH, state.MarketSZ)
}
