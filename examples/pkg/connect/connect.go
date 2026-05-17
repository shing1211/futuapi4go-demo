// Package connect provides HA (High Availability) connection management for Futu OpenD.
//
// Features:
//   - Multi-host failover with parallel TCP probe
//   - Per-host RSA encryption configuration
//   - Auto-reconnect on connection loss with exponential backoff
//   - Keep-alive monitoring
//   - Connection state machine
//
// Basic usage:
//
//	mc := connect.MustConnect(context.Background())
//	defer mc.Close()
//	quote, err := client.GetQuote(ctx, mc.Client, market, code)
//
// Advanced usage with callbacks:
//
//	mc.OnStateChange = func(old, new connect.State) {
//	    fmt.Printf("State changed: %v -> %v\n", old, new)
//	}
//	mc.OnError = func(err error) {
//	    fmt.Printf("Connection error (recovering): %v\n", err)
//	}
package connect

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/pb/getglobalstate"
	"github.com/shing1211/futuapi4go/pkg/pb/common"
)

func init() {
	godotenv.Load()
}

// =============================================================================
// Configuration
// =============================================================================

var (
	defaultAddr    = "127.0.0.1:11111"
	defaultTimeout = 3 * time.Second
	defaultRSAPath = "/etc/futu/keys/private_key.pem"
)

// Host represents a single OpenD gateway with its configuration.
type Host struct {
	Host  string
	Port  int
	IsRSA bool
}

// ConnectionInfo contains metadata about the established connection.
type ConnectionInfo struct {
	Host       string
	Port       int
	TCPMs      float64
	APIMs      float64
	ServerVer  string
	QuoteLogin bool
	TradeLogin bool
	MarketHK   bool
	MarketUS   bool
	MarketSH   bool
	MarketSZ   bool
	RSAUsed    bool

	// GlobalState fields
	ConnID         uint64
	ServerBuildNo  int32
	ServerTime     int64
	QotSvrIpAddr   string
	TrdSvrIpAddr   string
	MarketHKRight  int32
	MarketUSRight  int32
	MarketSHRight  int32
	MarketSZRight  int32
	MarketHKFuture int32
	MarketUSFuture int32
}

// Config reads environment variables and returns configured hosts, RSA key contents, and TCP timeout.
func Config() (hosts []Host, rsaKey string, timeout time.Duration) {
	hosts = parseHosts()

	keyPath := os.Getenv("FUTU_RSA_KEY")
	if keyPath == "" {
		keyPath = defaultRSAPath
	}
	if data, err := os.ReadFile(keyPath); err == nil {
		rsaKey = strings.TrimSpace(string(data))
	}
	if rsaKey == "" {
		rsaKey = os.Getenv("FUTU_RSA_PUBKEY") // fallback to legacy env var
	}
	if t := os.Getenv("FUTU_TCP_TIMEOUT"); t != "" {
		if sec, err := strconv.ParseFloat(t, 64); err == nil {
			timeout = time.Duration(sec * float64(time.Second))
		}
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return
}

func parseHosts() []Host {
	envHosts := os.Getenv("FUTU_OPEND_HOSTS")
	result := []Host{}

	if envHosts != "" {
		for _, entry := range strings.Split(envHosts, ",") {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				continue
			}
			parts := strings.Split(entry, ":")
			hostStr := parts[0]
			port := 11111
			isRSA := true // default for remote hosts

			if len(parts) > 1 {
				if p, err := strconv.Atoi(parts[1]); err == nil {
					port = p
				}
			}
			if len(parts) > 2 {
				isRSA = strings.ToLower(parts[2]) == "true"
			}

			// Localhost defaults to no RSA unless explicitly set
			if isLocalhost(hostStr) && len(parts) <= 2 {
				isRSA = false
			}

			result = append(result, Host{Host: hostStr, Port: port, IsRSA: isRSA})
		}
	}

	if len(result) == 0 {
		// Fallback to FUTU_ADDR
		addr := os.Getenv("FUTU_ADDR")
		if addr == "" {
			addr = defaultAddr
		}
		hostStr, port, isRSA := parseSingleAddr(addr)
		result = append(result, Host{Host: hostStr, Port: port, IsRSA: isRSA})
	}

	return result
}

func parseSingleAddr(addr string) (host string, port int, isRSA bool) {
	isRSA = false // localhost assumed no RSA
	if idx := strings.LastIndex(addr, ":"); idx >= 0 {
		host = addr[:idx]
		if p, err := strconv.Atoi(addr[idx+1:]); err == nil {
			port = p
		}
	} else {
		host = addr
		port = 11111
	}
	if isLocalhost(host) {
		isRSA = false
	}
	return
}

