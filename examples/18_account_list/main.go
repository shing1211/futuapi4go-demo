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
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	for _, acc := range accounts {
		fmt.Printf("ACC: id=%d type=%d env=%d firm=%d\n",
			acc.AccID, acc.AccType, acc.TrdEnv, acc.SecurityFirm)
	}
}
