# Option / API Coverage Plan — Tiers 2–5

This document tracks the remaining SDK proto APIs in the `futuapi4go` SDK that
do not yet have a corresponding demo example in the `futuapi4go-demo` repo.
Tier 1 (Option family) was completed in the `106–110` examples; see
`docs/EXAMPLES.md`.

Each row is one example the demo will add. Numbers continue sequentially
from `111` to keep numbering uniform.

## Progress checklist

- [x] Tier 1 — Option family (`106_option_volatility_rank`, `107_option_screener`, `108_option_risk`, `109_option_events`, `110_option_indicators`)
- [ ] Tier 2 — `111_institutional_flow`, `112_shareholder_data`, `113_top_brokers`, `114_short_data`
- [ ] Tier 3 — `115_macro_calendar`, `116_macro_indicators`, `117_earnings_dividends`, `118_research_ratings`, `119_ranks_distributions`, `120_search_discovery`, `121_user_security_single`
- [ ] Tier 4 — `122_kl_realtime`, `123_history_kl_variants`, `124_rehab_history`, `125_combo_order`, `126_quick_trade`, `127_today_orders_fills`, `128_system_verification`
- [ ] Tier 5 — `129_stock_screen`

---

## Tier 2 — Institutional / shareholder / flow suites

Goal: document each coherent data family in a single compound example so a
new developer can see the full surface at a glance.

| # | Path | SDK funcs | Notes |
|---|------|-----------|-------|
| 111 | `examples/111_institutional_flow` | `GetInstitutionList` + `GetInstitutionProfile` + `GetInstitutionHoldingList` + `GetInstitutionHoldingChange` + `GetInstitutionDistribution` | HK/US institutional-investor profiles, holdings, and inflow/outflow tracking |
| 112 | `examples/112_shareholder_data` | `GetShareholdersHolderDetail` + `GetShareholdersHoldingChanges` + `GetShareholdersInstitutional` + `GetShareholdersOverview` | Shareholder composition: top holders, holder change log, institutional breakdown, overview |
| 113 | `examples/113_top_brokers` | `GetTopTenBuySellBrokers` | Standalone — single API; today's top-10 broker net buy/sell per market |
| 114 | `examples/114_short_data` | `GetDailyShortVolume` + `GetShortInterest` + `GetShortSellingRank` | Short-volume / short-interest series + ranking |

---

## Tier 3 — Macro / research / ranks (mostly standalone)

| # | Path | SDK funcs | Notes |
|---|------|-----------|-------|
| 115 | `examples/115_macro_calendar` | `GetEconomicCalendar` + `GetFedWatchDotPlot` + `GetFedWatchTargetRate` | Economic release calendar + Fed funds rate path |
| 116 | `examples/116_macro_indicators` | `GetMacroIndicatorList` + `GetMacroIndicatorHistory` | Macro indicators (PMI, CPI, GDP, …) list and history series |
| 117 | `examples/117_earnings_dividends` | `GetEarningsCalendar` + `GetEarningsBeatRank` + `GetDividendCalendar` + `GetDividendRank` + `GetHighDividendSOERank` | Earnings / dividend calendars and rankings |
| 118 | `examples/118_research_ratings` | `GetResearchAnalystConsensus` + `GetResearchMorningstarReport` + `GetResearchRatingSummary` + `GetRatingChange` | Analyst ratings, Morningstar report, rating-change feed |
| 119 | `examples/119_ranks_distributions` | `GetHotList` + `GetHeatMapData` + `GetRiseFallDistribution` + `GetPeriodChangeRank` | Hot-list / heat-map / rise-fall distribution / period change rank |
| 120 | `examples/120_search_discovery` | `GetSearchQuote` + `GetSearchNews` | Symbol search and news search |
| 121 | `examples/121_user_security_single` | `GetUserSecurity` | vs `GetUserSecurityGroup` (#41) — this one returns individual watch-list items |

---

## Tier 4 — K-line / rehab variants & trading extras

| # | Path | SDK funcs | Notes |
|---|------|-----------|-------|
| 122 | `examples/122_kl_realtime` | `GetKL` | Distinct ProtoID `qotgetkl` (not `qotgetkline` covered by `GetKLines` #06) |
| 123 | `examples/123_history_kl_variants` | `GetHistoryKL` + `GetHistoryKLPoints` + `RequestHistoryKLQuota` | vs `RequestHistoryKL` (#15) — alternate KL-pull / points / quota APIs |
| 124 | `examples/124_rehab_history` | `GetRehab` | vs `RequestRehab` (#35) — distinct ProtoID `qotgetrehab` |
| 125 | `examples/125_combo_order` | `PlaceComboOrder` + `GetComboMaxTrdQtys` | RFQ workflow from `97e` but here actually calling placement on simulate |
| 126 | `examples/126_quick_trade` | `QuickBuy` + `QuickSell` + `QuickMarketBuy` + `QuickMarketSell` | One-click trading convenience wrappers (place market/limit in a single call) |
| 127 | `examples/127_today_orders_fills` | `GetTodayOrders` + `GetTodayFills` | vs `GetOrderList` / `GetOrderFillList` (#23/#43) — today-only filter variants |
| 128 | `examples/128_system_verification` | `Verification` (sys) | Standalone unlock-verification round-trip |

---

## Tier 5 — Screening ambiguity (1 item)

| # | Path | SDK funcs | Notes |
|---|------|-----------|-------|
| 129 | `examples/129_stock_screen` | `StockScreen` | `qotstockscreen` (distinct from `qotstockfilter` covered by `#30`) |

---

## Estimated work

- ~19 new examples (`111–129`).
- Each follows the established demo template:
  - `connect.MustConnect(ctx)` + `defer mc.Close()`
  - pointer helpers (`ptrInt32`, `ptrStr`, `ptrFloat64`)
  - `display.PrintJSON(rsp)` at the end
  - nil guards on all slice iteration
- All require new pb-package imports for the `*C2S` request types — the
  helpers from `97a`–`97e` and `106`–`110` already cover the import pattern.
- `go vet ./...` + `go build ./...` on the new dirs (with `GOTMPDIR` set on
  disk-constrained systems).
- Update `docs/EXAMPLES.md`:
  - Header count: `113 → 132`.
  - New section `Institutional & Macro Coverage (111–129)` (or grouped
    subsections matching the tiers above).
- No SDK code changes are required — every function listed here is already
  implemented in `github.com/shing1211/futuapi4go` and exposed via the
  `client/` convenience wrappers or the pkg-level `qot`/`trd`/`sys` funcs.

## Risks / open questions

- Some macro / research / institutional suites may have very low data
  density outside of HK / specific exchanges. Pick safe underlyings (NVDA
  for US, `00700` for HK) and document `RetType ≠ Succeed` paths.
- `Verification` (`128`) requires account credentials — may need a
  `ConnectWithEnv` switch that some users don't set up. Consider gating
  with an env-var check.
- `StockScreen` (`129`) is distinct from `StockFilter` (#30) but visually
  identical — doc must highlight the difference.