package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	cli := client.New(
		client.WithOnStateChange(func(old, new client.ConnState) {
			fmt.Printf("State: %d -> %d\n", old, new)
		}),
	)
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := cli.Connect("127.0.0.1:11111"); err != nil {
		log.Fatalf("connect: %v", err)
	}

	fmt.Printf("Current state: %d\n", cli.State())

	quote, err := client.GetQuote(ctx, cli, constant.Market_US, "US.AAPL")
	if err != nil {
		log.Printf("quote error: %v", err)
	} else {
		fmt.Printf("AAPL: %.2f\n", quote.Price)
	}

	if err := cli.Shutdown(5 * time.Second); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
