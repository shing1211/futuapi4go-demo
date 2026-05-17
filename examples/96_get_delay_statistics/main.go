package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== GetDelayStatistics Demo ===")
	fmt.Println("Retrieves performance delay statistics for Qot push,")
	fmt.Println("request-reply, and order placement operations.")
	fmt.Println()

	fmt.Println("--- SDK Fix Analysis ---")
	fmt.Println("The SDK (v0.5.13+) includes a custom proto2 marshaling")
	fmt.Println("workaround for GetDelayStatistics. The C2S request is")
	fmt.Println("encoded using non-packed repeated int32 fields (proto2")
	fmt.Println("format) instead of Go's default packed encoding (proto3).")
	fmt.Println()
	fmt.Println("This resolves the '解析protobuf协议失败' error that occurred")
	fmt.Println("with OpenD C++ parsers expecting proto2 wire format.")
	fmt.Println()

	fmt.Println("--- Calling GetDelayStatistics ---")
	stats, err := client.GetDelayStatistics(ctx, mc.Client)
	if err != nil {
		fe, ok := constant.AsFutuError(err)
		if ok {
			fmt.Printf("FutuError: code=%d msg=%s\n", fe.Code, fe.Message)
		}
		fmt.Printf("GetDelayStatistics: %v\n", err)
		fmt.Println()
		fmt.Println("The API may not be available on this OpenD version.")
		fmt.Println("Common causes:")
		fmt.Println("  1. OpenD version < 10.5 (older builds lack this API)")
		fmt.Println("  2. ProtoID mismatch between SDK and OpenD")
		fmt.Println("  3. API requires specific server-side configuration")
		fmt.Println()
		fmt.Println("This is expected — the API is optional and not all")
		fmt.Println("OpenD deployments support it. All other SDK functions")
		fmt.Println("continue to work normally.")
		return
	}

	fmt.Println("✓ GetDelayStatistics succeeded!")
	fmt.Println()

	fmt.Println("--- Qot Push Delay Statistics ---")
	fmt.Printf("  Push Type: %d | Avg Delay: %.2fms | Count: %d\n",
		stats.QotPushType, stats.DelayAvg, stats.Count)
	fmt.Println()
	fmt.Printf("  %-8s %-8s %-8s %-12s %s\n",
		"Begin", "End", "Count", "Proportion", "Cumulative")
	fmt.Println("  " + repeat("─", 50))
	for _, item := range stats.ItemList {
		fmt.Printf("  %-8d %-8d %-8d %-10.2f%% %-10.2f%%\n",
			item.Begin, item.End, item.Count, item.Proportion, item.CumulativeRatio)
	}

	fmt.Println()
	fmt.Println("--- Request-Reply Statistics ---")
	if len(stats.ReqReplyList) == 0 {
		fmt.Println("  (no data available yet)")
	} else {
		fmt.Printf("  %-8s %-6s %-12s %-12s %-12s %s\n",
			"ProtoID", "Count", "TotalCost", "OpenDCost", "NetDelay", "Local?")
		fmt.Println("  " + repeat("─", 60))
		for _, r := range stats.ReqReplyList {
			fmt.Printf("  %-8d %-6d %-12.2f %-12.2f %-12.2f %v\n",
				r.ProtoID, r.Count, r.TotalCostAvg, r.OpenDCostAvg, r.NetDelayAvg, r.IsLocalReply)
		}
	}

	fmt.Println()
	fmt.Println("--- Place Order Statistics ---")
	if len(stats.PlaceOrderList) == 0 {
		fmt.Println("  (no data available yet)")
	} else {
		fmt.Printf("  %-12s %-12s %-12s %-12s %s\n",
			"OrderID", "TotalCost", "OpenDCost", "NetDelay", "UpdateCost")
		fmt.Println("  " + repeat("─", 60))
		for _, p := range stats.PlaceOrderList {
			fmt.Printf("  %-12s %-12.2f %-12.2f %-12.2f %.2f\n",
				p.OrderID, p.TotalCost, p.OpenDCost, p.NetDelay, p.UpdateCost)
		}
	}

	fmt.Println()
	fmt.Println("=== Delay Statistics Complete ===")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(stats)
}

func repeat(s string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = s[0]
	}
	return string(b)
}
