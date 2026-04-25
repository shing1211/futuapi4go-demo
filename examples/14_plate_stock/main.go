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

	// US tech sector plate
	stocks, err := client.GetPlateSecurity(context.Background(), cli, int32(constant.Market_US), "LIST20882")
	if err != nil {
		log.Fatalf("GetPlateSecurity failed: %v", err)
	}
	for _, s := range stocks {
		fmt.Printf("STOCK: code=%s name=%s type=%d\n",
			s.Code, s.Name, s.Type)
	}
}
