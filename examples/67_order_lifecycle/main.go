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

	ctx := context.Background()

	fmt.Println("=== Order Lifecycle Demo (US Simulated) ===")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, cli)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	var accMarket constant.TrdMarket
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			for _, auth := range acc.TrdMarketAuthList {
				if auth == constant.TrdMarket_US.Int32() {
					accID = acc.AccID
					accMarket = constant.TrdMarket_US
					fmt.Printf("Using US simulated AccID=%d\n", accID)
					break
				}
			}
		}
		if accID != 0 {
			break
		}
	}

	if accID == 0 {
		for _, acc := range accounts {
			if acc.TrdEnv == 0 {
				for _, auth := range acc.TrdMarketAuthList {
					if auth == constant.TrdMarket_HK.Int32() {
						accID = acc.AccID
						accMarket = constant.TrdMarket_HK
						fmt.Printf("Using HK simulated AccID=%d\n", accID)
						break
					}
				}
			}
			if accID != 0 {
				break
			}
		}
	}

	if accID == 0 {
		accID = accounts[0].AccID
		accMarket = constant.TrdMarket(accounts[0].TrdMarketAuthList[0])
		fmt.Printf("Using AccID=%d (market=%d)\n", accID, accMarket)
	}

	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, cli, pwd); err != nil {
			fmt.Printf("UnlockTrading warning: %v\n", err)
		} else {
			fmt.Println("Trading unlocked successfully")
		}
	}

	fmt.Println("\n=== Step 1: Check Account Funds ===")
	funds, err := client.GetFunds(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetFunds failed: %v\n", err)
	} else {
		fmt.Printf("Cash: %.2f | Power: %.2f | Assets: %.2f\n",
			funds.Cash, funds.Power, funds.TotalAssets)
	}

	fmt.Println("\n=== Step 2: Check Current Positions ===")
	positions, err := client.GetPositionList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetPositionList failed: %v\n", err)
	} else {
		fmt.Printf("Current positions: %d\n", len(positions))
		for _, p := range positions {
			fmt.Printf("  %s: Qty=%.0f Cost=%.2f Cur=%.2f\n",
				p.Code, p.Quantity, p.CostPrice, p.CurPrice)
		}
	}

	fmt.Println("\n=== Step 3: Place a Limit Order (Demo) ===")
	var stock string
	var price float64
	var secMarket constant.TrdSecMarket
	var qty float64

	switch accMarket {
	case constant.TrdMarket_US:
		stock = "AAPL"
		price = 180.0
		secMarket = constant.TrdSecMarket_US
		qty = 1.0
	case constant.TrdMarket_HK:
		stock = "00100"
		price = 700.0
		secMarket = constant.TrdSecMarket_HK
		qty = 100.0
	default:
		stock = "AAPL"
		price = 180.0
		secMarket = constant.TrdSecMarket_US
		qty = 1.0
	}

	fmt.Printf("Placing order: Buy %.0f share(s) of %s @ $%.2f\n", qty, stock, price)

	result, err := client.PlaceOrder(ctx, cli, accID, accMarket,
		stock, constant.TrdSide_Buy, constant.OrderType_Normal, price, qty, secMarket)
	if err != nil {
		fmt.Printf("PlaceOrder failed: %v\n", err)
	} else {
		fmt.Printf("Order placed successfully! OrderID=%d\n", result.OrderID)
	}

	fmt.Println("\n=== Step 4: List Open Orders ===")
	orders, err := client.GetOrderList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetOrderList failed: %v\n", err)
	} else {
		fmt.Printf("Open orders: %d\n", len(orders))
		for _, o := range orders {
			fmt.Printf("  OrderID=%d %s TrdSide=%d Qty=%.0f Price=%.2f OrderStatus=%d\n",
				o.OrderID, o.Code, o.TrdSide, o.Qty, o.Price, o.OrderStatus)
		}
	}

	fmt.Println("\n=== Step 5: Modify/Cancel Order (Demo) ===")
	if len(orders) > 0 {
		orderID := orders[0].OrderID
		newPrice := price * 1.05
		fmt.Printf("Modifying order %d to price $%.2f...\n", orderID, newPrice)

		modResult, err := client.ModifyOrder(ctx, cli, accID, accMarket,
			orderID, constant.ModifyOrderOp_Normal, newPrice, 0)
		if err != nil {
			fmt.Printf("ModifyOrder failed: %v\n", err)
		} else {
			fmt.Printf("Modify result: OrderID=%d\n", modResult.OrderID)
		}
	}

	fmt.Println("\n=== Order Lifecycle Complete ===")
}