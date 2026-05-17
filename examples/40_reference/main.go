package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetreference"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	refs, err := client.GetReference(context.Background(), mc.Client, constant.Market_US, "NVDA", qotgetreference.ReferenceType_ReferenceType_Warrant)
	if err != nil {
		log.Fatalf("GetReference failed: %v", err)
	}
	for _, r := range refs {
		fmt.Printf("REFERENCE: code=%s name=%s type=%d\n",
			r.Code, r.Name, r.Type)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(refs)
}
