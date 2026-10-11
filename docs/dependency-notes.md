# Dependency Notes

## google.golang.org/protobuf — transitive dependency

`google.golang.org/protobuf` v1.36.12 is a transitive dependency via the futuapi4go SDK.

### Dependency chain

```
examples/02_ticker
  → github.com/shing1211/futuapi4go/pkg/push
    → google.golang.org/protobuf/proto
```

### Risk assessment

**Low.** v1.36.12 is the current release. No known CVEs. This is the modern protobuf API (not the legacy `golang/protobuf` wrapper).

### Resolution path

Upgrade `futuapi4go` when a newer version of the SDK pins a newer protobuf release.

---

## golang.org/x/xerrors — not present

`go mod why -m golang.org/x/xerrors` confirms: `(main module does not need module golang.org/x/xerrors)`.

Unlike the other demo repos (tiger-go-demo, longbridge-go-demo), this dependency chain does not reach xerrors. No action needed.

## github.com/golang/protobuf — not present

`go mod why -m github.com/golang/protobuf` confirms: `(main module does not need module github.com/golang/protobuf)`.

The futuapi4go SDK uses the modern `google.golang.org/protobuf` directly, not the legacy wrapper. No action needed.
