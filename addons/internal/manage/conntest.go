package manage

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"strconv"
	"strings"
	"sync"

	"time"
)

// Variables rather than constants so a test can run the whole thing in
// seconds.
var (

	// connTestConnectWait is how long a tunnel has to carry its first echo.
	connTestConnectWait = 60 * time.Second
	// connTestSoak is how many one-second echoes a tunnel must carry. Long
	// enough for a path that cuts flows after a few seconds to show it.
	connTestSoak = 60
	// connTestSlack is how long past its schedule the kharej keeps waiting
	// for the verdict before it stops on its own.
	connTestSlack = 5 * time.Minute
)

// connTestBulk is the transfer that measures speed.
const connTestBulk = 1 << 20

// connTestBinary is the engine to run; a variable so a test can point it at a
// built binary rather than at the test itself.
var connTestBinary = os.Executable

// ConnTestLink is everything the kharej needs to build its end of every
// tunnel under test, and to find the Iran side's coordinator.
type ConnTestLink struct {
	V     int    `json:"v"`
	Kind  string `json:"k"` // always connTestKind
	ID    string `json:"id"`
	Host  string `json:"h"`  // the Iran server's address
	Coord int    `json:"c"`  // the coordinator's port, TCP and UDP
	Tok   string `json:"t"`  // the coordinator's secret
	TCP   int    `json:"te"` // kharej: the TCP echo reverse tunnels forward to
	UDP   int    `json:"ue"` // kharej: the UDP echo the udp transport forwards to
	L3    int    `json:"le"` // kharej: the TCP echo on each direct tunnel's address
	Until int64  `json:"u"`  // unix seconds: the kharej stops by then whatever happens
	// Preset is the performance preset every tunnel under test was built
	// with, for the kharej's report; each case carries it too.
	Preset string `json:"pr,omitempty"`
	// The spoofing check: the forged source, and the UDP port each side's
	// receiver counts probes on. Empty SpoofSrc is no check.
	SpoofSrc string `json:"ss,omitempty"`
	SpoofK   int    `json:"sk,omitempty"` // the kharej's receiver
	SpoofI   int    `json:"si,omitempty"` // the Iran server's receiver
	SpoofN   int    `json:"sn,omitempty"` // probes each way
	// Cases are the setup links of the tunnels under test, exactly as the
	// wizard would print them, so the kharej builds each one by the path a
	// real tunnel takes.
}

// The link the operator copies is only where the Iran side's coordinator is
// and the test's secret: bk://t. and some twenty-seven characters. The
// rest — every tunnel's settings — the kharej asks the coordinator for (see
// ctConfig), and every tunnel's token is derived from the secret rather than
// carried, which is what keeps that answer inside one packet: a path that cuts
// a flow after a few packets still carries it.
//
// The body is base64url of: version (1), port (2, big endian), secret (12),
// then the host — 4 bytes for an IPv4 address, or the name as text.
const (
	connTestScheme   = shareScheme + "t."
	connTestSecret   = 12
	connTestLinkVer  = 1
	connTestFixedLen = 1 + 2 + connTestSecret
)

// Short is the link to copy.
func (l ConnTestLink) Short() string {
	secret, _ := base64.RawURLEncoding.DecodeString(l.Tok)
	b := []byte{connTestLinkVer, byte(l.Coord >> 8), byte(l.Coord)}
	b = append(b, secret...)
	if ip := net.ParseIP(l.Host).To4(); ip != nil {
		b = append(b, ip...)
	} else {
		b = append(b, l.Host...)
	}
	return connTestScheme + base64.RawURLEncoding.EncodeToString(b)
}

// connTestAddr is what a copied link says: where to ask, and the secret.
type connTestAddr struct {
	Host  string
	Coord int
	Tok   string
}

// ctNewSecret is a test's secret, as the link and the coordinator carry it.
func ctNewSecret() string {
	b := make([]byte, connTestSecret)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// ctCaseToken is one tunnel's token, derived from the test's secret so that
// neither the link nor the coordinator's answer has to carry it.
func ctCaseToken(secret, kind, tr string) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte("bk-conntest/" + kind + "/" + tr))
	return hex.EncodeToString(m.Sum(nil))[:32]
}

