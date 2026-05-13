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

	plates, err := client.GetOwnerPlate(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetOwnerPlate failed: %v", err)
	}
	for _, p := range plates {
		fmt.Printf("OWNER PLATE: %s\n", p)
	}
}
