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

	sec1 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")}
	sec2 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("AAPL")}
	sec3 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("TSLA")}

	snapshots, err := client.GetSecuritySnapshot(context.Background(), mc.Client, []*qotcommon.Security{sec1, sec2, sec3})
	if err != nil {
		log.Fatalf("GetSecuritySnapshot failed: %v", err)
	}
	for _, s := range snapshots {
		fmt.Printf("SNAP: %s %s price=%.2f open=%.2f high=%.2f low=%.2f vol=%d\n",
			s.Security.GetCode(), s.Name, s.CurPrice, s.OpenPrice, s.HighPrice, s.LowPrice, s.Volume)
	}
}

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string { return &v }
