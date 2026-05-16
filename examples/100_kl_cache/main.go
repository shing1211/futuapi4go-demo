package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/cache"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
)

func main() {
	cli := client.New()
	defer cli.Close()

	if err := cli.Connect("127.0.0.1:11111"); err != nil {
		log.Fatalf("connect: %v", err)
	}

	klCache := cache.NewKLCache(
		cache.WithMaxEntries(100),
		cache.WithTTL(5*time.Minute),
	)

	security := &qotcommon.Security{
		Market: int32Ptr(int32(constant.Market_HK)),
		Code:   strPtr("00700"),
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	klCache.StartCleanup(1*time.Minute, stopCh)

	start := time.Now()
	klines1, ok := klCache.Get(security, int32(constant.KLType_K_Day), 0)
	if ok {
		fmt.Printf("Cache HIT: %d KLines (retrieved in %v)\n", len(klines1), time.Since(start))
	} else {
		fmt.Printf("Cache MISS: first call, need to fetch...\n")
		klCache.Set(security, int32(constant.KLType_K_Day), 0, nil)
	}

	start = time.Now()
	klines2, ok := klCache.Get(security, int32(constant.KLType_K_Day), 0)
	if ok {
		fmt.Printf("Cache HIT: %d KLines (retrieved in %v)\n", len(klines2), time.Since(start))
	} else {
		fmt.Println("Cache MISS (expected if Set stored nil)")
	}

	fmt.Printf("Cache size: %d entries\n", klCache.Size())
	klCache.Clear()
	fmt.Printf("After Clear - cache size: %d entries\n", klCache.Size())
}

func int32Ptr(v int32) *int32 { return &v }
func strPtr(s string) *string { return &s }
