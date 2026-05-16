package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/trd"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	audit := trd.NewAuditLogger(logger)

	cli := client.New()
	defer cli.Close()

	if err := cli.Connect("127.0.0.1:11111"); err != nil {
		log.Fatalf("connect: %v", err)
	}

	req := &trd.PlaceOrderRequest{
		AccID:     12345,
		Code:      "US.AAPL",
		TrdSide:   constant.TrdSide_Buy,
		OrderType: constant.OrderType_Normal,
		Price:     150.0,
		Qty:       100.0,
	}

	warnings := trd.ValidateOrder(&trd.OrderValidationInput{
		Order:       req,
		MarketOpen:  false,
		BuyingPower: 50000,
		MaxBuyQty:   1000,
		MaxSellQty:  1000,
	})

	if trd.HasErrors(warnings) {
		for _, w := range warnings {
			log.Printf("FAIL: %s\n", w.Message)
		}
	}

	time.Sleep(time.Second)
	audit.LogPlaceOrder(req, 0, nil)
}
