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

	if err := client.RegQotPush(context.Background(), mc.Client,
		constant.Market_US, "NVDA",
		[]constant.SubType{constant.SubType_Quote},
		[]constant.RehabType{constant.RehabType_None},
		true,  // isReg
		true,  // isFirstPush
	); err != nil {
		log.Fatalf("RegQotPush failed: %v", err)
	}
	fmt.Println("RegQotPush registered for NVDA quote.")
}
