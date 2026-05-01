# Examples

> Each `main.go` demonstrates one SDK function. Run with `go run ./examples/NN_name`.

For full documentation, see the [root README.md](../README.md).

## Basic Function Examples (00-65)

| # | Example | SDK Function |
|---|---------|-------------|
| 00 | [`00_connect`](./00_connect) | `client.Connect` |
| 01 | [`01_quote`](./01_quote) | `client.GetQuote` |
| 02 | [`02_ticker`](./02_ticker) | `chanpkg.SubscribeTicker` |
| 03 | [`03_orderbook`](./03_orderbook) | `chanpkg.SubscribeOrderBook` |
| 04 | [`04_rt`](./04_rt) | `chanpkg.SubscribeRT` |
| 05 | [`05_broker`](./05_broker) | `chanpkg.SubscribeBroker` |
| 06 | [`06_kline_single`](./06_kline_single) | `client.GetKLines` |
| 07 | [`07_kline_multi`](./07_kline_multi) | `chanpkg.SubscribeKLines` |
| 08 | [`08_orderbook_req`](./08_orderbook_req) | `client.GetOrderBook` |
| 09 | [`09_ticker_req`](./09_ticker_req) | `client.GetTicker` |
| 10 | [`10_rt_req`](./10_rt_req) | `client.GetRT` |
| 11 | [`11_broker_req`](./11_broker_req) | `client.GetBroker` |
| 12 | [`12_capital_flow`](./12_capital_flow) | `client.GetCapitalFlow` |
| 13 | [`13_plate_set`](./13_plate_set) | `client.GetPlateSet` |
| 14 | [`14_plate_stock`](./14_plate_stock) | `client.GetPlateSecurity` |
| 15 | [`15_history_kline`](./15_history_kline) | `client.RequestHistoryKL` |
| 16 | [`16_market_state`](./16_market_state) | `client.GetMarketState` |
| 17 | [`17_global_state`](./17_global_state) | `client.GetGlobalState` |
| 18 | [`18_account_list`](./18_account_list) | `client.GetAccountList` |
| 19 | [`19_account_list`](./19_account_list) | `client.GetAccountInfo` |
| 20 | [`20_funds`](./20_funds) | `client.GetFunds` |
| 21 | [`21_positions`](./21_positions) | `client.GetPositionList` |
| 22 | [`22_place_order`](./22_place_order) | `client.PlaceOrder` |
| 23 | [`23_order_list`](./23_order_list) | `client.GetOrderList` |
| 24 | [`24_snapshot`](./24_snapshot) | `client.GetSecuritySnapshot` |
| 25 | [`25_trade_date`](./25_trade_date) | `client.GetTradeDate` |
| 26 | [`26_price_reminder`](./26_price_reminder) | `client.GetPriceReminder` |
| 27 | [`27_cancel_order`](./27_cancel_order) | `client.ModifyOrder` (cancel) |
| 28 | [`28_owner_plate`](./28_owner_plate) | `client.GetOwnerPlate` |
| 29 | [`29_capital_distribution`](./29_capital_distribution) | `client.GetCapitalDistribution` |
| 30 | [`30_stock_filter`](./30_stock_filter) | `client.StockFilter` |
| 31 | [`31_ipo_list`](./31_ipo_list) | `client.GetIpoList` |
| 32 | [`32_future_info`](./32_future_info) | `client.GetFutureInfo` |
| 33 | [`33_suspend`](./33_suspend) | `client.GetSuspend` |
| 34 | [`34_holding_change`](./34_holding_change) | `client.GetHoldingChangeList` |
| 35 | [`35_rehab`](./35_rehab) | `client.RequestRehab` |
| 36 | [`36_code_change`](./36_code_change) | `client.GetCodeChange` |
| 37 | [`37_warrant`](./37_warrant) | `client.GetWarrant` |
| 38 | [`38_option_chain`](./38_option_chain) | `client.GetOptionChain` |
| 39 | [`39_option_expiration`](./39_option_expiration) | `client.GetOptionExpirationDate` |
| 40 | [`40_reference`](./40_reference) | `client.GetReference` |
| 41 | [`41_user_security`](./41_user_security) | `client.GetUserSecurityGroup` |
| 42 | [`42_history_order`](./42_history_order) | `client.GetHistoryOrderList` |
| 43 | [`43_order_fill`](./43_order_fill) | `client.GetOrderFillList` |
| 44 | [`44_history_fill`](./44_history_fill) | `client.GetHistoryOrderFillList` |
| 45 | [`45_acc_trading_info`](./45_acc_trading_info) | `client.GetAccTradingInfo` |
| 46 | [`46_user_info`](./46_user_info) | `client.GetUserInfo` |
| 47 | [`47_subscribe_quote`](./47_subscribe_quote) | `chanpkg.SubscribeQuote` |
| 48 | [`48_subscribe_kline_single`](./48_subscribe_kline_single) | `chanpkg.SubscribeKLine` |
| 49 | [`49_subscribe_price_reminder`](./49_subscribe_price_reminder) | `chanpkg.SubscribePriceReminder` |
| 50 | [`50_unsubscribe`](./50_unsubscribe) | `client.Unsubscribe` |
| 51 | [`51_unsubscribe_all`](./51_unsubscribe_all) | `client.UnsubscribeAll` |
| 52 | [`52_query_subscription`](./52_query_subscription) | `client.QuerySubscription` |
| 53 | [`53_reg_qot_push`](./53_reg_qot_push) | `client.RegQotPush` |
| 54 | [`54_cancel_all_order`](./54_cancel_all_order) | `client.CancelAllOrder` |
| 55 | [`55_max_trd_qtys`](./55_max_trd_qtys) | `client.GetMaxTrdQtys` |
| 56 | [`56_order_fee`](./56_order_fee) | `client.GetOrderFee` |
| 57 | [`57_margin_ratio`](./57_margin_ratio) | `client.GetMarginRatio` |
| 58 | [`58_flow_summary`](./58_flow_summary) | `client.GetFlowSummary` |
| 59 | [`59_static_info`](./59_static_info) | `client.GetStaticInfo` |
| 60 | [`60_modify_user_security`](./60_modify_user_security) | `client.ModifyUserSecurity` |
| 61 | [`61_sub_info`](./61_sub_info) | `client.GetSubInfo` |
| 62 | [`62_set_price_reminder`](./62_set_price_reminder) | `client.SetPriceReminder` |
| 63 | [`63_sub_acc_push`](./63_sub_acc_push) | `client.SubAccPush` |
| 64 | [`64_reconfirm_order`](./64_reconfirm_order) | `client.ReconfirmOrder` |
| 65 | [`65_history_kl_quota`](./65_history_kl_quota) | `client.RequestHistoryKLQuota` |

