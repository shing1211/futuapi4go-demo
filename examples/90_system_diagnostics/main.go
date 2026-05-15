package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== OpenD Connection Diagnostic Demo ===")
	fmt.Println("Comprehensive diagnostic of the OpenD connection and protocol support")
	fmt.Println()

	fmt.Println("--- Connection Info ---")
	state, err := client.GetGlobalState(ctx, mc.Client)
	if err != nil {
		fmt.Printf("GetGlobalState: %v\n", err)
	} else {
		fmt.Printf("  Server Build:  v%d (build %d)\n", state.ServerVer, state.ServerBuildNo)
		fmt.Printf("  Server Time:   %d\n", state.Time)
		fmt.Printf("  Local Time:    %.0f\n", state.LocalTime)
		fmt.Printf("  Qot Logined:   %v\n", state.QotLogined)
		fmt.Printf("  Trd Logined:   %v\n", state.TrdLogined)
		fmt.Printf("  Program:       %s (status=%d)\n", state.ProgramStatusDesc, state.ProgramStatus)

	}
	fmt.Printf("  Client Ver:    %d\n", mc.Client.GetServerVer())
	fmt.Printf("  Login User ID: %d\n", mc.Client.GetLoginUserID())
	fmt.Printf("  Conn ID:       %d\n", mc.Client.GetConnID())

	fmt.Println()
	fmt.Println("--- Protocol Support Check ---")
	keyProtos := []struct {
		id    uint32
		name  string
	}{
		{constant.ProtoID_Qot_GetBasicQot, "GetQuote"},
		{constant.ProtoID_Qot_GetKL, "GetKLines"},
		{constant.ProtoID_Qot_GetOptionChain, "GetOptionChain"},
		{constant.ProtoID_Trd_PlaceOrder, "PlaceOrder"},
		{constant.ProtoID_Trd_SubAccPush, "SubAccPush"},
		{constant.ProtoID_Trd_UpdateOrder, "PushOrderUpdate (2208)"},
		{constant.ProtoID_Trd_UpdateOrderFill, "PushOrderFill (2218)"},
		{constant.ProtoID_Qot_UpdatePriceReminder, "PushPriceReminder (3019)"},
	}

	supported := 0
	unsupported := 0
	for _, p := range keyProtos {
		if mc.Client.CanSendProto(p.id) {
			fmt.Printf("  ✓ %s\n", p.name)
			supported++
		} else {
			fmt.Printf("  ✗ %s — NOT supported by this OpenD\n", p.name)
			unsupported++
		}
	}
	fmt.Printf("  Result: %d supported, %d not supported\n", supported, unsupported)

	fmt.Println()
	fmt.Println("--- Test Command ---")
	testResult, err := client.TestCmd(ctx, mc.Client, "echo", "futuapi4go-demo diagnostic")
	if err != nil {
		fmt.Printf("  TestCmd: %v\n", err)
	} else {
		fmt.Printf("  Cmd:    %s\n", testResult.Cmd)
		fmt.Printf("  Result: %s\n", testResult.Result)
	}

	fmt.Println()
	fmt.Println("--- User Info ---")
	userInfo, err := client.GetUserInfo(ctx, mc.Client)
	if err != nil {
		fmt.Printf("GetUserInfo: %v\n", err)
	} else {
		fmt.Printf("  User ID:     %d\n", userInfo.UserID)
	}

	accounts, err := client.GetAccountList(ctx, mc.Client)
	if err != nil {
		fmt.Printf("GetAccountList: %v\n", err)
	} else {
		fmt.Printf("\n  Accounts: %d\n", len(accounts))
		for _, acc := range accounts {
			fmt.Printf("    AccID=%d type=%d env=%d markets=%v\n",
				acc.AccID, acc.AccType, acc.TrdEnv, acc.TrdMarketAuthList)
		}
	}

	fmt.Println()
	fmt.Println("--- Market State ---")
	usState, err := client.GetMarketState(ctx, mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		fmt.Printf("GetMarketState (US): %v\n", err)
	} else {
		fmt.Printf("  US Market:  state=%d (%s)\n", usState, marketStateStr(usState))
	}

	hkState, err := client.GetMarketState(ctx, mc.Client, constant.Market_HK, "00700")
	if err != nil {
		fmt.Printf("GetMarketState (HK): %v\n", err)
	} else {
		fmt.Printf("  HK Market:  state=%d (%s)\n", hkState, marketStateStr(hkState))
	}

	fmt.Println()
	fmt.Println("--- Connection Type ---")
	fmt.Printf("  Client: %T\n", mc.Client)
	if mc.Info != nil {
		fmt.Printf("  Connected to: %s:%d\n", mc.Info.Host, mc.Info.Port)
		fmt.Printf("  RSA Used: %v\n", mc.Info.RSAUsed)
	} else {
		fmt.Println("  Connected to: (no info)")
	}
	fmt.Printf("  State: %v\n", mc.State)

	fmt.Println()
	fmt.Println("=== Diagnostic Complete ===")
}

func marketStateStr(state int32) string {
	switch state {
	case 0:
		return "Pre-Market"
	case 1:
		return "Trading"
	case 2:
		return "Post-Market"
	case 3:
		return "Closed"
	default:
		return fmt.Sprintf("Extended(%d)", state)
	}
}
