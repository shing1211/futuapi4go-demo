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

	user, err := client.GetUserInfo(context.Background(), cli)
	if err != nil {
		log.Fatalf("GetUserInfo failed: %v", err)
	}
	fmt.Printf("USER: id=%d nickname=%s avatar=%s apiLevel=%s\n",
		user.UserID, user.NickName, user.AvatarUrl, user.ApiLevel)
}
