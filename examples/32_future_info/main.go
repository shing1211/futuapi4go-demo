package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	// HK futures: HSI (恒生指数期货) continuous contract
	infos, err := client.GetFutureInfo(context.Background(), mc.Client, "HSImain")
	if err != nil {
		log.Fatalf("GetFutureInfo failed: %v", err)
	}
	for _, f := range infos {
		fmt.Printf("FUTURE: code=%s name=%s owner=%s expire=%s\n",
			f.Code, f.Name, f.Owner, f.Expire)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(infos)
}
