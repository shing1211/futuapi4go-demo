package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/tracing"
	oteladapter "github.com/shing1211/futuapi4go/pkg/tracing/otel"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	// 1. Create OTel stdout exporter and TracerProvider
	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("futuapi4go-demo"),
		)),
	)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	// 2. Install SDK tracing backend — spans auto-generated from here
	tracing.SetTracer(oteladapter.NewTracer("futuapi4go-demo",
		oteladapter.WithTracerProvider(tp)))

	// 3. Normal SDK usage — each API call and lifecycle event creates spans
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	cli := mc.Client
	if err := client.Subscribe(context.Background(), cli, constant.Market_US, "NVDA",
		[]constant.SubType{constant.SubType_Quote}); err != nil {
		fmt.Printf("Subscribe failed (expected without OpenD): %v\n", err)
		return
	}

	quote, err := client.GetQuote(context.Background(), cli, constant.Market_US, "NVDA")
	if err != nil {
		fmt.Printf("GetQuote failed (expected without OpenD): %v\n", err)
	} else {
		fmt.Printf("NVDA: price=%.2f open=%.2f high=%.2f low=%.2f\n",
			quote.Price, quote.Open, quote.High, quote.Low)
	}

	fmt.Println("\nOpenTelemetry spans exported above (connect, subscribe, getQuote, close)")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(quote)
}
