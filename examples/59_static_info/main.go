package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	infos, err := client.GetStaticInfo(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetStaticInfo failed: %v", err)
	}
	for _, info := range infos {
		fmt.Printf("STATIC: code=%s name=%s type=%d lotSize=%d listTime=%s\n",
			info.Code, info.Name, info.Type, info.LotSize, info.ListTime)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(infos)
}
