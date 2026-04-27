package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
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

	warrants, err := client.GetWarrant(context.Background(), cli,
		constant.Market_US, "NVDA",
		0, 10,                      // begin, num
		constant.WarrantSortField_ChangeRate, true, // sortField, ascend
		constant.WarrantType_None,  // optType=All
		qotcommon.Issuer_Issuer_Unknow, // issuer=All
		constant.WarrantStatus_None,   // status=All
	)
	if err != nil {
		log.Fatalf("GetWarrant failed: %v", err)
	}
	for _, w := range warrants {
		fmt.Printf("WARRANT: code=%s name=%s price=%.2f type=%d issuer=%d\n",
			w.Stock.GetCode(), w.Name, w.CurPrice, w.Type, w.Issuer)
	}
}
