package manage

import (
	"bytes"
	"context"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestAdditionalTestLinksAndPayloadIntegrity(t *testing.T) {
	p := externalPlan{Host: "192.0.2.1", Peer: "192.0.2.2", Token: ctNewSecret(), Coord: 12345}
	a, e := parseExternalTestLink(externalTestLink(p))
	if e != nil || a.Host != p.Host || a.Coord != p.Coord || a.Tok != p.Token {
		t.Fatal(a, e)
	}
	if _, e = parseExternalTestLink("bk://et.e30"); e == nil {
		t.Fatal("invalid link accepted")
	}
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer l.Close()
	go ctServeTCPEcho(l)
	c := &connTestCase{entry: l.Addr().(*net.TCPAddr).Port}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if e = externalPayloadSweep(ctx, c); e != nil {
		t.Fatal(e)
	}
	u, e := net.ListenPacket("udp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer u.Close()
	go ctServeUDPEcho(u)
	c = &connTestCase{entry: u.LocalAddr().(*net.UDPAddr).Port, udp: true}
	if e = externalPayloadSweep(ctx, c); e != nil {
		t.Fatal(e)
	}
}
func TestAllAdditionalTunnelTypesAppearInBothMenuAndTestCatalog(t *testing.T) {
	for _, kind := range []string{"gre", "l2tp-ip", "l2tp-udp", "awg", "ssh", "ssh-reverse", "rgt-tcp", "rgt-udp", "rgt-direct", "paqet", "alghadir"} {
		if _, ok := externaltunnel.Find(kind); !ok {
			t.Fatal(kind)
		}
	}
}
func TestAdditionalTableNeverTreatsMissingDependencyAsSuccess(t *testing.T) {
	s := ConnTestTable([]ConnTestResult{{Kind: "extra", Transport: "awg", Status: ctSkipped, Detail: "missing awg"}})
	if !strings.Contains(s, "SKIPPED") || !strings.Contains(s, "missing awg") {
		t.Fatal(s)
	}
	if bytes.Contains([]byte(s), []byte("🟢")) {
		t.Fatal(s)
	}
}

func TestAdditionalSetupFailuresAreNotMistakenForNetworkFiltering(t *testing.T) {
	table := additionalTestTable([]ConnTestResult{{Kind: "extra", Transport: "l2tp-ip", Status: ctDown, Detail: "kernel module unavailable"}})
	if !strings.Contains(table, "kernel module unavailable") || strings.Contains(table, "filtered for every transport") {
		t.Fatal(table)
	}
}

func TestAdditionalPeerRejectsPlansThatCouldChangeUnrelatedNetworkState(t *testing.T) {
	makePlan := func() externalPlan {
		s := externaltunnel.New("xt-sample-gre-tcp", "gre", "iran")
		s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
		s.Listen, s.Target = "127.0.0.1:21001", "127.0.0.1:21002"
		s.SourcePort = 21003
		s.SSHKey, s.KnownHosts = "", ""
		return externalPlan{Host: s.LocalIP, Peer: s.PeerIP, TCP: 21002, UDP: 21004, Token: ctNewSecret(), Cases: []externaltunnel.Spec{s}}
	}
	if e := validateExternalPeerPlan(makePlan()); e != nil {
		t.Fatal(e)
	}
	for _, alter := range []func(*externalPlan){
		func(p *externalPlan) { p.Cases[0].Name = "permanent-gre-tcp" },
		func(p *externalPlan) { p.Cases[0].Side = "kharej" },
		func(p *externalPlan) { p.Cases[0].PeerIP = "192.0.2.99" },
		func(p *externalPlan) { p.Cases[0].Target = "127.0.0.1:22" },
		func(p *externalPlan) { p.Cases[0].IranIP = "10.0.0.1/0"; p.Cases[0].KharejIP = "10.0.0.2/0" },
		func(p *externalPlan) { p.Cases[0].WAN = "eth0" },
		func(p *externalPlan) { p.Cases[0].Backend = "192.0.2.100:80" },
		func(p *externalPlan) { p.Cases = append(p.Cases, p.Cases[0]) },
	} {
		p := makePlan()
		alter(&p)
		if validateExternalPeerPlan(p) == nil {
			t.Fatal("unsafe test plan accepted")
		}
	}
}

func TestAdditionalCoordinatorWaitsForAuthenticatedPreparationReport(t *testing.T) {
	used := map[int]bool{}
	token := ctNewSecret()
	c, e := startExternalCoordinator(ctPickPort(used, true), token)
	if e != nil {
		t.Fatal(e)
	}
	defer c.close()
	if got := c.answer("hello "+token, "192.0.2.2"); got != "update" {
		t.Fatal(got)
	}
	select {
	case <-c.joined:
		t.Fatal("started before preparation")
	default:
	}
	report := externalIssueReport(map[string]string{"paqet-tcp": "missing libpcap"})
	if got := c.answer("ready invalid-token "+report, "192.0.2.2"); got != "" {
		t.Fatal(got)
	}
	if got := c.answer("ready "+token+" "+report, "192.0.2.2"); got != "ok" {
		t.Fatal(got)
	}
	select {
	case <-c.joined:
	default:
		t.Fatal("valid ready did not join")
	}
	issues := c.peerIssues()
	if issues["paqet-tcp"] != "missing libpcap" {
		t.Fatal(issues)
	}
	issues["paqet-tcp"] = "altered"
	if c.peerIssues()["paqet-tcp"] != "missing libpcap" {
		t.Fatal("report not copied")
	}
}
func TestAdditionalBothServersReceivePreparationFailuresWithoutSkippedRows(t *testing.T) {
	s, e := startExternalTestReady("127.0.0.1", "127.0.0.2", []string{"rgt-tcp"}, nil, nil, map[string]string{"rgt-tcp": "Iran download failed"})
	if e != nil {
		t.Fatal(e)
	}
	defer s.close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := make(chan []ConnTestResult, 1)
	errc := make(chan error, 1)
	go func() {
		rows, _, e := runExternalKharejReady(ctx, externalTestLink(s.plan), io.Discard, nil, func(context.Context, []externaltunnel.Spec, io.Writer) map[string]string {
			return map[string]string{"rgt-tcp": "Kharej download failed"}
		})
		result <- rows
		errc <- e
	}()
	select {
	case <-s.coord.joined:
	case <-ctx.Done():
		t.Fatal("peer not ready")
	}
	rows := s.run(ctx, nil)
	peer := <-result
	if e := <-errc; e != nil {
		t.Fatal(e)
	}
	for _, list := range [][]ConnTestResult{rows, peer} {
		if len(list) != 1 || list[0].Status != ctDown || list[0].Tried != 0 || !strings.Contains(list[0].Detail, "Iran download failed") || !strings.Contains(list[0].Detail, "Kharej download failed") {
			t.Fatal(list)
		}
	}
	table := additionalTestTable(rows)
	if strings.Contains(table, "SKIPPED") || !strings.Contains(table, "SETUP-FAIL") {
		t.Fatal(table)
	}
}
