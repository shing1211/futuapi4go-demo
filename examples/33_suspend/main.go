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

	sec := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("NVDA")}
	susp, err := client.GetSuspend(context.Background(), mc.Client, []*qotcommon.Security{sec}, "2026-01-01", "2026-04-24")
	if err != nil {
		log.Fatalf("GetSuspend failed: %v", err)
	}
	for _, s := range susp {
		fmt.Printf("SUSPEND: time=%s ts=%.0f\n",
			s.Time, s.Timestamp)
	}
}

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string { return &v }
