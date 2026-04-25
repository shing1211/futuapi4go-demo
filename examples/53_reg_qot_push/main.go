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

	if err := client.RegQotPush(context.Background(), cli,
		int32(constant.Market_US), "NVDA",
		[]int32{int32(constant.SubType_Quote)},
		[]int32{0},
		true,  // isReg
		true,  // isFirstPush
	); err != nil {
		log.Fatalf("RegQotPush failed: %v", err)
	}
	fmt.Println("RegQotPush registered for NVDA quote.")
}
