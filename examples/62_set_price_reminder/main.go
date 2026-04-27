package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
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

	key, err := client.SetPriceReminder(context.Background(), cli,
		constant.Market_US, "NVDA",
		constant.PriceReminderOp_Add,         // op: 1=Add
		constant.PriceReminderType_Price,     // reminderType: 1=Price above or below
		constant.PriceReminderFreq_Once,     // freq: 0=Once
		150.0,                               // value: trigger when price reaches 150
		"Watch for NVDA at 150",
	)
	if err != nil {
		log.Fatalf("SetPriceReminder failed: %v", err)
	}
	fmt.Printf("Price reminder set! Key: %d\n", key)
}
