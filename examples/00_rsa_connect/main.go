// 00_rsa_connect demonstrates connecting to a remote FutuOpenD with RSA encryption.
//
// RSA encryption is required when connecting to OpenD across networks.
// Configure per-host RSA settings in FUTU_OPEND_HOSTS:
//
//	FUTU_OPEND_HOSTS=127.0.0.1:11111:false,172.18.208.88:11111:true
//
// The connect package automatically handles:
//   - RSA key detection from FUTU_RSA_KEY
//   - Per-host RSA mode (true/false)
//   - Auto-fallback (tries opposite RSA mode if first fails)
//
// Key format: Go SDK accepts three PEM formats:
//   - "-----BEGIN PUBLIC KEY-----"     (PKIX, recommended)
//   - "-----BEGIN RSA PRIVATE KEY-----" (PKCS1)
//   - "-----BEGIN PRIVATE KEY-----"     (PKCS8)
package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	fmt.Println("Connected with RSA!")
	fmt.Printf("  Host:      %s:%d\n", mc.Info.Host, mc.Info.Port)
	fmt.Printf("  RSA Used:  %v\n", mc.Info.RSAUsed)

	fmt.Println()

	display.PrintAll(context.Background(), mc)
}
