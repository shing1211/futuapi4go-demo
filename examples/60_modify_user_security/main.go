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

	// Get existing groups first
	groups, err := client.GetUserSecurityGroup(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("GetUserSecurityGroup failed: %v", err)
	}

	// Add a stock to the first group (op=1 means add)
	if len(groups) > 0 {
		if err := client.ModifyUserSecurity(context.Background(), mc.Client,
			groups[0].Name,
			1, // op: 1=Add
			constant.Market_US,
			[]string{"NVDA"},
		); err != nil {
			log.Fatalf("ModifyUserSecurity failed: %v", err)
		}
		fmt.Printf("Added NVDA to group '%s'.\n", groups[0].Name)
	} else {
		// Create a new group and add stock (op=3 means add group)
		if err := client.ModifyUserSecurity(context.Background(), mc.Client,
			"MyWatchlist",
			3, // op: 3=Add group
			constant.Market_US,
			[]string{"NVDA"},
		); err != nil {
			log.Fatalf("ModifyUserSecurity (create group) failed: %v", err)
		}
		fmt.Println("Created group 'MyWatchlist' with NVDA.")
	}
}
