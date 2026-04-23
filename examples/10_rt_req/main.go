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

	rt, err := client.GetRT(cli, int32(constant.Market_US), "NVDA")
	if err != nil {
		log.Fatalf("GetRT failed: %v", err)
	}
	for _, r := range rt {
		fmt.Printf("RT: time=%s price=%.2f vol=%d avg=%.2f\n",
			r.Time, r.Price, r.Volume, r.AvgPrice)
	}
}
