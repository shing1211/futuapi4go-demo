# Example Reference

> All 114 standalone examples (plus 5 option sub-examples `97a`–`97e`) in the [futuapi4go-demo](https://github.com/shing1211/futuapi4go-demo) repository.
> Each example is a `main.go` in its own directory under `examples/`.
> Run with: `go run ./examples/NN_name`

## Connection Examples

| # | Example | SDK Function | Description |
|---|---------|-------------|-------------|
| 00 | [`00_connect`](../examples/00_connect) | `client.Connect` | Plain TCP connection (local OpenD) |
| 00 | [`00_rsa_connect`](../examples/00_rsa_connect) | `client.Connect` + RSA | RSA-encrypted connection (remote OpenD) |
| 00 | [`00_ws_connect`](../examples/00_ws_connect) | `client.ConnectWS` | WebSocket connection |

## Basic Function Examples (01–64)

| # | Example | SDK Function |
|---|---------|-------------|
| 01 | [`01_quote`](../examples/01_quote) | `client.GetQuote` |
| 02 | [`02_ticker`](../examples/02_ticker) | `chanpkg.SubscribeTicker` |
| 03 | [`03_orderbook`](../examples/03_orderbook) | `chanpkg.SubscribeOrderBook` |
| 04 | [`04_rt`](../examples/04_rt) | `chanpkg.SubscribeRT` |
| 05 | [`05_broker`](../examples/05_broker) | `chanpkg.SubscribeBroker` |
| 06 | [`06_kline_single`](../examples/06_kline_single) | `client.GetKLines` |
| 07 | [`07_kline_multi`](../examples/07_kline_multi) | `chanpkg.SubscribeKLines` |
| 08 | [`08_orderbook_req`](../examples/08_orderbook_req) | `client.GetOrderBook` |
| 09 | [`09_ticker_req`](../examples/09_ticker_req) | `client.GetTicker` |
| 10 | [`10_rt_req`](../examples/10_rt_req) | `client.GetRT` |
| 11 | [`11_broker_req`](../examples/11_broker_req) | `client.GetBroker` |
| 12 | [`12_capital_flow`](../examples/12_capital_flow) | `client.GetCapitalFlow` |
| 13 | [`13_plate_set`](../examples/13_plate_set) | `client.GetPlateSet` |
| 14 | [`14_plate_stock`](../examples/14_plate_stock) | `client.GetPlateSecurity` |
| 15 | [`15_history_kline`](../examples/15_history_kline) | `client.RequestHistoryKL` |
| 16 | [`16_market_state`](../examples/16_market_state) | `client.GetMarketState` |
| 17 | [`17_global_state`](../examples/17_global_state) | `client.GetGlobalState` |
| 18 | [`18_account_list`](../examples/18_account_list) | `client.GetAccountList` |
| 19 | [`19_account_list`](../examples/19_account_list) | `client.GetAccountInfo` |
| 20 | [`20_funds`](../examples/20_funds) | `client.GetFunds` |
| 21 | [`21_positions`](../examples/21_positions) | `client.GetPositionList` |
| 22 | [`22_place_order`](../examples/22_place_order) | `client.PlaceOrder` |
| 23 | [`23_order_list`](../examples/23_order_list) | `client.GetOrderList` |
| 24 | [`24_snapshot`](../examples/24_snapshot) | `client.GetSecuritySnapshot` |
| 25 | [`25_trade_date`](../examples/25_trade_date) | `client.GetTradeDate` |
| 26 | [`26_price_reminder`](../examples/26_price_reminder) | `client.GetPriceReminder` |
| 27 | [`27_cancel_order`](../examples/27_cancel_order) | `client.ModifyOrder` (cancel) |
| 28 | [`28_owner_plate`](../examples/28_owner_plate) | `client.GetOwnerPlate` |
| 29 | [`29_capital_distribution`](../examples/29_capital_distribution) | `client.GetCapitalDistribution` |
| 30 | [`30_stock_filter`](../examples/30_stock_filter) | `client.StockFilter` |
| 31 | [`31_ipo_list`](../examples/31_ipo_list) | `client.GetIpoList` |
| 32 | [`32_future_info`](../examples/32_future_info) | `client.GetFutureInfo` |
| 33 | [`33_suspend`](../examples/33_suspend) | `client.GetSuspend` |
| 34 | [`34_holding_change`](../examples/34_holding_change) | `client.GetHoldingChangeList` |
| 35 | [`35_rehab`](../examples/35_rehab) | `client.RequestRehab` |
| 36 | [`36_code_change`](../examples/36_code_change) | `client.GetCodeChange` |
| 37 | [`37_warrant`](../examples/37_warrant) | `client.GetWarrant` |
| 38 | [`38_option_chain`](../examples/38_option_chain) | `client.GetOptionChain` |
| 39 | [`39_option_expiration`](../examples/39_option_expiration) | `client.GetOptionExpirationDate` |
| 40 | [`40_reference`](../examples/40_reference) | `client.GetReference` |
| 41 | [`41_user_security`](../examples/41_user_security) | `client.GetUserSecurityGroup` |
| 42 | [`42_history_order`](../examples/42_history_order) | `client.GetHistoryOrderList` |
| 43 | [`43_order_fill`](../examples/43_order_fill) | `client.GetOrderFillList` |
| 44 | [`44_history_fill`](../examples/44_history_fill) | `client.GetHistoryOrderFillList` |
| 45 | [`45_acc_trading_info`](../examples/45_acc_trading_info) | `client.GetAccTradingInfo` |
| 46 | [`46_user_info`](../examples/46_user_info) | `client.GetUserInfo` |
| 47 | [`47_subscribe_quote`](../examples/47_subscribe_quote) | `chanpkg.SubscribeQuote` |
| 48 | [`48_subscribe_kline_single`](../examples/48_subscribe_kline_single) | `chanpkg.SubscribeKLine` |
| 49 | [`49_subscribe_price_reminder`](../examples/49_subscribe_price_reminder) | `chanpkg.SubscribePriceReminder` |
| 50 | [`50_unsubscribe`](../examples/50_unsubscribe) | `client.Unsubscribe` |
| 51 | [`51_unsubscribe_all`](../examples/51_unsubscribe_all) | `client.UnsubscribeAll` |
| 52 | [`52_query_subscription`](../examples/52_query_subscription) | `client.QuerySubscription` |
| 53 | [`53_reg_qot_push`](../examples/53_reg_qot_push) | `client.RegQotPush` |
| 54 | [`54_cancel_all_order`](../examples/54_cancel_all_order) | `client.CancelAllOrder` |
| 55 | [`55_max_trd_qtys`](../examples/55_max_trd_qtys) | `client.GetMaxTrdQtys` |
| 56 | [`56_order_fee`](../examples/56_order_fee) | `client.GetOrderFee` |
| 57 | [`57_margin_ratio`](../examples/57_margin_ratio) | `client.GetMarginRatio` |
| 58 | [`58_flow_summary`](../examples/58_flow_summary) | `client.GetFlowSummary` |
| 59 | [`59_static_info`](../examples/59_static_info) | `client.GetStaticInfo` |
| 60 | [`60_modify_user_security`](../examples/60_modify_user_security) | `client.ModifyUserSecurity` |
| 61 | [`61_sub_info`](../examples/61_sub_info) | `client.GetSubInfo` |
| 62 | [`62_set_price_reminder`](../examples/62_set_price_reminder) | `client.SetPriceReminder` |
| 63 | [`63_sub_acc_push`](../examples/63_sub_acc_push) | `client.SubAccPush` |
| 64 | [`64_reconfirm_order`](../examples/64_reconfirm_order) | `client.ReconfirmOrder` |

## Smoke Test

| # | Example | SDK Functions |
|---|---------|---------------|
| 65 | [`65_smoke_test`](../examples/65_smoke_test) | `GetGlobalState` + `GetAccountList` + `GetQuote` | Fast happy-path check for CI (exits nonzero on failure) |

## Gap Fill Examples (66–69)

| # | Example | SDK Functions |
|---|---------|---------------|
| 66 | [`66_multi_symbol_kline`](../examples/66_multi_symbol_kline) | `Subscribe` + `GetKLines` + `RequestHistoryKL` |
| 67 | [`67_order_lifecycle`](../examples/67_order_lifecycle) | `PlaceOrder` → `GetOrderList` → `ModifyOrder` |
| 68 | [`68_market_hours_check`](../examples/68_market_hours_check) | `GetMarketState` + `GetTradeDate` |
| 69 | [`69_subscribe_handler`](../examples/69_subscribe_handler) | `Subscribe` + push handlers (Ticker/KLine/OrderBook) |

## Futures & Options Examples (70–75)

| # | Example | Description |
|---|---------|-------------|
| 70 | [`70_futures_account_list`](../examples/70_futures_account_list) | `GetAccList(TrdCategory_Future)` for futures accounts |
| 71 | [`71_futures_cash`](../examples/71_futures_cash) | Futures margin and cash queries |
| 72 | [`72_futures_positions`](../examples/72_futures_positions) | `GetPositionList(TrdMarket_Futures)` |
| 73 | [`73_options_account_list`](../examples/73_options_account_list) | Options rights check via `GetAccList` |
| 74 | [`74_options_cash`](../examples/74_options_cash) | Options buying power and margin |
| 75 | [`75_options_positions`](../examples/75_options_positions) | Stock + options combined positions |

## Advanced Combo Examples (76–80)

| # | Example | Description |
|---|---------|-------------|
| 76 | [`76_pre_trade_checks`](../examples/76_pre_trade_checks) | Market state + account funds + position validation |
| 77 | [`77_realtime_dashboard`](../examples/77_realtime_dashboard) | Real-time monitoring with ticker subscriptions |
| 78 | [`78_dca_grid_bot`](../examples/78_dca_grid_bot) | Dollar Cost Averaging + Grid strategy |
| 79 | [`79_momentum_scanner`](../examples/79_momentum_scanner) | StockFilter + Snapshot + K-lines momentum analysis |
| 80 | [`80_vwap_executor`](../examples/80_vwap_executor) | OrderBook + VWAP calculation + execution |

## Advanced Trading Examples (81–85)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 81 | [`81_options_trading`](../examples/81_options_trading) | `GetOptionChain` + `GetQuote` + `GetFunds` | Options chain analysis with covered call workflow |
| 82 | [`82_trailing_stop`](../examples/82_trailing_stop) | `PlaceOrder` (TrailingStop) + `GetOrderList` | Trailing stop order with builder & parameter reference |
| 83 | [`83_trade_push_monitor`](../examples/83_trade_push_monitor) | `SubAccPush` + `SetPushHandler` + `ParsePushOrderUpdate` + `ParsePushOrderFill` | Real-time trade push event monitoring |
| 84 | [`84_order_builder`](../examples/84_order_builder) | `TradeAPI.NewOrder` (fluent `OrderBuilder`) | Fluent OrderBuilder: 3 patterns (limit, market, GTC) |
| 85 | [`85_risk_analyzer`](../examples/85_risk_analyzer) | `GetFunds` (PDT fields) + `GetMarginRatio` + `GetMaxTrdQtys` | Portfolio risk dashboard, PDT compliance, margin analysis |

## Infrastructure & Diagnostics Examples (86–90)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 86 | [`86_history_downloader`](../examples/86_history_downloader) | `history.NewDownloader` + `DownloadWithStats` + `NewConcurrentDownloader` + `DownloadMultiple` | Bulk historical data with pagination, retry, progress tracking |
| 87 | [`87_option_tools`](../examples/87_option_tools) | `option.ParseCode` + `option.FindAtm` + `option.FilterByExpiry` + `option.StrikeDistance` | Option code parsing, ATM finding, strike filtering utilities |
| 88 | [`88_convenience_trading`](../examples/88_convenience_trading) | `client.PlaceOrder` + `client.GetFunds` + `client.GetPositionList` | Simplest possible one-liner trading operations |
| 89 | [`89_quota_manager`](../examples/89_quota_manager) | `SystemAPI.GetUsedQuota` + `SubscribeSymbols` + `Unsubscribe` + `UnsubscribeAll` | Subscription quota management |
| 90 | [`90_system_diagnostics`](../examples/90_system_diagnostics) | `GetGlobalState` + `CanSendProto` + `TestCmd` + `GetConnID` + `GetMarketState` | Comprehensive OpenD connection diagnostic |

## Quantitative Trading Strategies (91–95)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 91 | [`91_orderbook_imbalance`](../examples/91_orderbook_imbalance) | `SubscribeOrderBook` + `GetOrderBook` + `OrderBookDetail` (iceberg detection) | Real-time bid/ask imbalance, liquidity pressure scoring, iceberg/spoof detection |
| 92 | [`92_pairs_trading`](../examples/92_pairs_trading) | `GetKLines` (60d) + `GetQuote` + Pearson correlation + spread z-score | Statistical arbitrage: correlation, spread normalization, mean-reversion signals |
| 93 | [`93_smart_money`](../examples/93_smart_money) | `GetCapitalFlow` + `GetCapitalDistribution` + `GetBroker` + `GetOrderBook` | Institutional accumulation/distribution score from 4 fused data sources |
| 94 | [`94_portfolio_rebalance`](../examples/94_portfolio_rebalance) | `GetPositionList` + `GetFunds` + `GetQuote` + `GetMaxTrdQtys` + `PlaceOrder` | Multi-asset portfolio rebalancer with PDT-aware drift correction |
| 95 | [`95_earnings_vol_strategy`](../examples/95_earnings_vol_strategy) | `GetOptionExpirationDate` + `GetOptionChain` + `GetKLines` (120d) | Earnings straddle: implied move vs historical move, vol strategy builder |

## Special Cases (96)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 96 | [`96_get_delay_statistics`](../examples/96_get_delay_statistics) | `client.GetDelayStatistics` | Performance delay statistics with graceful error handling |

## Infrastructure, Tracing & Advanced Features (97–101)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 97 | [`97_opentelemetry_tracing`](../examples/97_opentelemetry_tracing) | `tracing.SetTracer` + `otel.NewTracer` | OpenTelemetry tracing setup: stdout exporter, TracerProvider, auto-generated spans |
| 97a | [`97a_option_quote`](../examples/97a_option_quote) | `client.GetOptionQuote` | Greeks (Δ/Γ/Θ/ν/ρ), bid/ask, OI for a multi-leg option quote |
| 97b | [`97b_option_strategy`](../examples/97b_option_strategy) | `client.GetOptionStrategy` | Strategy catalogue (straddle / strangle / spread …) for an underlying |
| 97c | [`97c_strategy_analysis`](../examples/97c_strategy_analysis) | `client.GetOptionStrategyAnalysis` | Per-strategy P&L analysis: max profit/loss, breakevens, prob of profit |
| 97d | [`97d_strategy_spread`](../examples/97d_strategy_spread) | `client.GetOptionStrategySpread` | Spread prices for vertical / ratio option strategies |
| 97e | [`97e_combo_order`](../examples/97e_combo_order) | `client.GetComboMaxTrdQtys` + `client.PlaceComboOrder` | Combo-order RFQ workflow (max qty + placement preview) |
| 98 | [`98_state_shutdown`](../examples/98_state_shutdown) | `WithOnStateChange` + `State` + `Shutdown` | Connection state machine monitoring and graceful shutdown |
| 99 | [`99_audit_validation`](../examples/99_audit_validation) | `ValidateOrder` + `HasErrors` + `NewAuditLogger` | Order pre-flight validation and structured audit logging |
| 100 | [`100_kl_cache`](../examples/100_kl_cache) | `cache.NewKLCache` + `cache.NewKLCachedClient` | K-line LRU+TTL cache: hit/miss behavior, cleanup, size monitoring |
| 101 | [`101_otel_metrics`](../examples/101_otel_metrics) | `otel.NewOTelMeter` | OpenTelemetry metrics: connection, API calls, latency, rate limiting, breaker state |

## Prediction Markets / Event Contract (102–105)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 102 | [`102_ec_explore`](../examples/102_ec_explore) | `FilterCompetition` + `GetEventContractCategory` + `GetEventContractSeriesList` + `GetEventContractEventList` + `GetEventContract` + `GetEventContractMilestoneList` | Discover EC hierarchy: category → series → event → contracts → milestones |
| 103 | [`103_ec_market_data`](../examples/103_ec_market_data) | `GetEventContractSnapshot` + `GetEventContractOrderBook` + `GetEventContractTicker` + `GetEventContractKline` + `RequestHistoryEventContractKL` | Real-time + historical EC quotes (YES/NO bid/ask, predSide, klineSource) |
| 104 | [`104_ec_subscribe`](../examples/104_ec_subscribe) | `SubEventContract` + `RegisterHandler` + `ParseUpdateEventContract{OrderBook,Kline,Ticker}` | Streaming EC push subscription & parsing |
| 105 | [`105_ec_combo`](../examples/105_ec_combo) | `GetEventContractComboList` + `GetEventContractComboRfq` | Combine EC contracts into a Combo and request a firm RFQ quote |

## Option & Screener APIs (106–110)

| # | Example | SDK Functions | Description |
|---|---------|---------------|-------------|
| 106 | [`106_option_volatility_rank`](../examples/106_option_volatility_rank) | `client.GetOptionVolatility` + `client.GetOptionRank` + `client.GetOptionMarketStatistic` | Implied/historical vol, option ranking by volume, market-wide call/put stat |
| 107 | [`107_option_screener`](../examples/107_option_screener) | `client.GetOptionEarningsScreener` + `client.GetOptionZeroDteScreener` + `client.GetOptionSellerScreener` + `client.OptionScreen` | Four screening engines: earnings, 0-DTE, covered-call, v10.6 multi-criteria |
| 108 | [`108_option_risk`](../examples/108_option_risk) | `client.GetOptionExerciseProbability` + `client.GetOptionUnderlyingHisStatistic` + `client.GetOptionUnderlyingHisVolatility` + `client.GetOptionUnderlyingOverview` + `client.GetOptionUnderlyingRank` | Per-strike exercise probability + underlying-level IV/OI/vol stats and ranking |
| 109 | [`109_option_events`](../examples/109_option_events) | `client.GetOptionEvent` + `client.GetOptionEventAlert` + `client.SetOptionEventAlert` | Recent unusual option events; alert CRUD (add / deleteAll) lifecycle |
| 110 | [`110_option_indicators`](../examples/110_option_indicators) | `client.GetIndicatorList` + `client.RequestIndicatorCalc` | Search the MyLang indicator catalogue and compute an indicator on a K-line series |