func gzipB64(raw []byte) string {
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, _ = zw.Write(raw)
	_ = zw.Close()
	return base64.RawURLEncoding.EncodeToString(buf.Bytes())
}

func unGzipB64(s string) ([]byte, error) {
	gz, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(s))
	if err != nil {
		return nil, err
	}
	zr, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(io.LimitReader(zr, 256<<10))
}

// ConnTestResult is one tunnel's verdict.
type ConnTestResult struct {
	Kind      string  `json:"k"` // "reverse" or "direct"
	Transport string  `json:"t"`
	Status    string  `json:"s"` // ok, unstable, down, skipped
	Detail    string  `json:"d,omitempty"`
	Connect   float64 `json:"c,omitempty"` // seconds to the first echo
	OK        int     `json:"o"`           // echoes that came back
	Tried     int     `json:"n"`           // echoes sent
	RTTms     int     `json:"r,omitempty"`
	Mbps      float64 `json:"m,omitempty"`
	// Total is how many echoes the test sends each tunnel.
	Total int `json:"z,omitempty"`
}

// Result statuses.
const (
	ctTesting  = "testing" // still running
	ctOK       = "ok"
	ctUnstable = "unstable"
	ctDown     = "down"
	ctSkipped  = "skipped"
)

// connTestCase is one tunnel under test, as the Iran side holds it.
type connTestCase struct {
	kind, tr string
	name     string
	udp      bool   // the forwarded traffic is UDP (the udp transport)
	entry    int    // reverse: where on this server the tunnel's users arrive
	peerIP   string // direct: the kharej's tunnel address, where its echo listens
	skip     string
	rtts     []time.Duration // every echo's round trip, for the recommendation
}

// The two directions the spoofing check runs, as its rows name them.
const (
	ctSpoofToKharej = "iran-kharej"
	ctSpoofToIran   = "kharej-iran"
)

// ctProbe pushes traffic through one tunnel and judges it.
// report is told the tunnel's state whenever it changes: after every echo,
// and once more with the verdict, which is also returned.
func ctProbe(ctx context.Context, c *connTestCase, l3Echo int, report func(ConnTestResult)) ConnTestResult {
	r := ConnTestResult{Kind: c.kind, Transport: c.tr, Status: ctTesting, Total: connTestSoak}
	say := func(string) {
		if report != nil {
			report(r)
		}
	}
	network, addr := "tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(c.entry))
	if c.udp {
		network = "udp"
	}
	if c.kind == "direct" {
		addr = net.JoinHostPort(c.peerIP, strconv.Itoa(l3Echo))
	}
	dial := func() (net.Conn, error) { return net.DialTimeout(network, addr, 3*time.Second) }

	// Up: the first echo that comes back.
	start := time.Now()
	var conn net.Conn
	for conn == nil {
		if ctx.Err() != nil {
			r.Status, r.Detail = ctDown, "stopped"
			return r
		}
		if time.Since(start) > connectWaitFor(c) {
			r.Status, r.Detail = ctDown, fmt.Sprintf("never carried traffic in %s", connectWaitFor(c))
			if c.tr == "spoof" {
				r.Detail += " — the forged source was dropped: by the provider on the way out, " +
					"or by a strict rp_filter on the kharej"
			}
			say("never came up")
			return r
		}
		if cn, err := dial(); err == nil {
			if ctEcho(cn, c.udp) == nil {
				conn = cn
				break
			}
			cn.Close()
		}
		ctSleep(ctx, time.Second)
	}
	r.Connect = time.Since(start).Seconds()
	say(fmt.Sprintf("up after %.0fs", r.Connect))

	// Staying up: an echo a second over the same connection, and a new one
	// only when that one is lost — which is exactly what a user's connection
	// would see.
	var rtt time.Duration
	var firstLoss time.Duration
	soakStart := time.Now()
	for i := 0; i < connTestSoak && ctx.Err() == nil; i++ {
		ctSleep(ctx, time.Until(soakStart.Add(time.Duration(i)*time.Second)))
		r.Tried++
		if conn == nil {
			cn, err := dial()
			if err != nil {
				if firstLoss == 0 {
					firstLoss = time.Since(soakStart)
				}
				say("")
				continue
			}
			conn = cn
		}
		t0 := time.Now()
		if err := ctEcho(conn, c.udp); err != nil {
			conn.Close()
			conn = nil
			if firstLoss == 0 {
				firstLoss = time.Since(soakStart)
			}
			say("")
			continue
		}
		r.OK++
		took := time.Since(t0)
		rtt += took
		c.rtts = append(c.rtts, took)
		say("")
	}
	if conn != nil {
		conn.Close()
	}
	if r.OK > 0 {
		r.RTTms = int((rtt / time.Duration(r.OK)).Milliseconds())
	}

	// Speed, over a fresh connection, for the tunnels that carry TCP.
	bulkErr := error(nil)
	if !c.udp && ctx.Err() == nil {
		r.Mbps, bulkErr = ctBulk(dial)
	}

	switch {
	case r.OK >= r.Tried-1 && r.Tried > 0 && bulkErr == nil:
		r.Status = ctOK
	case r.OK == 0:
		r.Status = ctDown
		r.Detail = fmt.Sprintf("came up, then carried nothing (lost after %.0fs)", firstLoss.Seconds())
	default:
		r.Status = ctUnstable
		if firstLoss > 0 {
			r.Detail = fmt.Sprintf("first loss after %.0fs", firstLoss.Seconds())
		}
		if bulkErr != nil {
			r.Detail = strings.TrimPrefix(r.Detail+"; bulk transfer failed ("+ctShortErr(bulkErr)+")", "; ")
		}
	}
	say("done: " + r.Status)
	return r
}