func isLocalhost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1" || strings.HasPrefix(host, "127.")
}

// =============================================================================
// TCP Probe
// =============================================================================

type probeResult struct {
	host Host
	ms   *float64
}

func probeAllParallel(hosts []Host, timeout time.Duration) []*probeResult {
	results := make([]*probeResult, len(hosts))
	var wg sync.WaitGroup

	for i, h := range hosts {
		wg.Add(1)
		go func(i int, h Host) {
			defer wg.Done()
			ms := tcpProbe(h.Host, h.Port, timeout)
			results[i] = &probeResult{host: h, ms: ms}
		}(i, h)
	}

	wg.Wait()
	return results
}

func tcpProbe(host string, port int, timeout time.Duration) *float64 {
	t0 := time.Now()
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil
	}
	conn.Close()
	ms := time.Since(t0).Seconds() * 1000
	return &ms
}

func sortByLatency(results []*probeResult) []Host {
	reachable := make([]*probeResult, 0, len(results))
	for _, r := range results {
		if r.ms != nil {
			reachable = append(reachable, r)
		}
	}

	if len(reachable) == 0 {
		return nil
	}

	sort.Slice(reachable, func(i, j int) bool {
		return *reachable[i].ms < *reachable[j].ms
	})

	hosts := make([]Host, len(reachable))
	for i, r := range reachable {
		hosts[i] = r.host
		hosts[i].IsRSA = r.host.IsRSA
	}
	return hosts
}

// =============================================================================
// Backoff
// =============================================================================

const (
	backoffBase = time.Second
	backoffMax  = 60 * time.Second
)

type Backoff struct {
	attempt int
}

func (b *Backoff) Next() time.Duration {
	exp := backoffBase
	shift := b.attempt
	if shift > 63 {
		shift = 63
	}
	exp = time.Duration(uint64(exp) << shift)
	if exp > backoffMax {
		exp = backoffMax
	}
	jitter := time.Duration(rand.Int63n(int64(exp)))
	b.attempt++
	return jitter
}

func (b *Backoff) Reset() {
	b.attempt = 0
}

func NewTimer() *Backoff {
	return &Backoff{}
}

// =============================================================================
// State Machine
// =============================================================================

type State int

const (
	StateDisconnected State = iota
	StateConnecting
	StateConnected
	StateReconnecting
	StateFailed
)

