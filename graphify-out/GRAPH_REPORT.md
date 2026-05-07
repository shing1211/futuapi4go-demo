# Graph Report - .  (2026-05-07)

## Corpus Check
- 97 files · ~22,946 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 287 nodes · 178 edges · 128 communities (90 shown, 38 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 18 edges (avg confidence: 0.77)
- Token cost: 17,000 input · 6,500 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Project Setup & Trading|Project Setup & Trading]]
- [[_COMMUNITY_SDK Core & Protocols|SDK Core & Protocols]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_K-Line Data|K-Line Data]]
- [[_COMMUNITY_Error Handling|Error Handling]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 90|Community 90]]
- [[_COMMUNITY_Community 91|Community 91]]
- [[_COMMUNITY_Community 92|Community 92]]
- [[_COMMUNITY_Community 93|Community 93]]
- [[_COMMUNITY_Community 94|Community 94]]
- [[_COMMUNITY_Community 95|Community 95]]
- [[_COMMUNITY_Community 96|Community 96]]
- [[_COMMUNITY_Community 97|Community 97]]
- [[_COMMUNITY_Community 98|Community 98]]
- [[_COMMUNITY_Community 99|Community 99]]
- [[_COMMUNITY_Community 100|Community 100]]
- [[_COMMUNITY_Community 101|Community 101]]
- [[_COMMUNITY_Community 102|Community 102]]
- [[_COMMUNITY_Community 103|Community 103]]
- [[_COMMUNITY_Community 104|Community 104]]
- [[_COMMUNITY_Community 105|Community 105]]
- [[_COMMUNITY_Community 106|Community 106]]
- [[_COMMUNITY_Community 107|Community 107]]
- [[_COMMUNITY_Community 108|Community 108]]
- [[_COMMUNITY_Community 109|Community 109]]
- [[_COMMUNITY_Community 110|Community 110]]
- [[_COMMUNITY_Community 111|Community 111]]
- [[_COMMUNITY_Community 112|Community 112]]
- [[_COMMUNITY_Community 113|Community 113]]
- [[_COMMUNITY_Community 114|Community 114]]
- [[_COMMUNITY_Community 115|Community 115]]
- [[_COMMUNITY_Community 116|Community 116]]
- [[_COMMUNITY_Community 117|Community 117]]
- [[_COMMUNITY_Community 118|Community 118]]
- [[_COMMUNITY_Community 119|Community 119]]
- [[_COMMUNITY_Community 120|Community 120]]
- [[_COMMUNITY_Community 121|Community 121]]
- [[_COMMUNITY_Community 122|Community 122]]
- [[_COMMUNITY_Community 123|Community 123]]
- [[_COMMUNITY_Community 124|Community 124]]
- [[_COMMUNITY_Community 125|Community 125]]
- [[_COMMUNITY_Community 126|Community 126]]
- [[_COMMUNITY_Community 127|Community 127]]

## God Nodes (most connected - your core abstractions)
1. `futuapi4go-demo Project` - 12 edges
2. `futuapi4go SDK` - 6 edges
3. `main()` - 4 edges
4. `QOT_SUB` - 4 edges
5. `main()` - 3 edges
6. `main()` - 3 edges
7. `main()` - 3 edges
8. `main()` - 3 edges
9. `main()` - 3 edges
10. `main()` - 3 edges

## Surprising Connections (you probably didn't know these)
- `futuapi4go-demo Project` --references--> `Implementation Plan`  [EXTRACTED]
  AGENTS.md → docs/IMPLEMENTATION_PLAN.md
- `futuapi4go SDK` --has_version--> `SDK Version v0.5.4`  [EXTRACTED]
  AGENTS.md → README.md
- `futuapi4go SDK` --implements--> `System APIs (pkg/sys)`  [EXTRACTED]
  AGENTS.md → docs/FUTU_PROTO_REF.md
- `futuapi4go SDK` --implements--> `Trading APIs (pkg/trd)`  [EXTRACTED]
  AGENTS.md → docs/FUTU_PROTO_REF.md
- `futuapi4go SDK` --uses--> `WebSocket Connection Example`  [EXTRACTED]
  AGENTS.md → examples/00_ws_connect/README.md

