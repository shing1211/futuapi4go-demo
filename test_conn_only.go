//go:build ignore

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shing1211/futuapi4go/internal/client"
)

func main() {
	rsaKey := "/etc/futu/keys/private_key.pem"
	addr := "172.18.208.88:11111"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cli := client.New(client.WithRSAPublicKey(rsaKey))

	fmt.Println("Connecting...")
	if err := cli.Connect(addr); err != nil {
		fmt.Printf("Connect failed: %v\n", err)
		return
	}
	fmt.Printf("Connected! connID=%d encrypt=%v aesKey=%d\n",
		cli.GetConnID(), cli.IsEncrypt(), len(cli.GetAESKey()))

	fmt.Println("Closing...")
	cli.Close()
	fmt.Println("Done — clean exit")
}