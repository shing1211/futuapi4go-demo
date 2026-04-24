package main

import (
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

	// HK futures: HSI (恒生指数期货) continuous contract
	infos, err := client.GetFutureInfo(cli, "HSImain")
	if err != nil {
		log.Fatalf("GetFutureInfo failed: %v", err)
	}
	for _, f := range infos {
		fmt.Printf("FUTURE: code=%s name=%s owner=%s expire=%s\n",
			f.Code, f.Name, f.Owner, f.Expire)
	}
}
