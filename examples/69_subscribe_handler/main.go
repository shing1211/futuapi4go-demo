package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	futoclient "github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/proto"
)

var (
	tickerCount int32
	klineCount  int32
	orderbookCount int32
	stopCh      chan struct{}
)

func main() {
	cli := futoclient.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== Comprehensive Subscribe Handler Demo (US Stocks) ===")
	fmt.Println()

	// Create channels for push data
	tickerCh := make(chan *futoclient.PushTicker, 100)
	klineCh := make(chan *futoclient.PushKLine, 100)
	orderbookCh := make(chan *futoclient.PushOrderBook, 100)

	// Register push handlers
	cli.RegisterHandler(proto.ProtoID_Qot_UpdateTicker, func(protoID uint32, body []byte) {
		if ticker, err := futoclient.ParsePushTicker(body); err == nil {
			select {
			case tickerCh <- ticker:
			default:
			}
		}
	})

	cli.RegisterHandler(proto.ProtoID_Qot_UpdateKL, func(protoID uint32, body []byte) {
		if kline, err := futoclient.ParsePushKLine(body); err == nil {
			select {
			case klineCh <- kline:
			default:
			}
		}
	})

	cli.RegisterHandler(proto.ProtoID_Qot_UpdateOrderBook, func(protoID uint32, body []byte) {
		if ob, err := futoclient.ParsePushOrderBook(body); err == nil {
			select {
			case orderbookCh <- ob:
			default:
			}
		}
	})

	// Subscribe to US stocks
	symbols := []string{"AAPL", "TSLA", "NVDA"}

	fmt.Printf("Subscribing to %d symbols: %v\n", len(symbols), symbols)
	if err := futoclient.Subscribe(ctx, cli, constant.Market_US, symbols[0],
		constant.SubType_Ticker, constant.SubType_KL_Day, constant.SubType_OrderBook); err != nil {
		fmt.Printf("Subscribe failed: %v\n", err)
	}

	// Subscribe remaining symbols to ticker only
	for i := 1; i < len(symbols); i++ {
		futoclient.Subscribe(ctx, cli, constant.Market_US, symbols[i], constant.SubType_Ticker)
	}

	fmt.Println("\nWaiting for push data (5 seconds)...")
	fmt.Println("Press Ctrl+C to stop early")

	stopCh = make(chan struct{})

	// Start goroutines to process push data
	go processTickers(tickerCh)
	go processKlines(klineCh)
	go processOrderBooks(orderbookCh)

	// Wait for stop signal or timeout
	select {
	case <-time.After(5 * time.Second):
		close(stopCh)
	case <-ctx.Done():
	}

	// Summary
	fmt.Println("\n=== Summary ===")
	fmt.Printf("Received: %d tickers, %d klines, %d orderbooks\n",
		tickerCount, klineCount, orderbookCount)

	// Cleanup
	futoclient.UnsubscribeAll(ctx, cli)
}

func processTickers(ch <-chan *futoclient.PushTicker) {
	for {
		select {
		case <-stopCh:
			return
		case ticker := <-ch:
			atomic.AddInt32(&tickerCount, 1)
			if atomic.LoadInt32(&tickerCount) <= 3 {
				fmt.Printf("[Ticker] %s: Last=%d Vol=%d\n",
					ticker.Code, ticker.LastPrice, ticker.Volume)
			}
		}
	}
}

func processKlines(ch <-chan *futoclient.PushKLine) {
	for {
		select {
		case <-stopCh:
			return
		case kline := <-ch:
			atomic.AddInt32(&klineCount, 1)
			if atomic.LoadInt32(&klineCount) <= 3 {
				fmt.Printf("[KL] %s: O=%.2f H=%.2f L=%.2f C=%.2f\n",
					kline.Code, kline.Open, kline.High, kline.Low, kline.Close)
			}
		}
	}
}

func processOrderBooks(ch <-chan *futoclient.PushOrderBook) {
	for {
		select {
		case <-stopCh:
			return
		case ob := <-ch:
			atomic.AddInt32(&orderbookCount, 1)
			if atomic.LoadInt32(&orderbookCount) <= 2 {
				fmt.Printf("[OrderBook] %s: %d levels\n", ob.Code, len(ob.BidItems)+len(ob.AskItems))
			}
		}
	}
}