# Troubleshooting

Consolidated troubleshooting for the `futuapi4go-demo` project. If you hit a
problem running an example, check here before opening an issue.

---

## Common Errors

| Error | Cause / Fix |
|-------|-------------|
| `connection refused` | OpenD is not running. Set `FUTU_ADDR=127.0.0.1:11111` (or your OpenD address) and start OpenD first. |
| no data from `GetKLines`, `GetQuote`, etc. | Call `client.Subscribe` first — most market-data APIs require a subscription for push-type data (see example `47_subscribe_quote`). |
| `账户购买力不足` | Simulate account has no buying power — expected in sim mode. Use a real account or fund the sim account. |
| `模拟交易不支持` / `模拟账户不支持` | Function not supported in simulate mode — switch to real trading with `WithTradeEnv(constant.TrdEnv_Real)` and `FUTU_TRADE_PWD` set. |
| `未知的协议ID` | OpenD doesn't implement this API (e.g. `ReconfirmOrder`). |
| `没有解锁交易，请先解锁交易` | Need to unlock trading — set `FUTU_TRADE_PWD` (MD5 of your trade password, 32 hex chars). |
| `请求获取实时K线接口前，请先订阅` | Must subscribe to the K-line type before calling `GetKLines` (see `48_subscribe_kline_single`). |
| `暂不提供美股 OTC 市场行情` | Some US stocks are OTC and unsupported — skip with error handling. |

---

## Simulate Trading Limitations

These APIs are **not supported** in simulate trading mode:

| Example | Function | Error |
|---------|----------|-------|
| `43_order_fill` | GetOrderFillList | 模拟交易不支持成交数据 |
| `44_history_fill` | GetHistoryOrderFillList | 模拟交易不支持成交数据 |
| `56_order_fee` | GetOrderFee | 暂时不支持模拟交易 |
| `57_margin_ratio` | GetMarginRatio | 模拟账户不支持 |
| `58_flow_summary` | GetFlowSummary | 模拟账户不支持查询现金流水 |
| `64_reconfirm_order` | ReconfirmOrder | 未知的协议ID (OpenD doesn't implement) |

**Workaround:** use real trading environment (`WithTradeEnv(1)`) with
`FUTU_TRADE_PWD` set. Note that **order operations must never be retried** —
`PlaceOrder` / `ModifyOrder` / `CancelOrder` are non-idempotent.

---

## Known SDK Issues

### GetDelayStatistics — proto2 packed encoding (SDK fix v0.5.13+)

Some OpenD C++ parsers reject Go's default packed encoding for `repeated int32`
fields. The SDK includes a custom proto2 marshaling workaround in
`pkg/sys/system.go` (`marshalC2SProto2`).

**Demo behavior:** Example `96` calls `GetDelayStatistics` and handles both
success and failure gracefully. If OpenD rejects the request (older build or
missing API support), the demo prints a clear explanation and exits cleanly.

### GetTradeDate — all C2S fields are required

`GetTradeDate` has all required fields in its C2S. If the SDK doesn't populate
all required fields, OpenD returns `解析protobuf协议失败`. Works correctly with
OpenD v10.5.6508. The demo exits with a red error if this API fails.

---

## Known Caveats

- **US stocks** require `client.Subscribe` before `GetQuote` returns data; HK stocks do not.
- **Futures accounts** — use `GetAccList(TrdCategory_Future)`, not `GetAccountList` which only returns stock/options accounts.
- **`secMarket` required** — `PlaceOrder` and `GetMaxTrdQtys` require an explicit `TrdSecMarket` parameter.
- **Real trading** requires `FUTU_TRADE_PWD` (MD5 hash of the trading password, 32 hex chars) and `WithTradeEnv(constant.TrdEnv_Real)`.
- **RSA connections** — remote OpenD hosts default to RSA (`FUTU_OPEND_HOSTS` or a non-localhost `FUTU_ADDR`). Localhost defaults to plaintext. See `00_rsa_connect`.
- **Disk-constrained builds** — set `GOTMPDIR` to a location with free space if `/tmp` fills up during `go build`/`go test`.

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD server address | `127.0.0.1:11111` |
| `FUTU_TRADE_PWD` | MD5 hash of trading password (32 hex chars) | (not set) |
| `FUTU_RSA_PUBKEY` / `FUTU_RSA_KEY` | RSA public key PEM for remote encrypted connections | `/etc/futu/keys/private_key.pem` |
| `FUTU_OPEND_HOSTS` | Comma-separated `host:port:isRSA` list for HA failover | (empty → `FUTU_ADDR`) |
| `FUTU_TCP_TIMEOUT` | TCP probe timeout in seconds | `3` |
| `FUTU_WS_ADDR` / `FUTU_WS_SECRET` | WebSocket OpenD address / secret | (not set; WebSocket support pending) |

---

## Proto Reference

- Official docs: https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md
- OpenD downloads: https://www.futunn.com/download/fetch-lasted-link?name=opend-windows
