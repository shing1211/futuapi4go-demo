package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
)

type PriceData struct {
	Code      string
	LastPrice float64
	PrevPrice float64
	Change    float64
	ChangePct float64
	Volume    int64
}

var (
	prices      = make(map[string]*PriceData)
	priceMu    sync.Mutex
	updateCount int32
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

	fmt.Println("=== Real-Time Dashboard Demo (US Stocks) ===")
	fmt.Println()

	symbols := []string{"AAPL", "TSLA", "NVDA", "MSFT", "AMZN"}
	fmt.Printf("Monitoring: %v\n", symbols)

	for _, symbol := range symbols {
		if err := client.Subscribe(ctx, cli, constant.Market_US, symbol,
			[]constant.SubType{constant.SubType_Ticker}); err != nil {
			fmt.Printf("Subscribe %s failed: %v\n", symbol, err)
		}
	}

	tickerCh := make(chan *push.UpdateTicker, 100)
	for _, symbol := range symbols {
		go chanpkg.SubscribeTicker(cli, constant.Market_US, symbol, tickerCh)
	}

	fmt.Println("\nStarting real-time monitoring (10 seconds)...")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("%-10s %10s %10s %10s %10s\n", "Symbol", "Last", "Change", "Change%", "Volume")
	fmt.Println("─────────────────────────────────────────────────")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ticker.C:
			printPrices()
			count++
			if count >= 10 {
				goto done
			}
		case t := <-tickerCh:
			code := t.Security.GetCode()
			for _, tick := range t.TickerList {
				updatePrice(code, tick.GetPrice(), tick.GetVolume())
			}
		case <-sig:
			goto done
		}
	}

done:
	fmt.Println("\n=== Dashboard Session Complete ===")
	fmt.Printf("Total ticker updates: %d\n", updateCount)

	client.UnsubscribeAll(ctx, cli)
}

func updatePrice(code string, price float64, volume int64) {
	priceMu.Lock()
	defer priceMu.Unlock()

	updateCount++

	p, exists := prices[code]
	if !exists {
		p = &PriceData{Code: code}
		prices[code] = p
	}

	if p.LastPrice > 0 && price > 0 {
		p.PrevPrice = p.LastPrice
		p.Change = price - p.PrevPrice
		if p.PrevPrice > 0 {
			p.ChangePct = (p.Change / p.PrevPrice) * 100
		}
	}

	p.LastPrice = price
	p.Volume = volume
}

func printPrices() {
	priceMu.Lock()
	defer priceMu.Unlock()

	fmt.Println()
	for _, code := range []string{"AAPL", "TSLA", "NVDA", "MSFT", "AMZN"} {
		if p, ok := prices[code]; ok && p.LastPrice > 0 {
			fmt.Printf("%-10s %10.2f %+10.2f %+10.2f%% %10d\n",
				code, p.LastPrice, p.Change, p.ChangePct, p.Volume)
		} else {
			fmt.Printf("%-10s %10s %10s %10s %10s\n", code, "-", "-", "-", "-")
		}
	}
}