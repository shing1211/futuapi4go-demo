package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ipos, err := client.GetIpoList(context.Background(), mc.Client, constant.Market_HK)
	if err != nil {
		log.Fatalf("GetIpoList failed: %v", err)
	}
	for _, ip := range ipos {
		fmt.Printf("IPO: code=%s name=%s listDate=%s\n",
			ip.Code, ip.Name, ip.ListDate)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(ipos)
}
