package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	user, err := client.GetUserInfo(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("GetUserInfo failed: %v", err)
	}
	fmt.Printf("USER: id=%d nickname=%s avatar=%s apiLevel=%s\n",
		user.UserID, user.NickName, user.AvatarUrl, user.ApiLevel)

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(user)
}