func (s State) String() string {
	switch s {
	case StateDisconnected:
		return "Disconnected"
	case StateConnecting:
		return "Connecting"
	case StateConnected:
		return "Connected"
	case StateReconnecting:
		return "Reconnecting"
	case StateFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// =============================================================================
// Managed Connection
// =============================================================================

type ManagedConnection struct {
	Client *client.Client

	Info *ConnectionInfo

	State State

	OnStateChange func(old, new State)
	OnError       func(err error)
	OnConnect     func(*ConnectionInfo)

	mu            sync.RWMutex
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	keepAliveInt  time.Duration
	reconnectMode bool
	closed        int32
}

func (mc *ManagedConnection) Close() error {
	if !atomic.CompareAndSwapInt32(&mc.closed, 0, 1) {
		return nil
	}

	if mc.cancel != nil {
		mc.cancel()
	}

	mc.transitionTo(StateDisconnected)
	mc.wg.Wait()

	if mc.Client != nil {
		mc.Client.Close()
	}

	return nil
}

func (mc *ManagedConnection) transitionTo(to State) {
	mc.mu.Lock()
	old := mc.State
	mc.State = to
	mc.mu.Unlock()

	if mc.OnStateChange != nil {
		mc.OnStateChange(old, to)
	}
}

func (mc *ManagedConnection) setClient(cli *client.Client) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.Client != nil && mc.Client != cli {
		mc.Client.Close()
	}
	mc.Client = cli
}

func getGlobalState(ctx context.Context, cli *client.Client) (*getglobalstate.Response, error) {
	c2s := &getglobalstate.C2S{
		UserID: func() *uint64 { v := uint64(0); return &v }(),
	}

	pkt := &getglobalstate.Request{C2S: c2s}
	var rsp getglobalstate.Response

	if err := cli.Inner().RequestContext(ctx, 1002, pkt, &rsp); err != nil {
		return nil, err
	}

	if rsp.GetRetType() != int32(common.RetType_RetType_Succeed) {
		return nil, fmt.Errorf("get_global_state: retType=%d retMsg=%s", rsp.GetRetType(), rsp.GetRetMsg())
	}

	return &rsp, nil
}

func tryConnect(ctx context.Context, host Host, rsaKey string, useRSA bool) (*client.Client, error) {
	var cli *client.Client
	var err error

	if useRSA && rsaKey != "" {
		cli = client.New(client.WithRSAPublicKey(rsaKey))
	} else {
		cli = client.New()
	}

	addr := fmt.Sprintf("%s:%d", host.Host, host.Port)
	if err := cli.Connect(addr); err != nil {
		cli.Close()
		return nil, fmt.Errorf("connect to %s: %w", addr, err)
	}

	// Verify connection with GetGlobalState
	_, err = getGlobalState(ctx, cli)
	if err != nil {
		cli.Close()
		return nil, err
	}

	return cli, nil
}

// Connect establishes a managed connection with full HA features.
// It probes all configured hosts and connects to the fastest one.
// The connection will auto-reconnect on failure.
func Connect(ctx context.Context) (*ManagedConnection, error) {
	hosts, rsaKey, timeout := Config()

	log.Printf("[HA] Configured hosts:")
	for _, h := range hosts {
		rsaFlag := "RSA"
		if !h.IsRSA {
			rsaFlag = "plain"
		}
		log.Printf("[HA]   %s:%d (%s)", h.Host, h.Port, rsaFlag)
	}

	log.Printf("[HA] Probing %d hosts (timeout=%v)...", len(hosts), timeout)
	results := probeAllParallel(hosts, timeout)

	log.Printf("[HA] Probe results:")
	for _, r := range results {
		if r.ms != nil {
			log.Printf("[HA]   %s:%d reachable (%.2f ms)", r.host.Host, r.host.Port, *r.ms)
		} else {
			log.Printf("[HA]   %s:%d unreachable", r.host.Host, r.host.Port)
		}
	}

	sortedHosts := sortByLatency(results)

	if len(sortedHosts) == 0 {
		log.Printf("[HA] No reachable hosts!")
		return nil, fmt.Errorf("connect: no reachable OpenD gateways")
	}

	log.Printf("[HA] Sorted by latency:")
	for _, h := range sortedHosts {
		rsaFlag := "RSA"
		if !h.IsRSA {
			rsaFlag = "plain"
		}
		log.Printf("[HA]   # %s:%d (%s)", h.Host, h.Port, rsaFlag)
	}

	mc := &ManagedConnection{
		keepAliveInt:  30 * time.Second,
		reconnectMode: true,
	}
	mc.ctx, mc.cancel = context.WithCancel(context.Background())

	for _, host := range sortedHosts {
		log.Printf("[HA] Trying %s:%d (RSA=%v)...", host.Host, host.Port, host.IsRSA)
		cli, err := tryConnect(mc.ctx, host, rsaKey, host.IsRSA)
		if err == nil {
			log.Printf("[HA] Connected to %s:%d (RSA=%v)", host.Host, host.Port, host.IsRSA)
			mc.setClient(cli)
			mc.Info = &ConnectionInfo{
				Host:    host.Host,
				Port:    host.Port,
				RSAUsed: host.IsRSA,
			}
			mc.transitionTo(StateConnected)
			mc.startKeepAlive()
			mc.startReconnectMonitor()
			if mc.OnConnect != nil {
				mc.OnConnect(mc.Info)
			}
			return mc, nil
		}
		log.Printf("[HA] %s:%d (RSA=%v) failed: %v", host.Host, host.Port, host.IsRSA, err)

		// Fallback: try same host without RSA only if RSA was requested
		if host.IsRSA {
			log.Printf("[HA] Trying %s:%d (plain fallback)...", host.Host, host.Port)
			cli, err = tryConnect(mc.ctx, host, rsaKey, false)
			if err == nil {
				log.Printf("[HA] Connected to %s:%d (plain fallback)", host.Host, host.Port)
				mc.setClient(cli)
				mc.Info = &ConnectionInfo{
					Host:    host.Host,
					Port:    host.Port,
					RSAUsed: false,
				}
				mc.transitionTo(StateConnected)
				mc.startKeepAlive()
				mc.startReconnectMonitor()
				if mc.OnConnect != nil {
					mc.OnConnect(mc.Info)
				}
				return mc, nil
			}
			log.Printf("[HA] %s:%d (plain fallback) failed: %v", host.Host, host.Port, err)
		}
	}

	mc.cancel()
	log.Printf("[HA] All hosts exhausted")
	return nil, fmt.Errorf("connect: all hosts failed")
}

// MustConnect calls Connect and log.Fatal on error.
func MustConnect(ctx context.Context) *ManagedConnection {
	mc, err := Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return mc
}

// ConnectWS establishes a managed WebSocket connection.
// TODO: Implement WS HA support.
func ConnectWS(ctx context.Context, secretKey string) (*ManagedConnection, error) {
	return nil, fmt.Errorf("ConnectWS: not implemented yet")
}

// MustConnectWS calls ConnectWS and log.Fatal on error.
func MustConnectWS(ctx context.Context, secretKey string) *ManagedConnection {
	mc, err := ConnectWS(ctx, secretKey)
	if err != nil {
		log.Fatal(err)
	}
	return mc
}

func (mc *ManagedConnection) startKeepAlive() {
	mc.wg.Add(1)
	go func() {
		defer mc.wg.Done()

		ticker := time.NewTicker(mc.keepAliveInt)
		defer ticker.Stop()

		for {
			select {
			case <-mc.ctx.Done():
				return
			case <-ticker.C:
				if atomic.LoadInt32(&mc.closed) == 1 {
					return
				}
				if err := mc.keepAlive(); err != nil {
					if mc.OnError != nil {
						mc.OnError(fmt.Errorf("keepalive failed: %w", err))
					}
				}
			}
		}
	}()
}

func (mc *ManagedConnection) keepAlive() error {
	mc.mu.RLock()
	cli := mc.Client
	mc.mu.RUnlock()

	if cli == nil {
		return fmt.Errorf("client is nil")
	}

	ctx, cancel := context.WithTimeout(mc.ctx, 10*time.Second)
	defer cancel()

	_, err := getGlobalState(ctx, cli)
	return err
}

func (mc *ManagedConnection) startReconnectMonitor() {
	mc.wg.Add(1)
	go func() {
		defer mc.wg.Done()

		for {
			select {
			case <-mc.ctx.Done():
				return
			case <-time.After(5 * time.Second):
				if atomic.LoadInt32(&mc.closed) == 1 {
					return
				}
				mc.checkAndReconnect()
			}
		}
	}()
}

func (mc *ManagedConnection) checkAndReconnect() {
	mc.mu.RLock()
	cli := mc.Client
	state := mc.State
	mc.mu.RUnlock()

	if cli == nil || state != StateConnected {
		return
	}

	ctx, cancel := context.WithTimeout(mc.ctx, 5*time.Second)
	_, err := getGlobalState(ctx, cli)
	cancel()

	if err != nil {
		mc.transitionTo(StateReconnecting)
		mc.reconnect()
	}
}

func (mc *ManagedConnection) reconnect() {
	mc.mu.RLock()
	isReconnecting := mc.reconnectMode
	mc.mu.RUnlock()

	if !isReconnecting {
		return
	}

	hosts, rsaKey, timeout := Config()
	timer := NewTimer()

	for {
		select {
		case <-mc.ctx.Done():
			return
		default:
		}

		if atomic.LoadInt32(&mc.closed) == 1 {
			return
		}

		sortedHosts := sortByLatency(probeAllParallel(hosts, timeout))
		if len(sortedHosts) == 0 {
			if mc.OnError != nil {
				mc.OnError(fmt.Errorf("reconnect: no reachable hosts"))
			}
			timer.Next()
			continue
		}

		for _, host := range sortedHosts {
			cli, err := tryConnect(mc.ctx, host, rsaKey, host.IsRSA)
			if err == nil {
				mc.setClient(cli)
				mc.Info = &ConnectionInfo{
					Host:    host.Host,
					Port:    host.Port,
					RSAUsed: host.IsRSA,
				}
				mc.transitionTo(StateConnected)
				timer.Reset()
				if mc.OnError != nil {
					mc.OnError(fmt.Errorf("reconnected to %s:%d", host.Host, host.Port))
				}
				if mc.OnConnect != nil {
					mc.OnConnect(mc.Info)
				}
				return
			}

			// Try opposite RSA
			cli, err = tryConnect(mc.ctx, host, rsaKey, !host.IsRSA)
			if err == nil {
				mc.setClient(cli)
				mc.Info = &ConnectionInfo{
					Host:    host.Host,
					Port:    host.Port,
					RSAUsed: !host.IsRSA,
				}
				mc.transitionTo(StateConnected)
				timer.Reset()
				if mc.OnError != nil {
					mc.OnError(fmt.Errorf("reconnected to %s:%d", host.Host, host.Port))
				}
				if mc.OnConnect != nil {
					mc.OnConnect(mc.Info)
				}
				return
			}
		}

		if mc.OnError != nil {
			mc.OnError(fmt.Errorf("reconnect: all hosts failed"))
		}
		timer.Next()
	}
}