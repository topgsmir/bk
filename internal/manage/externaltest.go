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
	BatchSize         int `json:"batch_size,omitempty"`
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
	externalConnectionTestCatalog("Additional Tunnel Connection Test", externaltunnel.Kinds)
}
func Package3ConnectionTest() {
	externalConnectionTestCatalog("Test Tunnel Package 3", externaltunnel.Package3Kinds)
}
func externalConnectionTestCatalog(title string, catalog []externaltunnel.Kind) {
	tui.Clear()
	tui.Title(title)
	side := tui.ChooseOpt("This Server Is", []tui.Option{{Title: "Iran", Desc: "start temporary tunnels and copy the link"}, {Title: "Kharej", Desc: "paste the Iran test link"}})
	if side == 0 {
		externalIranMenuCatalog(catalog)
	} else if side == 1 {
		externalKharejMenu(tui.Prompt("bk://et. Test Link: "))
	}
}

const externalJoinWait = 30 * time.Minute

func externalIranMenuCatalog(catalog []externaltunnel.Kind) {
	host := strings.TrimSpace(tui.PromptDefault("Iran IPv4 Assigned To Its NIC", linkHost()))
	peer := strings.TrimSpace(tui.Prompt("Kharej IPv4 Assigned To Its NIC: "))
	if net.ParseIP(host).To4() == nil || net.ParseIP(peer).To4() == nil {
		tui.Error("Both IPv4 addresses are required.")
		tui.PressEnter()
		return
	}
	ctx, cancel := connTestContext()
	defer cancel()
	if package3Catalog(catalog) {
		if pick := tui.ChooseOpt("Package 3 Test", []tui.Option{{Title: "All Package 3 Methods", Desc: "full matrix; runs four cases at a time"}, {Title: "Select A Family", Desc: "Dagger, Solarpass, Backhaul or Eylan"}}); pick == 1 {
			var ok bool
			catalog, ok = choosePackage3Family(catalog)
			if !ok {
				return
			}
		} else if pick < 0 {
			return
		}
	}
	opts := []tui.Option{{Title: "All Selected Methods"}}
	for _, k := range catalog {
		opts = append(opts, tui.Option{Title: k.Title, Desc: k.Requirement})
	}
	selection := tui.ChooseOpt("Test Methods", opts)
	if selection < 0 {
		return
	}
	var kinds []string
	if selection > 0 {
		kinds = []string{catalog[selection-1].ID}
	}
	if selection == 0 {
		for _, k := range catalog {
			kinds = append(kinds, k.ID)
		}
	}
	ipv6 := [2]string{}
	for _, kind := range kinds {
		if kind == "d3-dc6" {
			ipv6[0] = strings.TrimSpace(tui.Prompt("Iran Assigned IPv6 (Blank = Report Unavailable): "))
			ipv6[1] = strings.TrimSpace(tui.Prompt("Kharej Assigned IPv6 (Blank = Report Unavailable): "))
			break
		}
	}
	var direct, reverse *externaltunnel.Spec
	for _, kind := range kinds {
		if kind == "ssh" {
			configured := promptExternalSSH("ssh", "iran", host, peer)
			direct = &configured
		}
		if kind == "ssh-reverse" {
			configured := externaltunnel.New("ssh-test", "ssh-reverse", "iran")
			configured.Port = externalNumber("Iran SSH Server Port (Kharej Connects Here)", 22)
			configured.SSHUser = tui.PromptDefault("Iran SSH User", configured.SSHUser)
			reverse = &configured
		}
	}
	tui.Info("Preparing Selected Methods On Both Servers Before Probing. Missing Verified Cores Are Installed Automatically; SSH May Ask For Login And Host Verification.")
	preparation := externaltunnel.PrepareDependencies(ctx, kinds, os.Stdout)
	if direct != nil {
		if e := ensureExternalSSH(ctx, *direct, os.Stdout); e != nil {
			preparation["ssh"] = e.Error()
		}
	}
	if ctx.Err() != nil {
		return
	}
	s, e := startExternalTestReadyIPv6(host, peer, kinds, direct, reverse, preparation, ipv6)
	if e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	defer s.close()
	path := "Additional tunnels"
	if package3Catalog(catalog) {
		path = "Test Tunnel Package 3"
	}
	tui.Info("Kharej: bk -> 0 -> " + path + " -> Kharej. Use The Same Updated Version On Both Servers. Copy This Secret Test Link:")
	fmt.Println(externalTestLink(s.plan))
	tui.Info("Waiting For Kharej To Prepare Dependencies And SSH (Up To 30 Minutes). Ctrl+C Stops.")
	select {
	case <-s.coord.joined:
	case <-ctx.Done():
		return
	case <-time.After(externalJoinWait):
		tui.Error("Kharej did not finish preparation and join.")
		tui.PressEnter()
		return
	}
	board := newCTBoard(os.Stdout, peer)
	rows := s.run(ctx, board.set)
	finishExternalBoard(board, rows, s.best)
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
	coord   *externalCoordinator
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
	return startExternalTestReady(host, peer, kinds, sshSettings, nil, nil)
}
func startExternalTestReady(host, peer string, kinds []string, sshSettings, reverseSettings *externaltunnel.Spec, preparation map[string]string) (*externalTest, error) {
	return startExternalTestReadyIPv6(host, peer, kinds, sshSettings, reverseSettings, preparation, [2]string{})
}
func startExternalTestReadyIPv6(host, peer string, kinds []string, sshSettings, reverseSettings *externaltunnel.Spec, preparation map[string]string, ipv6 [2]string) (*externalTest, error) {
	dir, e := os.MkdirTemp("", "bk-additional-test-")
	if e != nil {
		return nil, e
	}
	s := &externalTest{dir: dir}
	fail := func(e error) (*externalTest, error) { s.close(); return nil, e }
	used := map[int]bool{}
	for _, kind := range kinds {
		if externaltunnel.L2TPIPsec(kind) {
			for _, port := range []int{500, 4500, 1701} {
				used[port] = true
			}
		}
	}
	for _, settings := range []*externaltunnel.Spec{sshSettings, reverseSettings} {
		if settings != nil {
			used[settings.Port] = true
		}
	}
	s.plan = externalPlan{Version: 2, Host: host, Peer: peer, Token: ctNewSecret(), Coord: ctPickPort(used, true), TCP: ctPickPort(used, true), UDP: ctPickPort(used, true), Until: time.Now().Add(externalJoinWait + 5*time.Minute + connTestSlack).Unix()}
	for _, kind := range kinds {
		if externaltunnel.Dagger(kind) || externaltunnel.Solarpass(kind) || externaltunnel.Backhaul(kind) || externaltunnel.Eylan(kind) {
			s.plan.Version, s.plan.BatchSize = 3, 4
			s.plan.Until = time.Now().Add(externalJoinWait + 45*time.Minute + connTestSlack).Unix()
			break
		}
	}
	s.coord, e = startExternalCoordinator(s.plan.Coord, s.plan.Token)
	if e != nil {
		return fail(e)
	}
	if kinds == nil {
		for _, k := range externaltunnel.Kinds {
			kinds = append(kinds, k.ID)
		}
	}
	if s.plan.BatchSize > 0 {
		s.coord.enableWaves()
	}
	id := randomToken(6)
	index := 0
	baseID := 500000 + int(time.Now().UnixNano()%1000000)
	for _, kind := range kinds {
		protocols := []string{"tcp"}
		if externaltunnel.BothProtocols(kind) {
			protocols = append(protocols, "udp")
		}
		if externaltunnel.UDPOnly(kind) {
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
			spec.LocalIPv6, spec.PeerIPv6 = ipv6[0], ipv6[1]
			spec.Port = ctPickPort(used, true)
			spec.SourcePort = ctPickPort(used, true)
			if externaltunnel.Solarpass(kind) || externaltunnel.Backhaul(kind) || externaltunnel.Eylan(kind) {
				spec.Port = ctPickSolarBlock(used, 2)
				spec.SourcePort = ctPickSolarBlock(used, 4)
			}
			if externaltunnel.L2TPIPsec(kind) {
				spec.Port = 1701
			}
			spec.ID = baseID + index
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
			if externaltunnel.SSH(kind) {
				spec.Port = 22
				settings := sshSettings
				if kind == "ssh-reverse" {
					settings = reverseSettings
				}
				if settings != nil {
					spec.Port, spec.SSHUser, spec.SSHKey, spec.KnownHosts = settings.Port, settings.SSHUser, settings.SSHKey, settings.KnownHosts
				}
			}
			if why := preparation[kind]; why != "" {
				c.skip = why
			} else if e := externaltunnel.Check(spec); e != nil {
				c.skip = e.Error()
			} else if s.plan.BatchSize == 0 {
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
	stopAdditionalEngines(s.engines)
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
	if s.plan.BatchSize > 0 {
		return s.runPackage3Waves(ctx, progress)
	}
	return s.runMeasurements(ctx, progress, true)
}
func (s *externalTest) runMeasurements(ctx context.Context, progress func(int, ConnTestResult), publish bool) []ConnTestResult {
	rows := make([]ConnTestResult, len(s.cases))
	var mu sync.Mutex
	peerIssues := s.coord.peerIssues()
	for i, c := range s.cases {
		rows[i] = ConnTestResult{Kind: c.kind, Transport: c.tr, Status: ctTesting, Total: connTestSoak}
		if c.skip != "" {
			rows[i].Status, rows[i].Detail = ctDown, "SETUP FAILED (Iran): "+c.skip
		}
		if why := peerIssues[c.tr]; why != "" {
			rows[i].Status = ctDown
			rows[i].Detail += joinExternalDetail(rows[i].Detail, "SETUP FAILED (Kharej): "+why)
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
		if c.skip != "" || peerIssues[c.tr] != "" {
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
			if r.Status == ctOK && c.udp && strings.HasPrefix(c.tr, "s3-") {
				if err := externalSolarUDPConcurrency(ctx, c); err != nil {
					r.Status, r.Detail = ctUnstable, "concurrent UDP test: "+err.Error()
					report(r)
				}
			}
			if r.Status != ctOK {
				for _, engine := range s.engines {
					if engine.name == c.name {
						if why := externalEngineDiagnostic(engine); why != "" {
							r.Detail += joinExternalDetail(r.Detail, "Iran engine: "+why)
							report(r)
						}
					}
				}
			}
			for _, engine := range s.engines {
				if engine.name == c.name {
					if exited, why := engine.exited(); exited {
						if r.Status == ctOK {
							r.Status = ctUnstable
						}
						r.Detail = "local setup/engine: " + why
						report(r)
					}
				}
			}
		}(i, c)
	}
	wg.Wait()
	s.coord.finishMeasurements()
	select {
	case <-s.coord.finalReport:
	case <-ctx.Done():
	case <-time.After(10 * time.Second):
	}
	peerIssues = s.coord.peerIssues()
	mu.Lock()
	for i, c := range s.cases {
		if why := peerIssues[c.tr]; why != "" {
			if rows[i].Status == ctOK {
				rows[i].Status = ctUnstable
			}
			if !strings.Contains(rows[i].Detail, why) {
				rows[i].Detail += joinExternalDetail(rows[i].Detail, "Kharej: "+why)
			}

		}
	}
	mu.Unlock()
	if progress != nil {
		for i, r := range rows {
			progress(i, r)
		}
	}
	s.best = ctComputeBest(rows, s.cases, 0, s.dir)
	if publish {
		s.coord.publish(rows, s.best)
	}
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
				return fmt.Errorf("%d-byte datagram: %w", size, err)
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
	rows, best, e := runExternalKharejReady(ctx, raw, os.Stdout, func(rows []ConnTestResult) {
		if board == nil {
			board = newCTBoard(os.Stdout, "Iran")
		}
		board.setAll(rows)
	}, prepareExternalPeer)
	if e != nil {
		if board != nil {
			board.abandon()
		}
		tui.Error(e.Error())
	} else {
		if board == nil {
			board = newCTBoard(os.Stdout, "Iran")
		}
		finishExternalBoard(board, rows, best)
	}
	tui.PressEnter()
}
func runExternalKharej(ctx context.Context, raw string, out io.Writer, live func([]ConnTestResult)) ([]ConnTestResult, ConnTestBest, error) {
	return runExternalKharejReady(ctx, raw, out, live, nil)
}
func runExternalKharejReady(ctx context.Context, raw string, out io.Writer, live func([]ConnTestResult), prepare func(context.Context, []externaltunnel.Spec, io.Writer) map[string]string) ([]ConnTestResult, ConnTestBest, error) {
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
	if !fetched || (p.Version != 2 && p.Version != 3) || len(p.Cases) > 256 || p.Host != a.Host || p.Coord != a.Coord || p.Token != a.Tok || time.Now().Unix() > p.Until || p.Until > time.Now().Add(2*time.Hour).Unix() {
		return nil, ConnTestBest{}, fmt.Errorf("additional test configuration missing, expired or incompatible; update bk on both servers")
	}
	if e := validateExternalPeerPlan(p); e != nil {
		return nil, ConnTestBest{}, e
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
	defer func() { stopAdditionalEngines(engines) }()
	local := map[string]string{}
	preparation := map[string]string{}
	if prepare != nil {
		preparation = prepare(ctx, p.Cases, out)
	}
	if p.Version == 3 {
		return runPackage3Peer(ctx, a, p, dir, preparation, out, live)
	}
	for _, spec := range p.Cases {
		tr := spec.Kind
		if externaltunnel.BothProtocols(spec.Kind) {
			tr += "-" + spec.Protocol
		}
		key, known := spec.SSHKey, spec.KnownHosts
		spec = spec.Mirror()
		if externaltunnel.SSHInitiator(spec) && key != "" {
			spec.SSHKey, spec.KnownHosts = key, known
		}
		spec.Secret = ctCaseToken(p.Token, "extra", tr) + ctCaseToken(p.Token, "extra-key", tr)
		if why := preparation[spec.Kind]; why != "" {
			local[tr] = why
			continue
		}
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
		r, e := ctAsk(a.Host, a.Coord, "ready "+a.Tok+" "+externalIssueReport(local))
		if e == nil && r == "ok" {
			break
		}
		ctSleep(ctx, time.Second)
	}
	for ctx.Err() == nil {
		for _, engine := range engines {
			if exited, why := engine.exited(); exited {
				for _, spec := range p.Cases {
					if spec.Name == engine.name {
						local[externalTransport(spec)] = why
					}
				}
			}
		}
		if len(local) > 0 {
			_, _ = ctAsk(a.Host, a.Coord, "diagnostic "+a.Tok+" "+externalIssueReport(local))
		}
		r, e := ctAsk(a.Host, a.Coord, "result "+a.Tok)
		if e == nil && strings.HasPrefix(r, "live ") && live != nil {
			b, e := unGzipB64(strings.TrimPrefix(r, "live "))
			var rows []ConnTestResult
			if e == nil && json.Unmarshal(b, &rows) == nil {
				live(rows)
			}
		}
		if e == nil && strings.HasPrefix(r, "measured ") {
			body, err := unGzipB64(strings.TrimPrefix(r, "measured "))
			var rows []ConnTestResult
			if err == nil && json.Unmarshal(body, &rows) == nil {
				for _, row := range rows {
					if row.Status != ctOK {
						for _, engine := range engines {
							if strings.HasSuffix(engine.name, "-"+row.Transport) {
								if why := externalEngineDiagnostic(engine); why != "" {
									local[row.Transport] = why
								}
							}
						}
					}
				}
				_, _ = ctAsk(a.Host, a.Coord, "final "+a.Tok+" "+externalIssueReport(local))
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
				if why := local[row.Transport]; why != "" && row.Status != ctOK && !strings.Contains(row.Detail, why) {
					v.Results[i].Detail += "; kharej: " + why
				}
				if row.Status != ctOK {
					for _, engine := range engines {
						if strings.HasSuffix(engine.name, "-"+row.Transport) {
							if exited, why := engine.exited(); exited && !strings.Contains(v.Results[i].Detail, why) {
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

// Reuse the original live board, but print only measurements applicable to
// these cores. Original bk presets, FEC and MSS recommendations do not apply.
func finishExternalBoard(board *ctBoard, rows []ConnTestResult, best ConnTestBest) {
	board.abandon()
	var w strings.Builder
	if board.tty && len(board.shown) > 0 {
		fmt.Fprintf(&w, "\033[%dF\033[J", len(board.shown))
	}
	w.WriteString("ADDITIONAL TUNNEL CONNECTION TEST - " + board.title + "\n" + ctRule() + "\n\n")
	w.WriteString(additionalTestTable(rows))
	if best.Transport != "" {
		speed := "UDP throughput not measured"
		if best.Mbps > 0 {
			speed = fmt.Sprintf("%.1f Mbps", best.Mbps)
		}
		fmt.Fprintf(&w, "\nBest measured method: %s (%s, %d ms RTT, %d ms jitter, %.1f%% loss).\n", best.Transport, speed, best.RTTms, best.JitterMs, best.LossPct)
	}
	_, _ = io.WriteString(board.out, w.String())
}
func additionalTestTable(rows []ConnTestResult) string {
	table := ConnTestTable(rows)
	for _, r := range rows {
		if strings.Contains(r.Detail, "SETUP FAILED") {
			table = strings.Replace(table, ctRow(r), strings.Replace(ctRow(r), "DOWN        ", "SETUP-FAIL  ", 1), 1)
		}
	}
	table = strings.ReplaceAll(table, "Nothing carried traffic steadily between these two servers. The path is\nfiltered for every transport tried; another Iran or kharej server (another\nprovider, another IP) is the fix, not a setting.\n", "No additional method passed. Check the setup/dependency errors below and\nverify addresses, routing and network firewalls before trying another path.\n")
	for _, r := range rows {
		if r.Status != ctOK && r.Status != ctSkipped && r.Detail != "" {
			table += fmt.Sprintf("%s %s: %s\n", ctName(r.Transport), r.Status, r.Detail)
		}
	}
	return table
}

func stopAdditionalEngines(engines []*ctEngine) { stopAdditionalEnginesWithin(engines, 20*time.Second) }
func stopAdditionalEnginesWithin(engines []*ctEngine, timeout time.Duration) {
	for _, engine := range engines {
		_ = engine.cmd.Process.Signal(syscall.SIGTERM)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for _, engine := range engines {
		select {
		case <-engine.done:
		case <-ctx.Done():
			_ = syscall.Kill(-engine.cmd.Process.Pid, syscall.SIGKILL)
			<-engine.done
		}
	}
}

// A received test plan can configure only throwaway test resources and the
// loopback echo endpoints. Reject unrelated services, broad routes and local
// machine overrides before starting any process or listener as root.
func validateExternalPeerPlan(p externalPlan) error {
	if (p.Version == 3 && p.BatchSize != 4) || (p.Version == 2 && p.BatchSize != 0) {
		return fmt.Errorf("invalid test batch size")
	}
	if p.TCP < 1024 || p.TCP > 65535 || p.UDP < 1024 || p.UDP > 65535 || p.TCP == p.UDP || len(p.Cases) == 0 {
		return fmt.Errorf("invalid additional-test echo endpoints")
	}
	names := map[string]bool{}
	for _, spec := range p.Cases {
		tr := spec.Kind
		if externaltunnel.BothProtocols(spec.Kind) {
			tr += "-" + spec.Protocol
		}
		if spec.Side != "iran" || spec.LocalIP != p.Host || spec.PeerIP != p.Peer || !strings.HasPrefix(spec.Name, "xt-") || !strings.HasSuffix(spec.Name, "-"+tr) || names[spec.Name] {
			return fmt.Errorf("invalid temporary tunnel identity")
		}
		names[spec.Name] = true
		port := p.TCP
		if spec.Protocol == "udp" {
			port = p.UDP
		}
		host, entry, e := net.SplitHostPort(spec.Listen)
		n, e2 := strconv.Atoi(entry)
		if e != nil || e2 != nil || host != "127.0.0.1" || n < 1024 || n > 65535 || spec.Target != net.JoinHostPort("127.0.0.1", strconv.Itoa(port)) {
			return fmt.Errorf("test plan may forward only to its loopback echo service")
		}
		if spec.Backend != "" || spec.WAN != "" || spec.RouterMAC != "" || spec.SSHKey != "" || spec.KnownHosts != "" {
			return fmt.Errorf("test plan contains a machine-local override")
		}
		if !externaltunnel.SSH(spec.Kind) && (spec.Port < 1024 || spec.SourcePort < 1024) {
			return fmt.Errorf("test transport ports must be unprivileged")
		}
		spec = spec.Mirror()
		spec.Secret = ctCaseToken(p.Token, "extra", tr) + ctCaseToken(p.Token, "extra-key", tr)
		if e := spec.Validate(); e != nil {
			return e
		}
		if externaltunnel.Layer3(spec.Kind) {
			_, subnet, _ := net.ParseCIDR(spec.TunnelIP())
			prefix, _ := subnet.Mask.Size()
			if prefix != 30 {
				return fmt.Errorf("temporary tests require a /30 tunnel subnet")
			}
		}
	}
	return nil
}

func joinExternalDetail(existing, extra string) string {
	if existing != "" {
		return "; " + extra
	}
	return extra
}
func externalTransport(s externaltunnel.Spec) string {
	tr := s.Kind
	if externaltunnel.BothProtocols(s.Kind) {
		tr += "-" + s.Protocol
	}
	return tr
}

// Keep the actual failure, including an authenticated SSH server refusing a
// forwarding request. Ignore informational startup output on successful rows.
func externalEngineDiagnostic(engine *ctEngine) string {
	if exited, why := engine.exited(); exited {
		return why
	}
	file, e := os.Open(engine.log)
	if e != nil {
		return ""
	}
	defer file.Close()
	info, e := file.Stat()
	if e != nil {
		return ""
	}
	offset := info.Size() - 4096
	if offset < 0 {
		offset = 0
	}
	_, _ = file.Seek(offset, io.SeekStart)
	body, _ := io.ReadAll(io.LimitReader(file, 4096))
	lines := strings.Split(string(body), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		lower := strings.ToLower(line)
		if strings.Contains(lower, "failed") || strings.Contains(lower, "error") || strings.Contains(lower, "refused") {
			if len(line) > 1024 {
				line = line[:1024]
			}
			return line
		}
	}
	return ""
}

func ctPickSolarBlock(used map[int]bool, count int) int {
	for {
		port := ctPickPort(used, true)
		if port+count > 65535 {
			continue
		}
		free := true
		for n := 1; n < count; n++ {
			if used[port+n] {
				free = false
				break
			}
			tcp, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(port+n)))
			if err != nil {
				free = false
				break
			}
			udp, err := net.ListenPacket("udp", net.JoinHostPort("0.0.0.0", strconv.Itoa(port+n)))
			tcp.Close()
			if err != nil {
				free = false
				break
			}
			udp.Close()
		}
		if free {
			for n := 1; n < count; n++ {
				used[port+n] = true
			}
			return port
		}
	}
}

func externalSolarUDPConcurrency(ctx context.Context, c *connTestCase) error {
	results := make(chan error, 3)
	for client := 0; client < 3; client++ {
		go func(client int) {
			var dialer net.Dialer
			conn, err := dialer.DialContext(ctx, "udp4", net.JoinHostPort("127.0.0.1", strconv.Itoa(c.entry)))
			if err != nil {
				results <- err
				return
			}
			defer conn.Close()
			for packet := 0; packet < 10; packet++ {
				body := make([]byte, 512)
				if _, err = rand.Read(body); err != nil {
					results <- err
					return
				}
				body[0], body[1] = byte(client), byte(packet)
				conn.SetDeadline(time.Now().Add(3 * time.Second))
				if _, err = conn.Write(body); err != nil {
					results <- err
					return
				}
				reply := make([]byte, 1500)
				n, err := conn.Read(reply)
				if err != nil {
					results <- fmt.Errorf("client %d packet %d: %w", client, packet, err)
					return
				}
				if !bytes.Equal(body, reply[:n]) {
					results <- fmt.Errorf("client %d packet %d payload mismatch", client, packet)
					return
				}
			}
			results <- nil
		}(client)
	}
	var failure error
	for i := 0; i < 3; i++ {
		if err := <-results; err != nil {
			failure = err
		}
	}
	return failure
}