## Hyperedges (group relationships)
- **Trading System Components** — trading_mode_simulate, trading_mode_real, env_futu_trade_pwd, opend_server, simulator [EXTRACTED 1.00]
- **Push Subscription Examples** — example_02_ticker, example_03_orderbook, example_07_kline_multi [EXTRACTED 1.00]
- **Advanced Combo Examples** — example_76_pre_trade_checks, example_78_dca_grid_bot, example_79_momentum_scanner, example_80_vwap_executor [EXTRACTED 1.00]
- **Proto Known Issues** — issue_proto2_wire_format, issue_get_trade_date, proto_get_delay_statistics, proto_get_trade_date [EXTRACTED 1.00]
- **SDK API Packages** — proto_api_sys, proto_api_qot, proto_api_trd [EXTRACTED 1.00]

## Communities (128 total, 38 thin omitted)

### Community 0 - "Project Setup & Trading"
Cohesion: 0.11
Nodes (19): Environment Variable FUTU_ADDR, Environment Variable FUTU_TRADE_PWD, Example 00: client.Connect, Example 01: client.GetQuote, Example 02: chanpkg.SubscribeTicker, Example 03: chanpkg.SubscribeOrderBook, Example 22: client.PlaceOrder, Example 67: Order Lifecycle (+11 more)

### Community 1 - "SDK Core & Protocols"
Cohesion: 0.22
Nodes (9): Connection Pool O(1) Lookup, KLType, Market Data APIs (pkg/qot), System APIs (pkg/sys), Trading APIs (pkg/trd), QotMarket, futuapi4go SDK, SDK Version v0.5.4 (+1 more)

### Community 2 - "Community 2"
Cohesion: 0.25
Nodes (8): TRD_PLACE_ORDER, OpenD, OpenQuoteContext, OpenSecTradeContext, place_order, FUTU_API_PYTHON_REFERENCE, Trading, TrdEnv

### Community 3 - "Community 3"
Cohesion: 0.4
Nodes (6): InitConnect, KeepAlive, AES Encryption, PacketEncAlgo, FUTU_API_PROTO_REFERENCE, RSA Encryption

### Community 4 - "Community 4"
Cohesion: 0.4
Nodes (6): QOT_SUB, get_stock_quote, Json, ProtoFmt, Protobuf, subscribe

### Community 5 - "Community 5"
Cohesion: 0.7
Nodes (4): main(), marketStateString(), ptrInt32(), ptrStr()

### Community 6 - "Community 6"
Cohesion: 0.6
Nodes (4): main(), printPrices(), updatePrice(), PriceData

### Community 7 - "Community 7"
Cohesion: 0.83
Nodes (3): main(), ptrInt32(), ptrStr()

### Community 8 - "Community 8"
Cohesion: 0.83
Nodes (3): main(), ptrInt32(), ptrStr()

### Community 9 - "K-Line Data"
Cohesion: 0.83
Nodes (3): main(), ptrInt32(), ptrStr()

### Community 10 - "Error Handling"
Cohesion: 0.83
Nodes (3): main(), ptrInt32(), ptrStr()

### Community 11 - "Community 11"
Cohesion: 0.83
Nodes (3): checkMarketReadiness(), main(), marketStateString()

### Community 12 - "Community 12"
Cohesion: 0.67
Nodes (3): Example 07: chanpkg.SubscribeKLines, Example 15: client.RequestHistoryKL, Historical Data Downloader

### Community 13 - "Community 13"
Cohesion: 0.67
Nodes (3): Context with Timeout, Graceful Shutdown, Structured Error Handling

### Community 14 - "Community 14"
Cohesion: 0.67
Nodes (3): Version 0.5.1, Version 0.5.2, Unreleased Version

## Knowledge Gaps
- **80 isolated node(s):** `PriceData`, `SDK Version v0.5.4`, `Environment Variable FUTU_ADDR`, `OpenD Simulator`, `Simulate Trading Mode` (+75 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **38 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `futuapi4go-demo Project` connect `Project Setup & Trading` to `SDK Core & Protocols`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **Why does `futuapi4go SDK` connect `SDK Core & Protocols` to `Project Setup & Trading`?**
  _High betweenness centrality (0.004) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `QOT_SUB` (e.g. with `subscribe` and `get_stock_quote`) actually correct?**
  _`QOT_SUB` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `PriceData`, `SDK Version v0.5.4`, `Environment Variable FUTU_ADDR` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Project Setup & Trading` be split into smaller, more focused modules?**
  _Cohesion score 0.11 - nodes in this community are weakly interconnected._