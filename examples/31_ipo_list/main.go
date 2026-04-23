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

	ipos, err := client.GetIpoList(cli, int32(constant.Market_US))
	if err != nil {
		log.Fatalf("GetIpoList failed: %v", err)
	}
	for _, ip := range ipos {
		fmt.Printf("IPO: code=%s name=%s listDate=%s\n",
			ip.Code, ip.Name, ip.ListDate)
	}
}
