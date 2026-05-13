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

	funds, err := client.GetFunds(context.Background(), mc.Client, 0)
	if err != nil {
		log.Fatalf("GetFunds failed: %v", err)
	}
	fmt.Printf("Power=%.2f Cash=%.2f Assets=%.2f\n",
		funds.Power, funds.Cash, funds.TotalAssets)
}