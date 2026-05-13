package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	fmt.Println("Connected!")
	fmt.Printf("  Host:    %s:%d\n", mc.Info.Host, mc.Info.Port)
	fmt.Printf("  RSA:     %v\n", mc.Info.RSAUsed)
}
