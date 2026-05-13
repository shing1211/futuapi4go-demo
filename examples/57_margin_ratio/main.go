package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := mc.Client.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	_ = accID

	sec := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")}
	ratios, err := client.GetMarginRatio(context.Background(), mc.Client,
	accID, constant.TrdMarket_HK, []*qotcommon.Security{sec},
)
	if err != nil {
		log.Fatalf("GetMarginRatio failed: %v", err)
	}
	for _, r := range ratios {
		fmt.Printf("Margin: long=%v short=%v shortFee=%.4f longRatio=%.4f shortRatio=%.4f\n",
			r.IsLongPermit, r.IsShortPermit, r.ShortFeeRate, r.ImLongRatio, r.ImShortRatio)
	}
}

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string { return &v }
