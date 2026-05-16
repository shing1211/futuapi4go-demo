package main

import (
	"fmt"
	"time"

	"github.com/shing1211/futuapi4go/pkg/tracing/otel"
)

func main() {
	meter, err := otel.NewOTelMeter()
	if err != nil {
		fmt.Printf("meter error: %v\n", err)
		return
	}
	defer meter.Close()

	meter.RecordConnection("tcp")
	time.Sleep(50 * time.Millisecond)

	meter.RecordReconnect("tcp")
	time.Sleep(50 * time.Millisecond)

	meter.RecordAPICall("3001", "success", 45*time.Millisecond)
	meter.RecordAPICall("3004", "success", 120*time.Millisecond)
	meter.RecordAPICall("3001", "error", 0)

	meter.RecordPushMessage("3001")
	meter.RecordPushMessage("3004")

	meter.RecordAPIError("3001", "-1")

	meter.RecordRateLimited("3001")

	meter.RecordRetry("3001", "1")
	meter.RecordRetry("3001", "2")

	meter.RecordBreakerState("main", 0)
	meter.RecordBreakerState("main", 0.5)
	meter.RecordBreakerState("main", 1)

	meter.RecordOpenDUp(true)
	time.Sleep(50 * time.Millisecond)
	meter.RecordOpenDUp(false)

	fmt.Println("All OTel metrics recorded successfully")
}
