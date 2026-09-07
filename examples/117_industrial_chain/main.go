// 117_industrial_chain demonstrates the industrial chain analysis suite:
//   - GetIndustrialChainList    (list all industrial chains for a market)
//   - GetIndustrialChainDetail  (detail of a specific chain)
//   - GetIndustrialChainByPlate (chains associated with a plate)
//   - GetIndustrialPlateInfo    (info about a specific plate)
//   - GetIndustrialPlateStock   (stocks in a specific plate)
//
// Workflow:
//  1. List all available industrial chains for HK market.
//  2. Get detail for the first chain found.
//  3. Get plate info and constituent stocks for the first plate found.
//
// Note: Industrial chain APIs use numeric IDs (ChainId, PlateId).
// This demo discovers IDs dynamically from the list response.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotgetindustrialchainbyplate "github.com/shing1211/futuapi4go/pkg/pb/qotgetindustrialchainbyplate"
	qotgetindustrialchaindetail "github.com/shing1211/futuapi4go/pkg/pb/qotgetindustrialchaindetail"
	qotgetindustrialchainlist "github.com/shing1211/futuapi4go/pkg/pb/qotgetindustrialchainlist"
	qotgetindustrialplateinfo "github.com/shing1211/futuapi4go/pkg/pb/qotgetindustrialplateinfo"
	qotgetindustrialplatestock "github.com/shing1211/futuapi4go/pkg/pb/qotgetindustrialplatestock"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client
	ctx := context.Background()

	fmt.Println("=== Industrial Chain Analysis ===")

	// 1. GetIndustrialChainList — list all available chains for HK market.
	fmt.Println("\n-- 1) GetIndustrialChainList (HK market) --")
	chainListRsp, err := client.GetIndustrialChainList(ctx, cli, &qotgetindustrialchainlist.C2S{
		Market: ptrInt32(int32(qotcommon.QotMarket_QotMarket_HK_Security)),
	})
	if err != nil {
		log.Printf("GetIndustrialChainList failed: %v", err)
	} else {
		fmt.Printf("Found %d industrial chains\n", len(chainListRsp.DataList))
		for i, chain := range chainListRsp.DataList {
			if chain == nil {
				continue
			}
			fmt.Printf("  [%d] chainId=%d name=%s stocks=%d\n",
				i+1, chain.GetChainId(), chain.GetName(), chain.GetStocksNum())
		}
		display.PrintJSON(chainListRsp)
	}

	// 2. GetIndustrialChainDetail — detail of a specific chain.
	fmt.Println("\n-- 2) GetIndustrialChainDetail --")
	if len(chainListRsp.DataList) > 0 {
		firstChain := chainListRsp.DataList[0]
		chainID := firstChain.GetChainId()
		fmt.Printf("  Fetching detail for chainId=%d (%s)\n", chainID, firstChain.GetName())
		detailRsp, err := client.GetIndustrialChainDetail(ctx, cli, &qotgetindustrialchaindetail.C2S{
			ChainId: ptrInt64(chainID),
		})
		if err != nil {
			fmt.Printf("  GetIndustrialChainDetail: %v\n", err)
		} else {
			fmt.Printf("  Detail received — inspect JSON for plate/stock breakdown\n")
			display.PrintJSON(detailRsp)
		}
	} else {
		fmt.Println("  (skipped — no chains found)")
	}

	// 3. GetIndustrialChainByPlate — chains associated with a plate.
	fmt.Println("\n-- 3) GetIndustrialChainByPlate --")
	if len(chainListRsp.DataList) > 0 && len(chainListRsp.DataList[0].GetRelationSecurityList()) > 0 {
		sec := chainListRsp.DataList[0].GetRelationSecurityList()[0]
		fmt.Printf("  Using first related security: %s\n", sec.GetCode())
		byPlateRsp, err := client.GetIndustrialChainByPlate(ctx, cli, &qotgetindustrialchainbyplate.C2S{
			PlateId: ptrInt64(0), // Would need a valid plate ID
		})
		if err != nil {
			fmt.Printf("  GetIndustrialChainByPlate: %v\n", err)
		} else {
			fmt.Printf("  Found %d chains for plate\n", len(byPlateRsp.RelatedChainList))
			display.PrintJSON(byPlateRsp)
		}
		_ = sec
	} else {
		fmt.Println("  (skipped — no related securities found)")
	}

	// 4. GetIndustrialPlateInfo — info about a specific plate.
	fmt.Println("\n-- 4) GetIndustrialPlateInfo --")
	plateInfoRsp, err := client.GetIndustrialPlateInfo(ctx, cli, &qotgetindustrialplateinfo.C2S{
		PlateId: ptrInt64(0), // Would need a valid plate ID
	})
	if err != nil {
		fmt.Printf("  GetIndustrialPlateInfo: %v\n", err)
	} else {
		fmt.Printf("  Plate info received — inspect JSON for details\n")
		display.PrintJSON(plateInfoRsp)
	}

	// 5. GetIndustrialPlateStock — stocks in a specific plate.
	fmt.Println("\n-- 5) GetIndustrialPlateStock --")
	plateStockRsp, err := client.GetIndustrialPlateStock(ctx, cli, &qotgetindustrialplatestock.C2S{
		PlateId: ptrInt64(0), // Would need a valid plate ID
	})
	if err != nil {
		fmt.Printf("  GetIndustrialPlateStock: %v\n", err)
	} else {
		fmt.Printf("  Plate stocks: %d items\n", len(plateStockRsp.StockList))
		for i, stock := range plateStockRsp.StockList {
			if stock == nil || stock.Security == nil {
				continue
			}
			fmt.Printf("    [%d] code=%s name=%s\n", i+1, stock.Security.GetCode(), stock.GetName())
		}
		display.PrintJSON(plateStockRsp)
	}

	fmt.Println("\n── Done ─────────────────────────────────────────")
}

func ptrInt32(v int32) *int32  { return &v }
func ptrInt64(v int64) *int64  { return &v }
