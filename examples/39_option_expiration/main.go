package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	exps, err := client.GetOptionExpirationDate(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetOptionExpirationDate failed: %v", err)
	}
	for _, e := range exps {
		fmt.Printf("EXPIRATION: date=%s days=%d desc=%s\n",
			e.Date, e.Days, e.Desc)
	}
}