## Gap Fill Examples (66-69)

| # | Example | SDK Function |
|---|---------|-------------|
| 66 | [`66_multi_symbol_kline`](./66_multi_symbol_kline) | `client.GetKLines` + `RequestHistoryKL` (batch) |
| 67 | [`67_order_lifecycle`](./67_order_lifecycle) | `PlaceOrder` → `GetOrderList` → `ModifyOrder` |
| 68 | [`68_market_hours_check`](./68_market_hours_check) | `GetMarketState` + `GetTradeDate` |
| 69 | [`69_subscribe_handler`](./69_subscribe_handler) | `Subscribe` + push handlers |

## Futures & Options Examples (70-75)

| # | Example | SDK Function |
|---|---------|-------------|
| 70 | [`70_futures_account_list`](./70_futures_account_list) | `cli.Trade().GetAccList(TrdCategory_Future)` |
| 71 | [`71_futures_cash`](./71_futures_cash) | `client.GetAccTradingInfo` (futures margin) |
| 72 | [`72_futures_positions`](./72_futures_positions) | `cli.Trade().GetPositionList(TrdMarket_Futures)` |
| 73 | [`73_options_account_list`](./73_options_account_list) | `cli.Trade().GetAccList` + options rights check |
| 74 | [`74_options_cash`](./74_options_cash) | `client.GetAccountInfo` + `GetFunds` for options |
| 75 | [`75_options_positions`](./75_options_positions) | `client.GetPositionList` (stocks + options) |

## Advanced Combo Examples (76-80)

| # | Example | Strategy |
|---|---------|----------|
| 76 | [`76_pre_trade_checks`](./76_pre_trade_checks) | Market + Account + Position validation |
| 77 | [`77_realtime_dashboard`](./77_realtime_dashboard) | Real-time monitoring with subscriptions |
| 78 | [`78_dca_grid_bot`](./78_dca_grid_bot) | Dollar Cost Averaging + Grid strategy |
| 79 | [`79_momentum_scanner`](./79_momentum_scanner) | StockFilter + Snapshot + K-lines |
| 80 | [`80_vwap_executor`](./80_vwap_executor) | OrderBook + VWAP calculation + execution |

**80 examples total** — all SDK functions covered except `client.GetDelayStatistics` (known OpenD proto bug).