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

	exps, err := client.GetOptionExpirationDate(context.Background(), cli, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetOptionExpirationDate failed: %v", err)
	}
	for _, e := range exps {
		fmt.Printf("EXPIRATION: date=%s days=%d desc=%s\n",
			e.Date, e.Days, e.Desc)
	}
}
