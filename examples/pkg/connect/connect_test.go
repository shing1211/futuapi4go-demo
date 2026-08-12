package connect

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Backoff
// ---------------------------------------------------------------------------

func TestBackoff_NeverExceedsMax(t *testing.T) {
	b := &Backoff{}
	for i := 0; i < 1000; i++ {
		d := b.Next()
		if d < 0 {
			t.Fatalf("iteration %d: negative duration %v", i, d)
		}
		if d > backoffMax {
			t.Fatalf("iteration %d: Next() = %v exceeds max %v", i, d, backoffMax)
		}
	}
}

func TestBackoff_Window(t *testing.T) {
	if backoffWindow(0) != backoffBase {
		t.Fatalf("backoffWindow(0) = %v, want %v", backoffWindow(0), backoffBase)
	}
	if backoffWindow(1) != backoffBase*2 {
		t.Fatalf("backoffWindow(1) = %v, want %v", backoffWindow(1), backoffBase*2)
	}
	// Large attempts must cap at backoffMax without overflow (regression for
	// the old `1s << attempt` bug that panicked rand.Int63n).
	if got := backoffWindow(100); got != backoffMax {
		t.Fatalf("backoffWindow(100) = %v, want %v (capped)", got, backoffMax)
	}
}

func TestBackoff_WindowGrowsThenCaps(t *testing.T) {
	prev := backoffWindow(0)
	capped := false
	for i := 1; i < 64; i++ {
		cur := backoffWindow(i)
		if cur > backoffMax {
			t.Fatalf("backoffWindow(%d) = %v exceeds max %v", i, cur, backoffMax)
		}
		if cur == backoffMax {
			capped = true
		}
		if cur < prev {
			t.Fatalf("backoffWindow(%d) = %v shrank from %v", i, cur, prev)
		}
		prev = cur
	}
	if !capped {
		t.Fatalf("window never reached backoffMax in 64 attempts")
	}
}

func TestBackoff_JitterWithinWindow(t *testing.T) {
	b := &Backoff{}
	for i := 0; i < 20; i++ {
		d := b.Next()
		w := backoffWindow(b.attempt - 1)
		if d >= w {
			t.Fatalf("attempt %d: Next() = %v not < window %v", i, d, w)
		}
	}
}

func TestBackoff_Reset(t *testing.T) {
	b := &Backoff{}
	b.Next()
	b.Next() // advance
	b.Next() // advance further
	b.Reset()
	after := b.Next()

	// After Reset, the random window returns to the base window. We can't be
	// exact (jitter), but the first window is < 1s and subsequent grows.
	if after < 0 || after >= backoffBase {
		t.Fatalf("Next() after Reset = %v, want in [0, %v)", after, backoffBase)
	}
}

// ---------------------------------------------------------------------------
// sortByLatency
// ---------------------------------------------------------------------------

func TestSortByLatency_OrderedByFastest(t *testing.T) {
	fast := 1.0
	slow := 9.0
	results := []*probeResult{
		{host: Host{Host: "slow", IsRSA: true}, ms: &slow},
		{host: Host{Host: "fast", IsRSA: true}, ms: &fast},
	}
	got := sortByLatency(results)
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Host != "fast" || got[1].Host != "slow" {
		t.Fatalf("sorted order = %v, %v; want fast, slow", got[0].Host, got[1].Host)
	}
	// RSA flag preserved from input host
	if !got[0].IsRSA {
		t.Fatalf("fast host IsRSA = false, want true (preserved)")
	}
}

func TestSortByLatency_FiltersUnreachable(t *testing.T) {
	results := []*probeResult{
		{host: Host{Host: "dead1"}, ms: nil},
		{host: Host{Host: "alive"}, ms: &[]float64{2.0}[0]},
		{host: Host{Host: "dead2"}, ms: nil},
	}
	got := sortByLatency(results)
	if len(got) != 1 || got[0].Host != "alive" {
		t.Fatalf("got = %v, want only [alive]", got)
	}
}

func TestSortByLatency_EmptyReturnsNil(t *testing.T) {
	got := sortByLatency(nil)
	if got != nil {
		t.Fatalf("sortByLatency(nil) = %v, want nil", got)
	}
	got = sortByLatency([]*probeResult{{host: Host{Host: "dead"}, ms: nil}})
	if got != nil {
		t.Fatalf("all-unreachable = %v, want nil", got)
	}
}

// ---------------------------------------------------------------------------
// isLocalhost
// ---------------------------------------------------------------------------

