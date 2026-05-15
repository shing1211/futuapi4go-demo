package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Trade Push Monitor Demo ===")
	fmt.Println("Subscribes to trade push events and waits for updates")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, mc.Client)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			accID = acc.AccID
			break
		}
	}
	if accID == 0 {
		accID = accounts[0].AccID
	}
	fmt.Printf("Using AccID=%d\n", accID)

	accIDList := make([]uint64, len(accounts))
	for i, acc := range accounts {
		accIDList[i] = acc.AccID
	}

	fmt.Println("\n--- Subscribing to Trade Push ---")
	if err := client.SubAccPush(ctx, mc.Client, accIDList); err != nil {
		log.Fatalf("SubAccPush failed: %v", err)
	}
	fmt.Printf("Subscribed %d account(s) for trade push\n", len(accIDList))

	type pushStats struct {
		orderUpdates int
		orderFills   int
		priceAlerts  int
	}

	stats := &pushStats{}

	fmt.Println("\n--- Registering Push Handlers ---")

	mc.Client.SetPushHandler(constant.ProtoID_Trd_UpdateOrder, func(protoID uint32, body []byte) {
		stats.orderUpdates++
		update, err := client.ParsePushOrderUpdate(body)
		if err != nil {
			fmt.Printf("[Push Error] Parse order update: %v\n", err)
			return
		}
		fmt.Printf("[Order:%d] %s side=%d qty=%.0f price=%.2f status=%d\n",
			update.OrderID, update.Code, update.TrdSide,
			update.Qty, update.Price, update.OrderStatus)
	})

	mc.Client.SetPushHandler(constant.ProtoID_Trd_UpdateOrderFill, func(protoID uint32, body []byte) {
		stats.orderFills++
		fill, err := client.ParsePushOrderFill(body)
		if err != nil {
			fmt.Printf("[Push Error] Parse order fill: %v\n", err)
			return
		}
		fmt.Printf("[Fill:%d] Order %d: %.0f shares @ $%.2f\n",
			fill.FillID, fill.OrderID, fill.Qty, fill.Price)
	})

	mc.Client.SetPushHandler(constant.ProtoID_Qot_UpdatePriceReminder, func(protoID uint32, body []byte) {
		stats.priceAlerts++
		alert, err := client.ParsePushPriceReminder(body)
		if err != nil {
			fmt.Printf("[Push Error] Parse price reminder: %v\n", err)
			return
		}
		fmt.Printf("[PriceAlert] %s: cur=%.2f type=%d\n",
			alert.Code, alert.CurValue, alert.Type)
	})

	fmt.Println("Handlers registered for ProtoID 2208 (order), 2218 (fill), 3019 (price)")
	fmt.Println()

	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, mc.Client, pwd); err != nil {
			fmt.Printf("UnlockTrading warning: %v\n", err)
		}

		fmt.Println("--- Placing test order to trigger push events ---")
		result, err := client.PlaceOrder(ctx, mc.Client, accID, constant.TrdMarket_US,
			"US.NVDA", constant.TrdSide_Buy, constant.OrderType_Normal,
			120.0, 1, constant.TrdSecMarket_US)
		if err != nil {
			fmt.Printf("PlaceOrder (may trigger rejected push): %v\n", err)
		} else {
			fmt.Printf("Order placed: %d\n", result.OrderID)
		}
	} else {
		fmt.Println("No FUTU_TRADE_PWD set — skipping test order placement")
		fmt.Println("Handlers will still listen for pushes from other sources")
	}

	fmt.Println("\n--- Waiting for Push Events (60s) ---")
	fmt.Println("Press Ctrl+C to stop early")
	fmt.Println()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(60 * time.Second)

loop:
	for {
		select {
		case <-sig:
			fmt.Println("\nInterrupted by user")
			break loop
		case <-timeout:
			fmt.Println("\n60s timeout reached")
			break loop
		case <-ticker.C:
			fmt.Printf("  listening... (updates=%d fills=%d alerts=%d)\n",
				stats.orderUpdates, stats.orderFills, stats.priceAlerts)
		}
	}

	fmt.Println()
	fmt.Println("--- Push Event Summary ---")
	fmt.Printf("Order Updates:  %d\n", stats.orderUpdates)
	fmt.Printf("Order Fills:    %d\n", stats.orderFills)
	fmt.Printf("Price Alerts:   %d\n", stats.priceAlerts)
	fmt.Printf("Total:          %d\n", stats.orderUpdates+stats.orderFills+stats.priceAlerts)

	if stats.orderUpdates == 0 && stats.orderFills == 0 {
		fmt.Println()
		fmt.Println("No pushes received — expected for simulate mode.")
		fmt.Println("Push events require:")
		fmt.Println("  1. Real trading environment (TrdEnv_Real)")
		fmt.Println("  2. FUTU_TRADE_PWD for trade unlock")
		fmt.Println("  3. Active order placement/modification/cancellation")
	}

	fmt.Println("\n=== Monitor Complete ===")
}
