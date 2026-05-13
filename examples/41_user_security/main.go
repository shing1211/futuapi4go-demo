// 41_user_security demonstrates the user watchlist (security group) API.
//
// GetUserSecurityGroup lists the user's custom watchlists.
// GetUserSecurity retrieves the stocks in a given group.
//
// NOTE: The watchlist API may not be available for all account types.
// If "Unsupported watchlist group type" is returned, the account
// may not have watchlist functionality enabled.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	groups, err := client.GetUserSecurityGroup(context.Background(), mc.Client)
	if err != nil {
		fe, ok := constant.AsFutuError(err)
		if ok && fe.Code == -1 {
			fmt.Println("ℹ️  Watchlist (GetUserSecurityGroup) is not available for this account.")
			fmt.Println("   This feature may require a specific Futu account type or region.")
			fmt.Println("   Continuing with other operations...")
			return
		}
		log.Fatalf("GetUserSecurityGroup failed: %v", err)
	}

	if len(groups) == 0 {
		fmt.Println("No watchlist groups found.")
		fmt.Println("Create watchlists through the Futu app and they will appear here.")
		return
	}

	fmt.Printf("Found %d watchlist group(s):\n\n", len(groups))

	for _, g := range groups {
		groupTypeStr := "Custom"
		switch g.GroupType {
		case 1:
			groupTypeStr = "System"
		case 2:
			groupTypeStr = "Custom"
		case 3:
			groupTypeStr = "Smart"
		}
		fmt.Printf("📁 %s (type=%d)\n", g.Name, groupTypeStr)

		stocks, err := client.GetUserSecurity(context.Background(), mc.Client, g.Name)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to load stocks: %v\n", err)
			continue
		}
		if len(stocks) == 0 {
			fmt.Println("   (empty)")
			continue
		}

		for _, s := range stocks {
			fmt.Printf("   • %s — %s\n", s.Code, s.Name)
		}
		fmt.Println()
	}
}