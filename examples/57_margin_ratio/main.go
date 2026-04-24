package main

import (
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

	accounts, err := client.GetAccountList(cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	market := acc.TrdMarketAuthList[0]

	sec := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")}
	ratios, err := client.GetMarginRatio(cli,
		accID,
		market,
		[]*qotcommon.Security{sec},
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
