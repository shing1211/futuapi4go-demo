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

	sec := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("AAPL")}
	changes, err := client.GetCodeChange(context.Background(), mc.Client, []*qotcommon.Security{sec})
	if err != nil {
		log.Fatalf("GetCodeChange failed: %v", err)
	}
	for _, c := range changes {
		fmt.Printf("CODE CHANGE: type=%d code=%s rel=%s pub=%s eff=%s\n",
			c.Type, c.Security.GetCode(), c.RelatedSecurity.GetCode(), c.PublicTime, c.EffectiveTime)
	}
}

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string { return &v }
