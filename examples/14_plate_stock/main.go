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

	// US tech sector plate
	stocks, err := client.GetPlateSecurity(context.Background(), mc.Client, constant.Market_US, "LIST20882")
	if err != nil {
		log.Fatalf("GetPlateSecurity failed: %v", err)
	}
	for _, s := range stocks {
		fmt.Printf("STOCK: code=%s name=%s type=%d\n",
			s.Code, s.Name, s.Type)
	}
}