func TestIsLocalhost(t *testing.T) {
	cases := []struct {
		host string
		want bool
	}{
		{"localhost", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"127.8.9.10", true},
		{"192.168.1.1", false},
		{"10.0.0.1", false},
		{"opend.example.com", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isLocalhost(c.host); got != c.want {
			t.Errorf("isLocalhost(%q) = %v, want %v", c.host, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// parseSingleAddr
// ---------------------------------------------------------------------------

func TestParseSingleAddr_HostPort(t *testing.T) {
	host, port, isRSA := parseSingleAddr("10.0.0.5:22222")
	if host != "10.0.0.5" || port != 22222 || !isRSA {
		t.Fatalf("got host=%q port=%d rsa=%v, want host=10.0.0.5 port=22222 rsa=true (remote)", host, port, isRSA)
	}
}

func TestParseSingleAddr_DefaultPort(t *testing.T) {
	host, port, isRSA := parseSingleAddr("192.168.0.2")
	if host != "192.168.0.2" || port != 11111 || !isRSA {
		t.Fatalf("got host=%q port=%d rsa=%v, want default port 11111 rsa=true (remote defaults to RSA)", host, port, isRSA)
	}
}

func TestParseSingleAddr_RSAOnlyRemote(t *testing.T) {
	_, _, isRSA := parseSingleAddr("opend.example.com:11111")
	if !isRSA {
		t.Fatalf("remote host rsa = false, want true (remote defaults to RSA)")
	}
}

// ---------------------------------------------------------------------------
// parseHosts (env-driven)
// ---------------------------------------------------------------------------

func TestParseHosts_FUTUADDRFallback(t *testing.T) {
	t.Setenv("FUTU_OPEND_HOSTS", "")
	t.Setenv("FUTU_ADDR", "127.0.0.1:12345")
	hosts := parseHosts()
	if len(hosts) != 1 {
		t.Fatalf("len(hosts) = %d, want 1", len(hosts))
	}
	if hosts[0].Port != 12345 || hosts[0].IsRSA {
		t.Fatalf("hosts[0] = %+v, want port 12345 rsa=false", hosts[0])
	}
}

func TestParseHosts_MultiHost(t *testing.T) {
	t.Setenv("FUTU_OPEND_HOSTS", "local:11111:false,remote.example.com:11111:true")
	t.Setenv("FUTU_ADDR", "") // ensure no fallback
	hosts := parseHosts()
	if len(hosts) != 2 {
		t.Fatalf("len(hosts) = %d, want 2", len(hosts))
	}
	if hosts[0].Host != "local" || hosts[0].Port != 11111 || hosts[0].IsRSA {
		t.Fatalf("hosts[0] = %+v, want local:11111 rsa=false", hosts[0])
	}
	if hosts[1].Host != "remote.example.com" || hosts[1].IsRSA != true {
		t.Fatalf("hosts[1] = %+v, want remote.example.com rsa=true", hosts[1])
	}
}

func TestParseHosts_LocalhostDefaultsNoRSA(t *testing.T) {
	t.Setenv("FUTU_OPEND_HOSTS", "127.0.0.1:11111")
	t.Setenv("FUTU_ADDR", "")
	hosts := parseHosts()
	if len(hosts) != 1 || hosts[0].IsRSA {
		t.Fatalf("localhost should default to no RSA, got %+v", hosts)
	}
}

// ---------------------------------------------------------------------------
// State machine
// ---------------------------------------------------------------------------

func TestManagedConnection_StateTransitions(t *testing.T) {
	mc := &ManagedConnection{}
	var transitions []State
	mc.OnStateChange = func(old, _ State) { transitions = append(transitions, old) }

	mc.transitionTo(StateConnecting)
	mc.transitionTo(StateConnected)
	mc.transitionTo(StateDisconnected)

	if len(transitions) != 3 {
		t.Fatalf("len(transitions) = %d, want 3", len(transitions))
	}
	if transitions[0] != StateDisconnected || transitions[1] != StateConnecting || transitions[2] != StateConnected {
		t.Fatalf("transitions = %v, want [Disconnected Connecting Connected]", transitions)
	}
	if mc.State != StateDisconnected {
		t.Fatalf("final state = %v, want Disconnected", mc.State)
	}
}

func TestState_String(t *testing.T) {
	cases := []struct {
		s    State
		want string
	}{
		{StateDisconnected, "Disconnected"},
		{StateConnecting, "Connecting"},
		{StateConnected, "Connected"},
		{StateReconnecting, "Reconnecting"},
		{StateFailed, "Failed"},
		{State(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("State(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkBackoff_Next(b *testing.B) {
	bt := &Backoff{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = bt.Next()
	}
}

func BenchmarkSortByLatency_5hosts(b *testing.B) {
	results := []*probeResult{
		{host: Host{Host: "a"}, ms: &[]float64{5}[0]},
		{host: Host{Host: "b"}, ms: nil},
		{host: Host{Host: "c"}, ms: &[]float64{1}[0]},
		{host: Host{Host: "d"}, ms: &[]float64{3}[0]},
		{host: Host{Host: "e"}, ms: &[]float64{2}[0]},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = sortByLatency(results)
	}
}

func BenchmarkProbeAllParallel_3hosts(b *testing.B) {
	hosts := []Host{
		{Host: "localhost", Port: 11111, IsRSA: false},
		{Host: "127.0.0.1", Port: 11111, IsRSA: false},
		{Host: "::1", Port: 11111, IsRSA: false},
	}
	timeout := 200 * time.Millisecond
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = probeAllParallel(hosts, timeout)
	}
}