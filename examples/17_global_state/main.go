package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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

	state, err := client.GetGlobalState(context.Background(), cli)
	if err != nil {
		log.Fatalf("GetGlobalState failed: %v", err)
	}
	fmt.Printf("ServerVer: %d  BuildNo: %d\n", state.ServerVer, state.ServerBuildNo)
	fmt.Printf("QotLogined: %v  TrdLogined: %v\n", state.QotLogined, state.TrdLogined)
	fmt.Printf("Market HK=%d US=%d SH=%d SZ=%d\n",
		state.MarketHK, state.MarketUS, state.MarketSH, state.MarketSZ)
}
