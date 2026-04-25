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

	sec := &qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("AAPL")}
	changes, err := client.GetCodeChange(context.Background(), cli, []*qotcommon.Security{sec})
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
