# futuapi4go-demo — Usage Guide

[English](#english) | [中文](#chinese-中文-简体)

---

## English

### Prerequisites

| Requirement | Details |
|-------------|---------|
| Go | 1.26 or later |
| Futu OpenD | Running locally or remotely (default: `127.0.0.1:11111`) |
| RSA key (optional) | For remote encrypted connections |

Install OpenD:
- **Windows/macOS/Linux**: https://www.futunn.com/download/fetch-lasted-link?name=opend-windows
- Ensure OpenD is running and logged in before running examples.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD TCP address (`host:port`) | `127.0.0.1:11111` |
| `FUTU_WS_ADDR` | OpenD WebSocket address | `127.0.0.1:11117` |
| `FUTU_RSA_PUBKEY` | RSA public key PEM (file path or inline) | (not set) |
| `FUTU_WS_SECRET` | WebSocket secret key | (not set) |
| `FUTU_TRADE_PWD` | MD5 hash of trading password (32 hex chars) | (not set) |
| `FUTU_OPEND_HOSTS` | Comma-separated list of `host:port:RSA` for HA mode | (not set) |

**Windows PowerShell:**
```powershell
$env:FUTU_ADDR = "127.0.0.1:11111"
$env:FUTU_TRADE_PWD = "your-32-char-md5-hex"
```

**Linux/macOS:**
```bash
export FUTU_ADDR=127.0.0.1:11111
export FUTU_TRADE_PWD=your-32-char-md5-hex
```

### Running Examples

```bash
# Run a single example
go run ./examples/01_quote

# Run with verbose output
go run -v ./examples/00_connect

# Build all examples (no-op, verifies compilation)
go build ./...

# Lint
go vet ./...
```

### OpenD Configuration

#### Local (Default)
OpenD defaults to `127.0.0.1:11111` (TCP) and `127.0.0.1:11117` (WebSocket).
No extra configuration needed.

#### Remote Connection with RSA
```powershell
# Set the RSA public key (inline PEM or path to .pem file)
$env:FUTU_RSA_PUBKEY = "-----BEGIN PUBLIC KEY-----(...)

# Point to remote OpenD
$env:FUTU_ADDR = "203.0.113.42:11111"

go run ./examples/00_rsa_connect
```

See `examples/00_rsa_connect/main.go` for the full pattern.

#### HA Failover (Multiple Hosts)
```powershell
# Comma-separated: host:port:RSA (RSA=true/false)
$env:FUTU_OPEND_HOSTS = "203.0.113.42:11111:true,203.0.113.43:11111:true,127.0.0.1:11111:false"

go run ./examples/00_connect_ha
```

The HA client probes all hosts in parallel, connects to the fastest, and auto-reconnects on failure.

#### WebSocket
```powershell
$env:FUTU_WS_ADDR = "127.0.0.1:11117"
$env:FUTU_WS_SECRET = "your-secret-key"

go run ./examples/00_ws_connect
```

### Trading Modes

The SDK defaults to **simulate trading** (`TrdEnv=0`). Many order/flow APIs
are not available in simulate mode.

```go
// Simulate trading (default)
cli := client.New()

// Real trading
cli := client.New().WithTradeEnv(constant.TrdEnv_Real)
```

Real trading requires `FUTU_TRADE_PWD` set to the MD5 hash of your trading password.

**Simulate-only APIs** (will error in real trading mode without proper account):
- Order fills, history fills
- Cash flow / flow summary
- Margin ratio

### Core SDK Patterns

#### Basic Connect
```go
cli := client.New()
if err := cli.Connect("127.0.0.1:11111"); err != nil {
    log.Fatal(err)
}
defer cli.Close()
```

#### Get Quote (HK — no subscribe needed)
```go
quote, err := client.GetQuote(ctx, cli, constant.Market_HK, "00700")
log.Printf("00700: %.2f", quote.Price)
```

#### Get Quote (US — requires subscribe first)
```go
client.Subscribe(ctx, cli, constant.Market_US, "AAPL", []constant.SubType{constant.SubType_Basic})
quote, _ := client.GetQuote(ctx, cli, constant.Market_US, "AAPL")
log.Printf("AAPL: %.2f", quote.Price)
```

#### Streaming via Channels
```go
ch := make(chan *client.PushQuote, 100)
stop, _ := chanpkg.SubscribeQuote(ctx, cli, constant.Market_HK, "00700", ch)
defer stop()
for q := range ch {
    log.Printf("%s: %.2f", q.Code, q.CurPrice)
}
```

#### Place Order
```go
order, _ := trd.NewOrder(accID, constant.TrdMarket_HK, constant.TrdEnv_Simulate).
    Buy("00700", 100).
    At(350.0).
    Build()

resp, err := trd.PlaceOrder(ctx, cli.Inner(), order)
log.Printf("OrderID: %s", resp.OrderID)
```

### Example Categories

| Range | Category | Examples |
|-------|----------|----------|
| `00_*` | Connection | TCP, RSA, WebSocket, HA |
| `01–17` | Market Data | Quote, ticker, orderbook, K-line, flow, plates |
| `18–64` | Trading | Account, positions, orders, subscribe/push |
| `65–69` | Gap Fill | Multi-function workflows |
| `70–75` | Futures/Options | Futures/options accounts, positions, margin |
| `76–80` | Advanced Combo | Pre-trade checks, DCA grid, VWAP, momentum |
| `81–85` | Advanced Trading | Options, trailing stop, risk analysis |
| `86–95` | Quant Strategies | Orderbook imbalance, pairs, smart money |
| `96–101` | Infrastructure | Delay stats, OTel tracing, KL cache, metrics |
| `102–105` | Event Contracts | Prediction market / Moomoo US APIs |
| `97a–e` | Option Strategies | Greeks, spread, combo order |

See [EXAMPLES.md](EXAMPLES.md) for the complete list.

### K-Line Cache (Example 100)
```go
klCache := cache.NewKLCache(
    cache.WithMaxEntries(2000),
    cache.WithTTL(5*time.Minute),
)
cachedCli := cache.NewKLCachedClient(cli.Inner(), klCache)
klines, err := cachedCli.GetKL(ctx, rehabType, klType, security)
```

### Graceful Shutdown
```go
if err := cli.Shutdown(5 * time.Second); err != nil {
    log.Printf("shutdown error: %v", err)
}
```

### Distributed Tracing (Example 97)
```go
import (
    "github.com/shing1211/futuapi4go/pkg/tracing"
    "github.com/shing1211/futuapi4go/pkg/tracing/otel"
)

tracing.SetTracer(otel.NewTracer("my-app"))
// All API calls, connect/disconnect, and push handlers generate spans automatically.
```

See `examples/97_opentelemetry_tracing` for the full example.

---

## Chinese 中文 (简体)

### 前置要求

| 需求 | 详情 |
|------|------|
| Go | 1.26 或更高版本 |
| Futu OpenD | 本地或远程运行（默认：`127.0.0.1:11111`） |
| RSA 密钥（可选） | 远程加密连接使用 |

安装 OpenD：https://www.futunn.com/download/fetch-lasted-link?name=opend-windows

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `FUTU_ADDR` | OpenD TCP 地址（`host:port`） | `127.0.0.1:11111` |
| `FUTU_WS_ADDR` | OpenD WebSocket 地址 | `127.0.0.1:11117` |
| `FUTU_RSA_PUBKEY` | RSA 公钥 PEM（文件路径或内容） | （未设置） |
| `FUTU_WS_SECRET` | WebSocket 密钥 | （未设置） |
| `FUTU_TRADE_PWD` | 交易密码 MD5 哈希值（32 字符十六进制） | （未设置） |

### 运行示例

```bash
go run ./examples/01_quote          # 单个示例
go build ./...                      # 编译所有示例
go vet ./...                        # 代码检查
```

### 交易模式

SDK 默认为**模拟交易**（`TrdEnv=0`）。实盘交易：

```go
cli := client.New().WithTradeEnv(constant.TrdEnv_Real)
```

实盘交易需要设置 `FUTU_TRADE_PWD`（交易密码的 MD5 哈希值）。

### 常见模式

```go
// 基本连接
cli := client.New()
cli.Connect("127.0.0.1:11111")
defer cli.Close()

// 获取报价（港股无需订阅）
quote, _ := client.GetQuote(ctx, cli, constant.Market_HK, "00700")

// 美股需先订阅
client.Subscribe(ctx, cli, constant.Market_US, "AAPL", []constant.SubType{constant.SubType_Basic})
quote, _ := client.GetQuote(ctx, cli, constant.Market_US, "AAPL")

// 订阅推送（通道模式）
ch := make(chan *client.PushQuote, 100)
stop, _ := chanpkg.SubscribeQuote(ctx, cli, constant.Market_HK, "00700", ch)
defer stop()

// 下单
order := trd.NewOrder(accID, constant.TrdMarket_HK, constant.TrdEnv_Simulate).
    Buy("00700", 100).At(350.0).Build()
resp, _ := trd.PlaceOrder(ctx, cli.Inner(), order)
```

详见 [EXAMPLES.md](EXAMPLES.md) 完整示例列表。
