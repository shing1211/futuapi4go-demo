package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	key, err := client.SetPriceReminder(context.Background(), mc.Client,
		constant.Market_US, "NVDA",
		constant.PriceReminderOp_Add,         // op: 1=Add
		constant.PriceReminderType_Above,     // reminderType: 1=Above
		constant.PriceReminderFreq_Once,     // freq: 0=Once
		150.0,                               // value: trigger when price reaches 150
		"Watch for NVDA at 150",
	)
	if err != nil {
		log.Fatalf("SetPriceReminder failed: %v", err)
	}
	fmt.Printf("Price reminder set! Key: %d\n", key)
}
