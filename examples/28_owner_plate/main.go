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

	plates, err := client.GetOwnerPlate(cli, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetOwnerPlate failed: %v", err)
	}
	for _, p := range plates {
		fmt.Printf("OWNER PLATE: %s\n", p)
	}
}
