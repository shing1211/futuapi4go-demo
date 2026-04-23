package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
)

func main() {
	cli := client.New()
	defer cli.Close()

	if err := cli.Connect("127.0.0.1:11111"); err != nil {
		panic(err)
	}

	// Set up channel listeners for each data type
	quoteCh := make(chan *push.UpdateBasicQot, 100)
	tickerCh := make(chan *push.UpdateTicker, 100)
	orderBookCh := make(chan *push.UpdateOrderBook, 100)
	rtCh := make(chan *push.UpdateRT, 100)
	brokerCh := make(chan *push.UpdateBroker, 100)
	klCh := make(chan *push.UpdateKL, 100)

	chanpkg.SubscribeQuote(cli, constant.Market_US, "NVDA", quoteCh)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for all NVDA data (Ctrl+C to exit)...")
	for {
		select {
		case q := <-quoteCh:
			fmt.Printf("QUOTE [%s]: price=%.2f vol=%d\n",
				q.Security.GetCode(), q.CurPrice, q.Volume)
		case t := <-tickerCh:
			if len(t.TickerList) > 0 {
				fmt.Printf("TICKER: price=%.2f vol=%d\n",
					t.TickerList[0].GetPrice(), t.TickerList[0].GetVolume())
			}
		case ob := <-orderBookCh:
			if len(ob.OrderBookBidList) > 0 && len(ob.OrderBookAskList) > 0 {
				fmt.Printf("ORDERBOOK: bid=%.2f ask=%.2f\n",
					ob.OrderBookBidList[0].GetPrice(), ob.OrderBookAskList[0].GetPrice())
			}
		case rt := <-rtCh:
			if len(rt.RTList) > 0 {
				fmt.Printf("RT: price=%.2f avg=%.2f\n",
					rt.RTList[0].GetPrice(), rt.RTList[0].GetAvgPrice())
			}
		case b := <-brokerCh:
			if len(b.BidBrokerList) > 0 {
				fmt.Printf("BROKER: name=%s pos=%d\n",
					b.BidBrokerList[0].GetName(), b.BidBrokerList[0].GetPos())
			}
		case kl := <-klCh:
			for _, bar := range kl.KLList {
				fmt.Printf("KL: %s O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
					*bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice, *bar.Volume)
			}
		case <-sig:
			fmt.Println("Shutting down...")
			return
		}
	}
}
