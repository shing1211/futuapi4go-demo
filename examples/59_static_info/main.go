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

	infos, err := client.GetStaticInfo(cli, int32(constant.Market_US), "NVDA")
	if err != nil {
		log.Fatalf("GetStaticInfo failed: %v", err)
	}
	for _, info := range infos {
		fmt.Printf("STATIC: code=%s name=%s type=%d lotSize=%d listTime=%s\n",
			info.Code, info.Name, info.Type, info.LotSize, info.ListTime)
	}
}