// connectWaitFor is how long a tunnel is given to come up. The direct ones
// wait for the kharej to have built its interface as well.
func connectWaitFor(c *connTestCase) time.Duration {
	if c.kind == "direct" {
		return connTestConnectWait + 30*time.Second
	}
	return connTestConnectWait
}

// ctEcho sends a small random payload and requires it back unchanged.
func ctEcho(c net.Conn, udp bool) error {
	msg := make([]byte, 512)
	_, _ = rand.Read(msg)
	_ = c.SetDeadline(time.Now().Add(4 * time.Second))
	defer c.SetDeadline(time.Time{})
	if _, err := c.Write(msg); err != nil {
		return err
	}
	got := make([]byte, len(msg))
	if udp {
		n, err := c.Read(got)
		if err != nil {
			return err
		}
		got = got[:n]
	} else if _, err := io.ReadFull(c, got); err != nil {
		return err
	}
	if !bytes.Equal(got, msg) {
		return errors.New("the echo came back different")
	}
	return nil
}

// ctBulk sends connTestBulk bytes through a fresh connection and reads them
// back, and reports the speed.
func ctBulk(dial func() (net.Conn, error)) (float64, error) {
	c, err := dial()
	if err != nil {
		return 0, err
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(40 * time.Second))
	payload := make([]byte, connTestBulk)
	_, _ = rand.Read(payload)
	start := time.Now()
	werr := make(chan error, 1)
	go func() {
		_, err := c.Write(payload)
		werr <- err
	}()
	got := make([]byte, len(payload))
	if n, err := io.ReadFull(c, got); err != nil {
		return 0, fmt.Errorf("%d of %d KB back: %w", n>>10, len(payload)>>10, err)
	}
	if err := <-werr; err != nil {
		return 0, err
	}
	if !bytes.Equal(got, payload) {
		return 0, errors.New("the transfer came back different")
	}
	secs := time.Since(start).Seconds()
	return float64(2*connTestBulk*8) / secs / 1e6, nil
}

func ctSleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

