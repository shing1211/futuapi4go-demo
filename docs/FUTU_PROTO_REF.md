# Futu API Proto Reference

> The field-by-field reference for every API in the demo. Full protocol docs: https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md

## Table of Contents

- [Common Types](#common-types)
- [System APIs (pkg/sys)](#system-apis-pkgsys)
- [Market Data APIs (pkg/qot)](#market-data-apis-pkgqot)
- [Market Analysis APIs (pkg/qot)](#market-analysis-apis-pkgqot)
- [Options & Warrants APIs (pkg/qot)](#options--warrants-apis-pkgqot)
- [Historical Data APIs (pkg/qot)](#historical-data-apis-pkgqot)
- [Corporate Actions APIs (pkg/qot)](#corporate-actions-apis-pkgqot)
- [User Data APIs (pkg/qot)](#user-data-apis-pkgqot)
- [Trading APIs (pkg/trd)](#trading-apis-pkgtrd)
- [Push Subscriptions (pkg/qot)](#push-subscriptions-pkgqot)
- [Known Issues](#known-issues)

---

## Common Types

### QotMarket (行情市场)
| Value | Market |
|-------|--------|
| 1 | HK |
| 2 | US |
| 3 | SH (A-share) |
| 4 | SZ (A-share) |
| 7 | HKFuture |
| 8 | USFuture |
| 9 | SG (Singapore) |
| 10 | JP (Japan) |

### SecurityType (证券类型)
| Value | Type |
|-------|------|
| 1 | Stock |
| 2 | Warrant |
| 3 | Drt |
| 4 | Bond |
| 5 | Option |
| 6 | Futures |
| 7 | Index |
| 9 | Plate |

### RehabType (复权类型)
| Value | Type |
|-------|------|
| 0 | None |
| 1 | Forward |
| 2 | Backward |

### KLType (K线类型)
| Value | Type |
|-------|------|
| 1 | 1Min |
| 2 | Day |
| 3 | Week |
| 4 | Month |
| 5 | Year |
| 6 | 1Min |
| 7 | 5Min |
| 8 | 15Min |
| 9 | 30Min |
| 10 | 60Min |

---

## System APIs (pkg/sys)

### GetGlobalState — ProtoID 1002
```protobuf
message C2S {
    required uint64 userID = 1; // 填0
}
```

### GetUserInfo — ProtoID 1005
```protobuf
message C2S {} // 空消息
```

### GetDelayStatistics — ProtoID 1006
> ⚠️ OpenD rejects completely empty C2S with "解析protobuf协议失败"
```protobuf
message C2S {
    repeated int32 typeList = 1;       // 1=QotPush, 2=ReqReply, 3=PlaceOrder
    optional int32 qotPushStage = 2;  // 1=SR2SS, 2=SS2CR, 3=CR2CS, 4=SS2CS, 5=SR2CS
    repeated int32 segmentList = 3;   // 分段，默认100ms以下2ms分段，100ms以上500/1000/2000/-1
}
// 修复: SDK发送 typeList=[1] 以避免空消息
```

---

## Market Data APIs (pkg/qot)

### GetBasicQot — ProtoID 3004
```protobuf
message C2S {
    repeated Security securityList = 1;
}
// 返回: repeated BasicQot
```

### GetKL — ProtoID 3006
```protobuf
message C2S {
    required Security security = 1;
    required int32 rehabType = 2;  // RehabType
    required int32 klType = 3;     // KLType
    required int32 reqNum = 4;
}
```

### GetOrderBook — ProtoID 3012
```protobuf
message C2S {
    required Security security = 1;
    required int32 num = 2; // 档位数量
}
```

### GetTicker — ProtoID 3010
```protobuf
message C2S {
    required Security security = 1;
    required int32 num = 2;
}
```

### GetRT (Intraday) — ProtoID 3008
```protobuf
message C2S {
    required Security security = 1;
}
```

### GetBroker — ProtoID 3014
```protobuf
message C2S {
    required Security security = 1;
    required int32 num = 2;
}
```

### GetSecuritySnapshot — ProtoID 3203
```protobuf
message C2S {
    repeated Security securityList = 1;
}
```

### GetTradeDate — ProtoID 2205
> ⚠️ 所有字段 required，SDK必须全部填充
```protobuf
message C2S {
    required int32 market = 1;     // QotMarket
    required string beginTime = 2; // "2026-01-01"
    required string endTime = 3;  // "2026-12-31"
}
```

---

## Market Analysis APIs (pkg/qot)

### GetPlateSet — ProtoID 3204
```protobuf
message C2S {
    required int32 market = 1; // QotMarket
}
```

### GetPlateSecurity — ProtoID 3205
```protobuf
message C2S {
    required Security plate = 1;
}
```

### GetCapitalFlow — ProtoID 3211
```protobuf
message C2S {
    required Security security = 1;
    required int32 periodType = 2; // 1=Day, 2=Minute
}
```

### GetCapitalDistribution — ProtoID 3212
```protobuf
message C2S {
    required Security security = 1;
}
```

### GetOwnerPlate — ProtoID 3207
```protobuf
message C2S {
    repeated Security securityList = 1;
}
```

### GetReference — ProtoID 3206
```protobuf
message C2S {
    required Security security = 1;
    required int32 referenceType = 2; // 1=Future, 2=Option, 3=Warrant
}
```

### GetStaticInfo — ProtoID 2201 / 3202
```protobuf
message C2S {
    required int32 market = 1; // QotMarket
    optional int32 secType = 2; // SecurityType
    repeated Security securityList = 3;
}
```

### GetFutureInfo — ProtoID 2211 / 3218
```protobuf
message C2S {
    repeated Security securityList = 1;
}
```

### StockFilter — ProtoID 3215
```protobuf
message C2S {
    required int32 market = 1;
    required int32 begin = 2;
    required int32 num = 3;
    repeated BaseFilter baseFilterList = 4;
    // fieldName: 8=VolumeRatio, 9=TurnoverRate, etc.
}
```

---

## Options & Warrants APIs (pkg/qot)

### GetOptionExpirationDate — ProtoID 3224
```protobuf
message C2S {
    required Security owner = 1;       // 标的股票
    optional int32 indexOptionType = 2; // 0=所有, 1=指数期权
}
```

### GetOptionChain — ProtoID 3209
```protobuf
message C2S {
    required Security owner = 1;
    optional int32 indexOptionType = 2;
    optional int32 type = 3;  // 1=Standard, 2=Non-standard
    optional int32 expiryDate = 4;
    optional int32 strikePrice = 5;
    optional int32 optionType = 6; // 1=Call, 2=Put
}
```

### GetWarrant — ProtoID 3210
```protobuf
message C2S {
    required int32 begin = 1;
    required int32 num = 2;
    optional int32 sortField = 3;  // 11=ImpliedVolatility
    optional bool ascend = 4;
    repeated int32 typeList = 5; // 1=Call, 2=Put, 3=Bull, 4=Bear
    optional Security owner = 6;
}
```

---

## Historical Data APIs (pkg/qot)

### RequestHistoryKL — ProtoID 3103
> 支持自动翻页
```protobuf
message C2S {
    required Security security = 1;
    required int32 rehabType = 2;
    required int32 klType = 3;
    required string beginTime = 4;  // "2026-01-01"
    required string endTime = 5;    // "" 表示最新
    required int64 maxAckKLNum = 6;
    optional bytes nextReqKey = 7;  // 分页key
    optional bool isNext = 8;       // true=下一页
}
```

### GetHistoryKL — ProtoID 3225
> 单次时间范围查询
```protobuf
message C2S {
    required Security security = 1;
    required int32 rehabType = 2;
    required int32 klType = 3;
    required string beginTime = 4;
    required string endTime = 5;
}
```

### RequestHistoryKLQuota — ProtoID 3104
```protobuf
message C2S {
    optional bool getDetail = 1;
}
```

### GetRehab — ProtoID 2207 / 3208
```protobuf
message C2S {
    repeated Security securityList = 1;
}
```

---

## Corporate Actions APIs (pkg/qot)

### GetIpoList — ProtoID 2212 / 3217
```protobuf
message C2S {
    required int32 market = 1; // QotMarket
}
```

### GetCodeChange — ProtoID 2210
```protobuf
message C2S {
    repeated Security securityList = 1;
}
```

### GetSuspend — ProtoID 2209
```protobuf
message C2S {
    repeated Security securityList = 1;
    optional string beginTime = 2;
    optional string endTime = 3;
}
```

### GetHoldingChangeList — ProtoID 2213
```protobuf
message C2S {
    required Security security = 1;
    optional int32 holderCategory = 2; // 1=机构
    optional string beginTime = 3;
    optional string endTime = 4;
}
```

---

## User Data APIs (pkg/qot)

### GetUserSecurityGroup — ProtoID 2402 / 3222
```protobuf
message C2S {
    optional int32 groupType = 1; // 1=自选
}
```

### GetUserSecurity — ProtoID 2401 / 3213
```protobuf
message C2S {
    required string groupName = 1;
}
```

### SetPriceReminder — ProtoID 2405 / 3220
```protobuf
message C2S {
    required Security security = 1;
    required int32 op = 2;      // 1=Add, 2=Del, 3=Edit
    required int32 type = 3;    // 1=高于, 2=低于, 3=提醒
    required double value = 4;
    optional string note = 5;
}
```

### GetPriceReminder — ProtoID 2404 / 3221
```protobuf
// Security: required Security
// market: required int32 (QotMarket)
```

---

## Trading APIs (pkg/trd)

### GetAccList — ProtoID 2001
```protobuf
message C2S {
    optional int32 trdMarket = 1;  // TrdMarket filter
}
```

### GetFunds — ProtoID 2002
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    optional int32 currency = 3;  // 1=HKD, 2=USD, 3=CAD
}
```

### GetPositionList — ProtoID 2003
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    optional string code = 3;
    optional int32 positionID = 4;
}
```

### GetOrderList — ProtoID 2004
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    optional string orderID = 3;
    optional int32 envFlag = 4;  // 1=SIMULATE, 2=REAL
}
```

### GetOrderFillList — ProtoID 2005
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    optional string orderID = 3;
    optional int32 envFlag = 4;
}
```

### GetHistoryOrderList — ProtoID 2006
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    optional string beginTime = 3;
    optional string endTime = 4;
    optional int32 envFlag = 5;
    optional int32 orderID = 6;
}
```

### GetMaxTrdQtys — ProtoID 2206
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    required Security security = 3;
    required int32 orderType = 4;  // OrderType
    required double price = 5;
}
```

### PlaceOrder — ProtoID 2201
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    required int32 orderType = 3;   // OrderType.NORMAL / MARKET
    required Security security = 4;
    required int32 trdSide = 5;     // 1=BUY, 2=SELL
    required double price = 6;
    required int64 qty = 7;
    optional string orderID = 8;    // for modify/cancel
    optional int32 modifyOrderOp = 9; // 1=normal, 2=cancel
}
```

### GetFlowSummary — ProtoID 2214
```protobuf
message C2S {
    required int32 trdMarket = 1;
    required int32 accID = 2;
    optional string beginTime = 3;
    optional string endTime = 4;
}
```

---

## Push Subscriptions (pkg/qot)

### Subscribe — ProtoID 3001
```protobuf
message C2S {
    repeated Security securityList = 1;
    repeated int32 subTypeList = 2;     // SubType
    required bool isSubOrUnSub = 3;    // true=订阅
    required bool isRegOrUnRegPush = 4; // true=注册推送
}
// SubType: 0=Basic, 1=KL_1Min, 2=KL_5Min, 3=KL_15Min, 4=KL_30Min,
//          5=KL_60Min, 6=KL_Day, 7=KL_Week, 8=KL_Month, 9=KL_Year,
//          13=OrderBook, 14=Ticker, 15=Broker
```

### UnsubscribeAll — ProtoID 3002

---

## Known Issues

### Empty C2S Bug (OpenD serverVer=1003)

OpenD rejects completely empty protobuf messages with "解析protobuf协议失败".

**Affected APIs:**
- `GetDelayStatistics` — C2S all optional, empty = `{0x0a, 0x00}`
- `GetTradeDate` — C2S all required, but SDK may send invalid defaults

**Fix:** Ensure C2S has at least one field set before marshalling.

**Debug command:** Set `FUTU_ADDR=host:port go run main.go` to connect to a different OpenD instance.
