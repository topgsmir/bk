package manage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
)

const externalTestScheme = "bk://et."

type externalPlan struct {
	Version           int `json:"v"`
	Host, Peer, Token string
	Coord, TCP, UDP   int
	Until             int64
	Cases             []externaltunnel.Spec
}

func externalTestLink(p externalPlan) string {
	b, _ := json.Marshal(connTestAddr{Host: p.Host, Coord: p.Coord, Tok: p.Token})
	return externalTestScheme + base64.RawURLEncoding.EncodeToString(b)
}
func parseExternalTestLink(raw string) (connTestAddr, error) {
	var a connTestAddr
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, externalTestScheme) || len(raw) > 4096 {
		return a, fmt.Errorf("paste the full bk://et. test link")
	}
	b, e := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(raw, externalTestScheme))
	if e != nil {
		return a, e
	}
	if e = json.Unmarshal(b, &a); e != nil {
		return a, e
	}
	if net.ParseIP(a.Host).To4() == nil || a.Coord < 1 || a.Coord > 65535 || len(a.Tok) != 16 {
		return a, fmt.Errorf("invalid additional-tunnel test link")
	}
	return a, nil
}
func ExternalConnectionTest() {
	tui.Clear()
	tui.Title("Additional Tunnel Connection Test")
	side := tui.ChooseOpt("This Server Is", []tui.Option{{Title: "Iran", Desc: "start temporary tunnels and copy the link"}, {Title: "Kharej", Desc: "paste the Iran test link"}})
	if side == 0 {
		externalIranMenu()
	} else if side == 1 {
		externalKharejMenu(tui.Prompt("bk://et. Test Link: "))
	}
}
func externalIranMenu() {
	host := strings.TrimSpace(tui.PromptDefault("Iran IPv4 Assigned To Its NIC", linkHost()))
	peer := strings.TrimSpace(tui.Prompt("Kharej IPv4 Assigned To Its NIC: "))
	if net.ParseIP(host).To4() == nil || net.ParseIP(peer).To4() == nil {
		tui.Error("Both IPv4 addresses are required.")
		tui.PressEnter()
		return
	}
	tui.Info("Uses Existing Dependencies. Missing Cores Or Kernel Support Are Reported As SKIPPED/Setup Errors.")
	ctx, cancel := connTestContext()
	defer cancel()
	opts := []tui.Option{{Title: "All Additional Methods"}}
	for _, k := range externaltunnel.Kinds {
		opts = append(opts, tui.Option{Title: k.Title, Desc: k.Requirement})
	}
	selection := tui.ChooseOpt("Test Methods", opts)
	if selection < 0 {
		return
	}
	var kinds []string
	if selection > 0 {
		kinds = []string{externaltunnel.Kinds[selection-1].ID}
	}
	var sshSettings *externaltunnel.Spec
	if selection == 0 || (len(kinds) == 1 && kinds[0] == "ssh") {
		configured := externaltunnel.New("ssh-test", "ssh", "iran")
		configured.Port = externalNumber("Kharej SSH Port", 22)
		configured.SSHUser = tui.PromptDefault("Kharej SSH User", configured.SSHUser)
		configured.SSHKey = tui.PromptDefault("SSH Private Key Path", configured.SSHKey)
		configured.KnownHosts = tui.PromptDefault("Verified Known Hosts File", configured.KnownHosts)
		sshSettings = &configured
	}
	s, e := startExternalTestOptions(host, peer, kinds, sshSettings)
	if e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	defer s.close()
	fmt.Println()
	tui.Info("Kharej: bk → 0 → Additional tunnels → Kharej. Copy This Secret Test Link:")
	fmt.Println(externalTestLink(s.plan))
	select {
	case <-s.coord.joined:
	case <-ctx.Done():
		return
	case <-time.After(connTestJoinWait):
		tui.Error("Kharej did not join.")
		tui.PressEnter()
		return
	}
	board := newCTBoard(os.Stdout, peer)
	rows := s.run(ctx, board.set)
	board.finish(rows, s.best)
	select {
	case <-s.coord.fetched:
	case <-ctx.Done():
	case <-time.After(20 * time.Second):
	}
	s.close()
	tui.PressEnter()
}

type externalTest struct {
	dir     string
	plan    externalPlan
	coord   *ctCoordinator
	engines []*ctEngine
	cases   []*connTestCase
	best    ConnTestBest
}