// ctPickPort is a random port free on this server (for TCP, and for UDP as
// well when both is set) and not already given to another part of the test.
// The kharej's echo ports are picked here too; they are high and random, so a
// clash there is unlikely, and a tunnel that meets one reports it.
func ctPickPort(used map[int]bool, both bool) int {
	for i := 0; i < 500; i++ {
		var b [2]byte
		_, _ = rand.Read(b[:])
		p := 20000 + (int(b[0])<<8|int(b[1]))%40000
		if used[p] {
			continue
		}
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err != nil {
			continue
		}
		l.Close()
		if both {
			u, err := net.ListenPacket("udp", fmt.Sprintf(":%d", p))
			if err != nil {
				continue
			}
			u.Close()
		}
		used[p] = true
		return p
	}
	return 0
}

// ctEngine is one engine process under test.
type ctEngine struct {
	name string
	cmd  *exec.Cmd
	log  string
	done chan struct{}
}

// exited reports whether the engine has stopped on its own, and the last
// thing it said.
func (e *ctEngine) exited() (bool, string) {
	select {
	case <-e.done:
		return true, ctLastLine(e.log)
	default:
		return false, ""
	}
}

func ctLastLine(path string) string {
	b, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if l := strings.TrimSpace(lines[i]); l != "" {
			return l
		}
	}
	return ""
}

// The coordinator: one question and one answer per connection, over TCP and
// over UDP, so that a path cutting flows after a few packets — or one that
// carries only one of the two — still lets the two sides talk.
//
//	config <token>  → config <c>  the test, for the kharej to build its end from
//	hello <token>   → ok          the kharej is in; its address is noted
//	result <token>  → wait | done <verdict>
type ctCoordinator struct {
	tok     string
	tcp     net.Listener
	udp     net.PacketConn
	mu      sync.Mutex
	peer    string
	verdict string
	config  string // the answer to "config"; see ctConfig
	live    func() []ConnTestResult
	joined  chan struct{}
	fetched chan struct{}
	joinOne sync.Once
	fetchOn sync.Once
	// spoofArrived is the kharej's count of the Iran server's forged probes.
	spoofArrived chan int
	spoofOnce    sync.Once
	// pmtu is the path MTU the kharej measured.
	pmtu     chan int
	pmtuOnce sync.Once
}

// answer is the coordinator's whole protocol. An empty answer is silence: a
// request without the token gets nothing, not even a refusal.
//
//	spoof <token> <n>  → ok      the kharej counted n forged probes
//	pmtu <token> <n> <pad> → pm <n>  a path-MTU probe of n bytes
//	pmtu <token> <n>   → ok      the path MTU the kharej measured
func (c *ctCoordinator) answer(line, from string) string {
	f := strings.Fields(line)
	if len(f) < 2 || subtle.ConstantTimeCompare([]byte(f[1]), []byte(c.tok)) != 1 {
		return ""
	}
	// A path-MTU probe: its padding is the point, and the answer is small.
	if f[0] == "pmtu" && len(f) == 4 {
		return "pm " + f[2]
	}
	if len(f) > 3 {
		return ""
	}
	if f[0] == "pmtu" {
		n, err := strconv.Atoi(f[len(f)-1])
		if len(f) != 3 || err != nil || n < 0 {
			return ""
		}
		c.pmtuOnce.Do(func() { c.pmtu <- n })
		return "ok"
	}
	if f[0] == "spoof" {
		n, err := strconv.Atoi(f[len(f)-1])
		if len(f) != 3 || err != nil || n < 0 {
			return ""
		}
		// Asked again when an answer was lost; the first count is the count.
		c.spoofOnce.Do(func() { c.spoofArrived <- n })
		return "ok"
	}
	if len(f) != 2 {
		return ""
	}
	switch f[0] {
	case "config":
		c.mu.Lock()
		cfg := c.config
		c.mu.Unlock()
		if cfg == "" {
			return ""
		}
		return "config " + cfg
	case "hello":
		c.mu.Lock()
		if c.peer == "" {
			c.peer = strings.TrimPrefix(from, "::ffff:")
		}
		c.mu.Unlock()
		c.joinOne.Do(func() { close(c.joined) })
		return "ok"
	case "result":
		c.mu.Lock()
		v, live := c.verdict, c.live
		c.mu.Unlock()
		if v == "" {
			if live == nil {
				return "wait"
			}
			raw, _ := json.Marshal(live())
			return "live " + gzipB64(raw)
		}
		c.fetchOn.Do(func() { close(c.fetched) })
		return "done " + v
	}
	return ""
}

