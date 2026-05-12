// 00_rsa_connect demonstrates connecting to a remote FutuOpenD with RSA encryption.
//
// RSA encryption is required when connecting to OpenD across networks.
// The client needs OpenD's RSA PUBLIC KEY (in PEM format).
//
// How it works:
//   - OpenD is configured with an RSA key pair (public + private key)
//   - Client obtains OpenD's PUBLIC KEY (PEM format)
//   - Client calls New(WithRSAPublicKey(pem)) then Connect(addr)
//   - The InitConnect packet is encrypted with the public key
//   - OpenD decrypts with its private key, then AES session key is exchanged
//   - All subsequent communication is AES-encrypted
//
// Key format: Go SDK expects "PUBLIC KEY" (PKIX/PKCS#8) PEM — NOT "RSA PRIVATE KEY".
// If you have the private key PEM, convert like:
//
//	openssl rsa -in /etc/futu/keys/private_key.pem -pubout -out /tmp/opend_pubkey.pem
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shing1211/futuapi4go/client"
)

func main() {
	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "172.18.208.88:11111" // remote gateway — update as needed
	}

	pubKeyPEM := os.Getenv("FUTU_RSA_PUBKEY")
	if pubKeyPEM == "" {
		// Default: same path as Python SDK — accepts PKIX "PUBLIC KEY",
		// "RSA PRIVATE KEY" (PKCS1), or "PRIVATE KEY" (PKCS8) PEM.
		// All three formats work; Go extracts the public key automatically.
		data, err := os.ReadFile("/etc/futu/keys/private_key.pem")
		if err != nil {
			fmt.Println("NOTE: Set FUTU_RSA_PUBKEY=/path/to/key.pem")
			fmt.Println()
			fmt.Println("Go SDK accepts three PEM formats:")
			fmt.Println("  - \"-----BEGIN PUBLIC KEY-----\"     (PKIX, recommended)")
			fmt.Println("  - \"-----BEGIN RSA PRIVATE KEY-----\" (PKCS1)")
			fmt.Println("  - \"-----BEGIN PRIVATE KEY-----\"     (PKCS8)")
			fmt.Println()
			connectWithoutRSA(addr)
			return
		}
		pubKeyPEM = strings.TrimSpace(string(data))
	}

	cli := client.New(client.WithRSAPublicKey(pubKeyPEM))
	defer func() {
		cli.Close()
		fmt.Println("Disconnected.")
	}()

	fmt.Printf("Connecting to %s with RSA encryption...\n", addr)
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("RSA connection failed: %v", err)
	}

	fmt.Println("Connected!")
	fmt.Printf("  ConnID:     %d\n", cli.GetConnID())
	fmt.Printf("  ServerVer:  %d\n", cli.GetServerVer())
	fmt.Printf("  Encrypt:    %v\n", cli.IsEncrypt())
	fmt.Printf("  LoginUID:   %d\n", cli.GetLoginUserID())

	// Verify trade context works with RSA
	fmt.Println("\nVerifying trade context...")
	if cli.CanSendProto(2206) { // PID 2206 = GetAccList
		fmt.Println("  Trade proto available ✓")
	} else {
		fmt.Println("  Trade proto not available")
	}
}

// connectWithoutRSA demonstrates what happens when RSA is not used.
// Trade APIs will reject connections to remote OpenD with encryption enabled.
func connectWithoutRSA(addr string) {
	cli := client.New()
	defer cli.Close()

	fmt.Printf("Connecting to %s (no RSA)...\n", addr)
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	fmt.Println("Connected (no encryption).")
	fmt.Printf("  Encrypt: %v\n", cli.IsEncrypt())
	fmt.Println("\nWARNING: Trade APIs will fail on remote OpenD without RSA.")
	fmt.Println("Expected error: 'In order to ensure the security of trade, ")
	fmt.Println("cross-network communication, trade connect need to be encrypted.'")
}
