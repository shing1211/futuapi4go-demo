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

	plates, err := client.GetPlateSet(context.Background(), mc.Client, constant.Market_US)
	if err != nil {
		log.Fatalf("GetPlateSet failed: %v", err)
	}
	for _, p := range plates {
		fmt.Printf("PLATE: code=%s name=%s type=%d\n", p.Code, p.Name, p.PlateType)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(plates)
}
