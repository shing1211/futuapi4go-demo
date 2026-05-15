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
	"github.com/shing1211/futuapi4go/pkg/pb/common"
	"github.com/shing1211/futuapi4go/pkg/pb/getglobalstate"
)

func marketState(s int32) string {
	switch s {
	case 0:
		return "N/A"
	case 1:
		return "PreMarketBegin"
	case 2:
		return "Morning"
	case 3:
		return "Rest"
	case 4:
		return "Afternoon"
	case 5:
		return "Close"
	case 6:
		return "AfterHours"
	case 7:
		return "PreMarketEnd"
	case 8:
		return "PreMarketMarket"
	case 9:
		return "PreMarketRest"
	case 10:
		return "PreMarketAfter"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	fmt.Println("Connected with RSA!")
	fmt.Printf("  Host:      %s:%d\n", mc.Info.Host, mc.Info.Port)
	fmt.Printf("  RSA Used:  %v\n", mc.Info.RSAUsed)

	// Print GetGlobalState details
	c2s := &getglobalstate.C2S{
		UserID: func() *uint64 { v := uint64(0); return &v }(),
	}
	req := &getglobalstate.Request{C2S: c2s}
	var rsp getglobalstate.Response
	if err := mc.Client.Inner().RequestContext(context.Background(), 1002, req, &rsp); err != nil {
		fmt.Printf("  GetGlobalState error: %v\n", err)
	} else if rsp.GetRetType() != int32(common.RetType_RetType_Succeed) {
		fmt.Printf("  GetGlobalState failed: retType=%d, msg=%s\n", rsp.GetRetType(), rsp.GetRetMsg())
	} else {
		s2c := rsp.GetS2C()
		fmt.Println("  Global State:")
		fmt.Printf("    ServerVer:    %d\n", s2c.GetServerVer())
		fmt.Printf("    QotLogined:   %v\n", s2c.GetQotLogined())
		fmt.Printf("    TrdLogined:   %v\n", s2c.GetTrdLogined())
		fmt.Printf("    Market HK:    %s\n", marketState(s2c.GetMarketHK()))
		fmt.Printf("    Market US:    %s\n", marketState(s2c.GetMarketUS()))
		fmt.Printf("    Market SH:    %s\n", marketState(s2c.GetMarketSH()))
		fmt.Printf("    Market SZ:    %s\n", marketState(s2c.GetMarketSZ()))
		fmt.Printf("    Market HKF:   %s\n", marketState(s2c.GetMarketHKFuture()))
		fmt.Printf("    Market USF:   %s\n", marketState(s2c.GetMarketUSFuture()))
		fmt.Printf("    Trade supported: %v\n", s2c.GetTrdLogined())
	}

	if mc.Client.CanSendProto(2206) {
		fmt.Println("  Trade proto available ✓")
	}
}