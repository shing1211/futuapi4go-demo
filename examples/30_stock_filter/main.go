package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotstockfilter"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	results, err := client.StockFilter(context.Background(), mc.Client, constant.Market_US, 0, 5)
	if err != nil {
		log.Fatalf("StockFilter failed: %v", err)
	}

	for _, r := range results {
		fmt.Printf("STOCK: code=%s name=%s price=%.2f vol=%d\n",
			r.Security.GetCode(), r.Name, r.CurPrice, r.Volume)

		// Print all BaseData fields
		for _, b := range r.BaseDataList {
			if b == nil {
				continue
			}
			fieldName := qotstockfilter.StockField(b.GetFieldName())
			fmt.Printf("  Base:   %s = %.4f\n", fieldName, b.GetValue())
		}

		// Print all AccumulateData fields
		for _, a := range r.AccumulateDataList {
			if a == nil {
				continue
			}
			fieldName := qotstockfilter.AccumulateField(a.GetFieldName())
			fmt.Printf("  Accu:   %s = %.4f (days=%d)\n", fieldName, a.GetValue(), a.GetDays())
		}

		// Print all FinancialData fields
		for _, f := range r.FinancialDataList {
			if f == nil {
				continue
			}
			fieldName := qotstockfilter.FinancialField(f.GetFieldName())
			fmt.Printf("  Fin:    %s = %.4f (quarter=%d)\n", fieldName, f.GetValue(), f.GetQuarter())
		}

		// Print all CustomIndicatorData fields
		for _, c := range r.CustomIndicatorDataList {
			if c == nil {
				continue
			}
			fieldName := qotstockfilter.CustomIndicatorField(c.GetFieldName())
			fmt.Printf("  Custom: %s = %.4f (klType=%d)\n", fieldName, c.GetValue(), c.GetKlType())
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(results)
}
