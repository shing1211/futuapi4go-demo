package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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

	accounts, err := client.GetAccountList(context.Background(), cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	_ = accID // unused in this example

	fills, err := client.GetOrderFillList(cli, accID)
	if err != nil {
		log.Fatalf("GetOrderFillList failed: %v", err)
	}
	if len(fills) == 0 {
		fmt.Println("(no order fills)")
	}
	for _, f := range fills {
		fmt.Printf("FILL: id=%d %s price=%.2f qty=%.0f time=%s\n",
			f.FillID, f.Code, f.Price, f.Qty, f.CreateTime)
	}
}
