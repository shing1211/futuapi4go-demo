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

	plates, err := client.GetOwnerPlate(context.Background(), mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetOwnerPlate failed: %v", err)
	}
	for code, plateList := range plates {
		fmt.Printf("Security: %s\n", code)
		for _, p := range plateList {
			fmt.Printf("  Plate: %s (code=%s, type=%d)\n", p.Name, p.Code, p.PlateType)
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(plates)
}
