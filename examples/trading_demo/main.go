// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example: Trading Operations - Account and order management
//
// This example demonstrates trading functionality:
// - List trading accounts
// - Query account funds and positions
// - View order history
//
// Run: go run examples/trading_demo
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/trd"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}

	fmt.Println("=== Trading Operations Example ===")
	fmt.Printf("Connecting to %s...\n", addr)

	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	fmt.Printf("Connected! ConnID=%d\n\n", cli.GetConnID())

	runAccountList(cli)
	runFundsQuery(cli)
	runPositionList(cli)
	runOrderList(cli)

	fmt.Println("\n=== Example Complete ===")
}

func runAccountList(cli *client.Client) {
	fmt.Println("--- Account List ---")

	resp, err := trd.GetAccList(cli.Inner(), int32(constant.TrdCategory_Security), false)
	if err != nil {
		log.Printf("GetAccList failed: %v", err)
		return
	}

	for _, acc := range resp.AccList {
		fmt.Printf("  AccountID: %d (AuthList: %v)\n", acc.AccID, acc.TrdMarketAuthList)
	}
}

func runFundsQuery(cli *client.Client) {
	fmt.Println("\n--- Account Funds ---")

	resp, err := trd.GetFunds(cli.Inner(), &trd.GetFundsRequest{
		AccID:     0,
		TrdMarket: int32(constant.TrdMarket_HK),
		TrdEnv:    int32(constant.TrdEnv_Simulate),
	})
	if err != nil {
		log.Printf("GetFunds failed: %v", err)
		return
	}

	funds := resp.Funds
	fmt.Printf("  Cash: %.2f %s\n", funds.Cash, mapCurrency(funds.Currency))
	fmt.Printf("  Available: %.2f\n", funds.AvlWithdrawalCash)
	fmt.Printf("  Total Assets: %.2f\n", funds.TotalAssets)
}

func runPositionList(cli *client.Client) {
	fmt.Println("\n--- Position List ---")

	resp, err := trd.GetPositionList(cli.Inner(), &trd.GetPositionListRequest{
		AccID:     0,
		TrdMarket: int32(constant.TrdMarket_HK),
		TrdEnv:    int32(constant.TrdEnv_Simulate),
	})
	if err != nil {
		log.Printf("GetPositionList failed: %v", err)
		return
	}

	if len(resp.PositionList) == 0 {
		fmt.Println("  (No positions)")
		return
	}

	for _, pos := range resp.PositionList {
		fmt.Printf("  %s: %.0f shares @ avg %.2f\n",
			pos.Code, pos.Qty, pos.CostPrice)
	}
}

func runOrderList(cli *client.Client) {
	fmt.Println("\n--- Order History ---")

	resp, err := trd.GetOrderList(cli.Inner(), &trd.GetOrderListRequest{
		AccID:     0,
		TrdMarket: int32(constant.TrdMarket_HK),
		TrdEnv:    int32(constant.TrdEnv_Simulate),
	})
	if err != nil {
		log.Printf("GetOrderList failed: %v", err)
		return
	}

	if len(resp.OrderList) == 0 {
		fmt.Println("  (No orders)")
		return
	}

	for _, order := range resp.OrderList {
		fmt.Printf("  %s: %s %.0f@%.2f [%s]\n",
			order.Code,
			mapTrdSide(order.TrdSide),
			order.Qty,
			order.Price,
			mapOrderStatus(order.OrderStatus))
	}
}

func mapCurrency(currency int32) string {
	switch currency {
	case 1:
		return "HKD"
	case 2:
		return "USD"
	case 3:
		return "CNH"
	case 4:
		return "JPY"
	case 5:
		return "SGD"
	default:
		return "Unknown"
	}
}

func mapTrdSide(side int32) string {
	switch side {
	case int32(constant.TrdSide_Buy):
		return "BUY"
	case int32(constant.TrdSide_Sell):
		return "SELL"
	default:
		return "UNKNOWN"
	}
}

func mapOrderStatus(status int32) string {
	switch status {
	case int32(constant.OrderStatus_WaitingSubmit):
		return "Waiting"
	case int32(constant.OrderStatus_Submitted):
		return "Submitted"
	case int32(constant.OrderStatus_FilledPart):
		return "Partial Fill"
	case int32(constant.OrderStatus_FilledAll):
		return "Filled"
	case int32(constant.OrderStatus_CancelledPart):
		return "Partial Cancel"
	case int32(constant.OrderStatus_CancelledAll):
		return "Cancelled"
	default:
		return fmt.Sprintf("Status(%d)", status)
	}
}
