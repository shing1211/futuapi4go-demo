package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetreference"
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

	refs, err := client.GetReference(context.Background(), cli, constant.Market_US, "NVDA", qotgetreference.ReferenceType_ReferenceType_Warrant)
	if err != nil {
		log.Fatalf("GetReference failed: %v", err)
	}
	for _, r := range refs {
		fmt.Printf("REFERENCE: code=%s name=%s type=%d\n",
			r.Code, r.Name, r.Type)
	}
}
