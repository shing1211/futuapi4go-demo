package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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

	groups, err := client.GetUserSecurityGroup(context.Background(), cli)
	if err != nil {
		log.Fatalf("GetUserSecurityGroup failed: %v", err)
	}
	for _, g := range groups {
		fmt.Printf("GROUP: name=%s type=%d\n", g.Name, g.GroupType)
	}

	// List stocks in first group
	if len(groups) > 0 {
		stocks, err := client.GetUserSecurity(context.Background(), cli, groups[0].Name)
		if err != nil {
			log.Fatalf("GetUserSecurity failed: %v", err)
		}
		for _, s := range stocks {
			fmt.Printf("  STOCK: code=%s name=%s\n", s.Code, s.Name)
		}
	}
}