// setConfig sets what the kharej is told the test is.
func (c *ctCoordinator) setConfig(cfg string) {
	c.mu.Lock()
	c.config = cfg
	c.mu.Unlock()
}

// setLive hands the coordinator the running test's rows, for the kharej's
// screen.
func (c *ctCoordinator) setLive(rows func() []ConnTestResult) {
	c.mu.Lock()
	c.live = rows
	c.mu.Unlock()
}

// ctVerdict is what the kharej fetches when the test is over.
type ctVerdict struct {
	Results []ConnTestResult `json:"r"`
	Best    ConnTestBest     `json:"b"`
}

func (c *ctCoordinator) publish(results []ConnTestResult, best ConnTestBest) {
	raw, _ := json.Marshal(ctVerdict{Results: results, Best: best})
	c.mu.Lock()
	c.verdict = gzipB64(raw)
	c.mu.Unlock()
}

func (c *ctCoordinator) close() {
	c.tcp.Close()
	c.udp.Close()
}

// ctAsk puts one question to the coordinator, over TCP and then UDP.
func ctAsk(host string, port int, line string) (string, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	var last error
	for _, network := range []string{"tcp", "udp"} {
		conn, err := net.DialTimeout(network, addr, 5*time.Second)
		if err != nil {
			last = err
			continue
		}
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Write([]byte(line + "\n"))
		if err == nil {
			buf := make([]byte, 8192)
			var n int
			if network == "udp" {
				n, err = conn.Read(buf)
			} else {
				n, err = io.ReadFull(conn, buf)
				if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
					err = nil
				}
			}
			if err == nil && n > 0 {
				conn.Close()
				return strings.TrimSpace(string(buf[:n])), nil
			}
		}
		conn.Close()
		last = err
	}
	if last == nil {
		last = errors.New("no answer")
	}
	return "", last
}

// ctEchoes are the kharej's echo servers.
type ctEchoes struct {
	mu      sync.Mutex
	closers []io.Closer
}

func ctStartEchoes(ctx context.Context, link ConnTestLink) (*ctEchoes, error) {
	e := &ctEchoes{}
	tl, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(link.TCP)))
	if err != nil {
		return nil, fmt.Errorf("port %d is taken on this server — start a new test on the Iran server: %w", link.TCP, err)
	}
	e.add(tl)
	go ctServeTCPEcho(tl)
	ul, err := net.ListenPacket("udp", net.JoinHostPort("127.0.0.1", strconv.Itoa(link.UDP)))
	if err != nil {
		e.close()
		return nil, fmt.Errorf("port %d is taken on this server — start a new test on the Iran server: %w", link.UDP, err)
	}
	e.add(ul)
	go ctServeUDPEcho(ul)
	return e, nil
}

func (e *ctEchoes) add(c io.Closer) {
	e.mu.Lock()
	e.closers = append(e.closers, c)
	e.mu.Unlock()
}

func (e *ctEchoes) close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, c := range e.closers {
		c.Close()
	}
	e.closers = nil
}

func ctServeTCPEcho(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			return
		}
		go func() {
			defer c.Close()
			_, _ = io.Copy(c, c)
		}()
	}
}

func ctServeUDPEcho(l net.PacketConn) {
	buf := make([]byte, 65535)
	for {
		n, from, err := l.ReadFrom(buf)
		if err != nil {
			return
		}
		_, _ = l.WriteTo(buf[:n], from)
	}
}

// ctShortErr is an error without the addresses a net error repeats.
func ctShortErr(err error) string {
	msg := err.Error()
	if i := strings.LastIndex(msg, ": "); i >= 0 && strings.Contains(msg[:i], "->") {
		head, _, _ := strings.Cut(msg, ": ")
		return head + ": " + msg[i+2:]
	}
	return msg
}
