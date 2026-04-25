package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
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

	// Get existing groups first
	groups, err := client.GetUserSecurityGroup(context.Background(), cli)
	if err != nil {
		log.Fatalf("GetUserSecurityGroup failed: %v", err)
	}

	// Add a stock to the first group (op=1 means add)
	if len(groups) > 0 {
		if err := client.ModifyUserSecurity(context.Background(), cli,
			groups[0].Name,
			1, // op: 1=Add
			int32(constant.Market_US),
			[]string{"NVDA"},
		); err != nil {
			log.Fatalf("ModifyUserSecurity failed: %v", err)
		}
		fmt.Printf("Added NVDA to group '%s'.\n", groups[0].Name)
	} else {
		// Create a new group and add stock (op=3 means add group)
		if err := client.ModifyUserSecurity(context.Background(), cli,
			"MyWatchlist",
			3, // op: 3=Add group
			int32(constant.Market_US),
			[]string{"NVDA"},
		); err != nil {
			log.Fatalf("ModifyUserSecurity (create group) failed: %v", err)
		}
		fmt.Println("Created group 'MyWatchlist' with NVDA.")
	}
}
