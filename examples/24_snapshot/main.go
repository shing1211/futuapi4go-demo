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

	sec1 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")}
	sec2 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("AAPL")}
	sec3 := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("TSLA")}

	snapshots, err := client.GetSecuritySnapshot(context.Background(), cli, []*qotcommon.Security{sec1, sec2, sec3})
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
