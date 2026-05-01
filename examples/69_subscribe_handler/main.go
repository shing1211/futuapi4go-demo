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
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
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

	fmt.Println("=== Comprehensive Subscribe Handler Demo (US Stocks) ===")
	fmt.Println()

	symbols := []string{"AAPL", "TSLA", "NVDA"}
	fmt.Printf("Subscribing to %d symbols: %v\n", len(symbols), symbols)

	for _, symbol := range symbols {
		if err := client.Subscribe(ctx, cli, constant.Market_US, symbol,
			[]constant.SubType{constant.SubType_Ticker}); err != nil {
			fmt.Printf("Subscribe %s ticker failed: %v\n", symbol, err)
		}
	}

	if err := client.Subscribe(ctx, cli, constant.Market_US, symbols[0],
		[]constant.SubType{constant.SubType_K_Day, constant.SubType_OrderBook}); err != nil {
		fmt.Printf("Subscribe %s kline/orderbook failed: %v\n", symbols[0], err)
	}

	tickerCh := make(chan *push.UpdateTicker, 100)
	klineCh := make(chan *push.UpdateKL, 100)
	orderbookCh := make(chan *push.UpdateOrderBook, 100)

	go chanpkg.SubscribeTicker(cli, constant.Market_US, symbols[0], tickerCh)
	go chanpkg.SubscribeKLine(cli, constant.Market_US, symbols[0], constant.KLType_K_Day, klineCh)
	go chanpkg.SubscribeOrderBook(cli, constant.Market_US, symbols[0], orderbookCh)

	fmt.Println("\nWaiting for push data (5 seconds)...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	tickerCount := 0
	klineCount := 0
	orderbookCount := 0

	fmt.Println("─────────────────────────────────────────────────")

	for {
		select {
		case t := <-tickerCh:
			for _, tick := range t.TickerList {
				if tickerCount < 3 {
					fmt.Printf("[Ticker] %s: price=%.2f vol=%d dir=%d\n",
						symbols[0], tick.GetPrice(), tick.GetVolume(), tick.GetDir())
				}
				tickerCount++
			}
		case kl := <-klineCh:
			for _, kline := range kl.KLList {
				if klineCount < 3 {
					fmt.Printf("[KL] %s: O=%.2f H=%.2f L=%.2f C=%.2f\n",
						symbols[0], kline.GetOpenPrice(), kline.GetHighPrice(), kline.GetLowPrice(), kline.GetClosePrice())
				}
				klineCount++
			}
		case ob := <-orderbookCh:
			if orderbookCount < 2 {
				fmt.Printf("[OrderBook] %s: %d bids, %d asks\n",
					symbols[0], len(ob.OrderBookBidList), len(ob.OrderBookAskList))
			}
			orderbookCount++
		case <-time.After(5 * time.Second):
			goto done
		case <-sig:
			goto done
		}
	}

done:
	fmt.Println("\n=== Summary ===")
	fmt.Printf("Received: %d tickers, %d klines, %d orderbooks\n",
		tickerCount, klineCount, orderbookCount)

	client.UnsubscribeAll(ctx, cli)
}