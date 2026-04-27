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

	plates, err := client.GetPlateSet(context.Background(), cli, constant.Market_US)
	if err != nil {
		log.Fatalf("GetPlateSet failed: %v", err)
	}
	for _, p := range plates {
		fmt.Printf("PLATE: code=%s name=%s\n", p.Code, p.Name)
	}
}
