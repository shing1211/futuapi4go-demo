package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil {
		fmt.Printf("GetAccountList failed: %v\n", err)
		return
	}
	for _, acc := range accounts {
		fmt.Printf("ACC: id=%d type=%d env=%d firm=%d\n",
			acc.AccID, acc.AccType, acc.TrdEnv, acc.SecurityFirm)
	}
}
