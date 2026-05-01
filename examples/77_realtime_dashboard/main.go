package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	futoclient "github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/proto"
)

type PriceData struct {
	Code      string
	LastPrice float64
	PrevPrice float64
	Change    float64
	ChangePct float64
	Volume    int64
	Timestamp time.Time
}

var (
	prices      = make(map[string]*PriceData)
	priceMu     sync.Mutex
	updateCount int32
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println("=== Real-Time Dashboard Demo (US Stocks) ===")
	fmt.Println()

	symbols := []string{"AAPL", "TSLA", "NVDA", "MSFT", "AMZN"}
	fmt.Printf("Monitoring: %v\n", symbols)

	cli.RegisterHandler(proto.ProtoID_Qot_UpdateTicker, func(protoID uint32, body []byte) {
		if ticker, err := futoclient.ParsePushTicker(body); err == nil {
			updatePrice(ticker.Code, ticker.LastPrice, ticker.Volume)
		}
	})

	for _, symbol := range symbols {
		if err := futoclient.Subscribe(ctx, cli, constant.Market_US, symbol, constant.SubType_Ticker); err != nil {
			fmt.Printf("Subscribe %s failed: %v\n", symbol, err)
		}
	}

	fmt.Println("\nStarting real-time monitoring (10 seconds)...")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("%-10s %10s %10s %10s %10s\n", "Symbol", "Last", "Change", "Change%", "Volume")
	fmt.Println("─────────────────────────────────────────────────")

	stopCh = make(chan struct{})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ticker.C:
			printDashboard()
			count++
			if count >= 10 {
				close(stopCh)
				return
			}
		case <-ctx.Done():
			close(stopCh)
			return
		}
	}
}

func updatePrice(code string, price float64, volume int64) {
	priceMu.Lock()
	defer priceMu.Unlock()

	if _, exists := prices[code]; !exists {
		prices[code] = &PriceData{Code: code}
	}

	p := prices[code]
	prevPrice := p.LastPrice
	p.LastPrice = price
	p.Volume = volume
	p.Timestamp = time.Now()

	if prevPrice > 0 {
		p.Change = price - prevPrice
		if prevPrice > 0 {
			p.ChangePct = (p.Change / prevPrice) * 100
		}
	}
}

func printDashboard() {
	priceMu.Lock()
	defer priceMu.Unlock()

	fmt.Println("─────────────────────────────────────────────────")
	for code, p := range prices {
		changeStr := fmt.Sprintf("%+.2f", p.Change)
		pctStr := fmt.Sprintf("%+.2f%%", p.ChangePct)
		fmt.Printf("%-10s %10.2f %10s %10s %10d\n",
			code, p.LastPrice, changeStr, pctStr, p.Volume)
	}
	if len(prices) == 0 {
		fmt.Println("Waiting for data...")
	}
}