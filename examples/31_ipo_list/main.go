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
		if ip.CnExData != nil {
			fmt.Printf("  CN: applyCode=%s issueSize=%d ipoPrice=%.2f winningRatio=%.4f\n",
				ip.CnExData.GetApplyCode(), ip.CnExData.GetIssueSize(),
				ip.CnExData.GetIpoPrice(), ip.CnExData.GetWinningRatio())
		}
		if ip.HkExData != nil {
			fmt.Printf("  HK: priceMin=%.2f priceMax=%.2f listPrice=%.2f lotSize=%d\n",
				ip.HkExData.GetIpoPriceMin(), ip.HkExData.GetIpoPriceMax(),
				ip.HkExData.GetListPrice(), ip.HkExData.GetLotSize())
		}
		if ip.UsExData != nil {
			fmt.Printf("  US: priceMin=%.2f priceMax=%.2f issueSize=%d\n",
				ip.UsExData.GetIpoPriceMin(), ip.UsExData.GetIpoPriceMax(),
				ip.UsExData.GetIssueSize())
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(ipos)
}