func startExternalEngine(dir string, s externaltunnel.Spec) (*ctEngine, error) {
	bin, e := connTestBinary()
	if e != nil {
		return nil, e
	}
	cfg := filepath.Join(dir, s.Name+".json")
	if e = externaltunnel.Save(cfg, s); e != nil {
		return nil, e
	}
	logPath := filepath.Join(dir, s.Name+".log")
	f, e := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return nil, e
	}
	cmd := exec.Command(bin, "external", "run", "-c", cfg)
	cmd.Stdout, cmd.Stderr = f, f
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if e = cmd.Start(); e != nil {
		f.Close()
		return nil, e
	}
	engine := &ctEngine{name: s.Name, cmd: cmd, log: logPath, done: make(chan struct{})}
	go func() { _ = cmd.Wait(); f.Close(); close(engine.done) }()
	return engine, nil
}

// Passing a kind list is used by the privileged integration test; nil covers
// every new kind without changing the original Connection Test's cases.
func startExternalTest(host, peer string, kinds []string) (*externalTest, error) {
	return startExternalTestOptions(host, peer, kinds, nil)
}
func startExternalTestOptions(host, peer string, kinds []string, sshSettings *externaltunnel.Spec) (*externalTest, error) {
	dir, e := os.MkdirTemp("", "bk-additional-test-")
	if e != nil {
		return nil, e
	}
	s := &externalTest{dir: dir}
	fail := func(e error) (*externalTest, error) { s.close(); return nil, e }
	used := map[int]bool{}
	s.plan = externalPlan{Version: 1, Host: host, Peer: peer, Token: ctNewSecret(), Coord: ctPickPort(used, true), TCP: ctPickPort(used, true), UDP: ctPickPort(used, true), Until: time.Now().Add(connTestJoinWait + 5*time.Minute + connTestSlack).Unix()}
	s.coord, e = startCTCoordinator(s.plan.Coord, s.plan.Token)
	if e != nil {
		return fail(e)
	}
	if kinds == nil {
		for _, k := range externaltunnel.Kinds {
			kinds = append(kinds, k.ID)
		}
	}
	id := randomToken(6)
	index := 0
	for _, kind := range kinds {
		protocols := []string{"tcp"}
		if externaltunnel.Layer3(kind) || kind == "paqet" {
			protocols = append(protocols, "udp")
		}
		if kind == "rgt-udp" {
			protocols = []string{"udp"}
		}
		for _, proto := range protocols {
			tr := kind
			if len(protocols) > 1 {
				tr += "-" + proto
			}
			c := &connTestCase{kind: "extra", tr: tr, name: "xt-" + id + "-" + tr, udp: proto == "udp", entry: ctPickPort(used, true)}
			spec := externaltunnel.New(c.name, kind, "iran")
			spec.LocalIP, spec.PeerIP = host, peer
			spec.Port = ctPickPort(used, true)
			spec.SourcePort = ctPickPort(used, true)
			spec.ID = 500000 + int(time.Now().UnixNano()%1000000) + index
			spec.Secret = ctCaseToken(s.plan.Token, "extra", tr) + ctCaseToken(s.plan.Token, "extra-key", tr)
			spec.Protocol = proto
			block := fmt.Sprintf("10.204.%d.", index+1)
			spec.IranIP, spec.KharejIP = block+"1/30", block+"2/30"
			spec.Listen = net.JoinHostPort("127.0.0.1", strconv.Itoa(c.entry))
			target := s.plan.TCP
			if c.udp {
				target = s.plan.UDP
			}
			spec.Target = net.JoinHostPort("127.0.0.1", strconv.Itoa(target))
			if kind == "ssh" {
				spec.Port = 22
				if sshSettings != nil {
					spec.Port, spec.SSHUser, spec.SSHKey, spec.KnownHosts = sshSettings.Port, sshSettings.SSHUser, sshSettings.SSHKey, sshSettings.KnownHosts
				}
			}
			if e := externaltunnel.Check(spec); e != nil {
				c.skip = e.Error()
			} else {
				engine, e := startExternalEngine(dir, spec)
				if e != nil {
					c.skip = e.Error()
				} else {
					s.engines = append(s.engines, engine)
				}
			}
			s.cases = append(s.cases, c)
			s.plan.Cases = append(s.plan.Cases, spec)
			index++
		}
	}
	public := s.plan
	public.Cases = append([]externaltunnel.Spec(nil), s.plan.Cases...)
	for i := range public.Cases {
		public.Cases[i].Secret = ""
		public.Cases[i].SSHKey = ""
		public.Cases[i].KnownHosts = ""
	}
	b, e := json.Marshal(public)
	if e != nil {
		return fail(e)
	}
	s.coord.setConfig(gzipB64(b))
	return s, nil
}
func (s *externalTest) close() {
	if s == nil {
		return
	}
	ctStopAll(s.engines)
	s.engines = nil
	if s.coord != nil {
		s.coord.close()
		s.coord = nil
	}
	if s.dir != "" {
		os.RemoveAll(s.dir)
		s.dir = ""
	}
}
func (s *externalTest) run(ctx context.Context, progress func(int, ConnTestResult)) []ConnTestResult {
	rows := make([]ConnTestResult, len(s.cases))
	var mu sync.Mutex
	for i, c := range s.cases {
		rows[i] = ConnTestResult{Kind: c.kind, Transport: c.tr, Status: ctTesting, Total: connTestSoak}
		if c.skip != "" {
			rows[i].Status, rows[i].Detail = ctSkipped, c.skip
		}
	}
	s.coord.setLive(func() []ConnTestResult { mu.Lock(); defer mu.Unlock(); return append([]ConnTestResult(nil), rows...) })
	if progress != nil {
		for i, r := range rows {
			progress(i, r)
		}
	}
	var wg sync.WaitGroup
	for i, c := range s.cases {
		if c.skip != "" {
			continue
		}
		wg.Add(1)
		go func(i int, c *connTestCase) {
			defer wg.Done()
			report := func(r ConnTestResult) {
				mu.Lock()
				rows[i] = r
				mu.Unlock()
				if progress != nil {
					progress(i, r)
				}
			}
			r := ctProbe(ctx, c, 0, report)
			// The unchanged probe checks a persistent connection, recovery, RTT and
			// a byte-identical 1 MiB TCP transfer. Also test packet sizes through the
			// actual tunnel, so a port that opens but drops payload cannot pass.
			if r.Status == ctOK {
				if e := externalPayloadSweep(ctx, c); e != nil {
					r.Status, r.Detail = ctUnstable, "payload-size test: "+e.Error()
					report(r)
				}
			}
			for _, engine := range s.engines {
				if engine.name == c.name {
					if exited, why := engine.exited(); exited && r.Status != ctOK {
						r.Detail = "local setup/engine: " + why
						report(r)
					}
				}
			}
		}(i, c)
	}
	wg.Wait()
	s.best = ctComputeBest(rows, s.cases, 0, s.dir)
	s.coord.publish(rows, s.best)
	return rows
}
func externalPayloadSweep(ctx context.Context, c *connTestCase) error {
	network := "tcp"
	if c.udp {
		network = "udp"
	}
	d := net.Dialer{Timeout: 3 * time.Second}
	conn, e := d.DialContext(ctx, network, net.JoinHostPort("127.0.0.1", strconv.Itoa(c.entry)))
	if e != nil {
		return e
	}
	defer conn.Close()
	sizes := []int{64, 512, 1200}
	if !c.udp {
		sizes = append(sizes, 16<<10)
	}
	for _, size := range sizes {
		b := make([]byte, size)
		if _, e = rand.Read(b); e != nil {
			return e
		}
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if _, e = conn.Write(b); e != nil {
			return fmt.Errorf("%d bytes: %w", size, e)
		}
		r := make([]byte, size)
		if c.udp {
			n, err := conn.Read(r)
			if err != nil {
				return err
			}
			if n != size {
				return fmt.Errorf("%d-byte datagram truncated to %d", size, n)
			}
		} else {
			if _, e = io.ReadFull(conn, r); e != nil {
				return e
			}
		}
		if !bytes.Equal(b, r) {
			return fmt.Errorf("%d-byte payload corrupted", size)
		}
	}
	return nil
}
func externalKharejMenu(raw string) {
	ctx, cancel := connTestContext()
	defer cancel()
	var board *ctBoard
	rows, best, e := runExternalKharej(ctx, raw, os.Stdout, func(rows []ConnTestResult) {
		if board == nil {
			board = newCTBoard(os.Stdout, "Iran")
		}
		board.setAll(rows)
	})
	if e != nil {
		if board != nil {
			board.abandon()
		}
		tui.Error(e.Error())
	} else {
		if board == nil {
			board = newCTBoard(os.Stdout, "Iran")
		}
		board.finish(rows, best)
	}
	tui.PressEnter()
}
func runExternalKharej(ctx context.Context, raw string, out io.Writer, live func([]ConnTestResult)) ([]ConnTestResult, ConnTestBest, error) {
	a, e := parseExternalTestLink(raw)
	if e != nil {
		return nil, ConnTestBest{}, e
	}
	var p externalPlan
	fetched := false
	deadline := time.Now().Add(time.Minute)
	for !fetched && ctx.Err() == nil && time.Now().Before(deadline) {
		r, e := ctAsk(a.Host, a.Coord, "config "+a.Tok)
		if e == nil && strings.HasPrefix(r, "config ") {
			b, e := unGzipB64(strings.TrimPrefix(r, "config "))
			if e == nil && json.Unmarshal(b, &p) == nil {
				fetched = true
				break
			}
		}
		ctSleep(ctx, time.Second)
	}
	if !fetched || p.Version != 1 || len(p.Cases) > 32 || p.Host != a.Host || p.Coord != a.Coord || p.Token != a.Tok || time.Now().Unix() > p.Until || p.Until > time.Now().Add(time.Hour).Unix() {
		return nil, ConnTestBest{}, fmt.Errorf("additional test configuration missing, damaged or expired")
	}
	ctx, cancel := context.WithDeadline(ctx, time.Unix(p.Until, 0))
	defer cancel()
	dir, e := os.MkdirTemp("", "bk-additional-test-")
	if e != nil {
		return nil, ConnTestBest{}, e
	}
	defer os.RemoveAll(dir)
	echoes, e := ctStartEchoes(ctx, ConnTestLink{TCP: p.TCP, UDP: p.UDP})
	if e != nil {
		return nil, ConnTestBest{}, e
	}
	defer echoes.close()
	var engines []*ctEngine
	defer func() { ctStopAll(engines) }()
	local := map[string]string{}
	for _, spec := range p.Cases {
		tr := spec.Kind
		if externaltunnel.Layer3(spec.Kind) || spec.Kind == "paqet" {
			tr += "-" + spec.Protocol
		}
		spec = spec.Mirror()
		spec.Secret = ctCaseToken(p.Token, "extra", tr) + ctCaseToken(p.Token, "extra-key", tr)
		if e = externaltunnel.Check(spec); e != nil {
			local[tr] = e.Error()
			continue
		}
		engine, e := startExternalEngine(dir, spec)
		if e != nil {
			local[tr] = e.Error()
			continue
		}
		engines = append(engines, engine)
	}
	fmt.Fprintf(out, "Started %d temporary additional tunnels. No services or permanent configuration are created.\n", len(engines))
	for ctx.Err() == nil {
		r, e := ctAsk(a.Host, a.Coord, "hello "+a.Tok)
		if e == nil && r == "ok" {
			break
		}
		ctSleep(ctx, time.Second)
	}
	for ctx.Err() == nil {
		r, e := ctAsk(a.Host, a.Coord, "result "+a.Tok)
		if e == nil && strings.HasPrefix(r, "live ") && live != nil {
			b, e := unGzipB64(strings.TrimPrefix(r, "live "))
			var rows []ConnTestResult
			if e == nil && json.Unmarshal(b, &rows) == nil {
				live(rows)
			}
		}
		if e == nil && strings.HasPrefix(r, "done ") {
			b, e := unGzipB64(strings.TrimPrefix(r, "done "))
			if e != nil {
				return nil, ConnTestBest{}, e
			}
			var v ctVerdict
			if e = json.Unmarshal(b, &v); e != nil {
				return nil, ConnTestBest{}, e
			}
			for i, row := range v.Results {
				if why := local[row.Transport]; why != "" && row.Status != ctOK {
					v.Results[i].Detail += "; kharej: " + why
				}
				if row.Status != ctOK {
					for _, engine := range engines {
						if strings.HasSuffix(engine.name, "-"+row.Transport) {
							if exited, why := engine.exited(); exited {
								v.Results[i].Detail += "; kharej engine: " + why
							}
						}
					}
				}
			}
			return v.Results, v.Best, nil
		}
		ctSleep(ctx, 2*time.Second)
	}
	return nil, ConnTestBest{}, ctx.Err()
}
