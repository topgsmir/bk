package manage

import (
	"bytes"
	"context"
	"github.com/topgsmir/bk/internal/externaltunnel"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	for _, kind := range []string{"gre", "l2tp-ip", "l2tp-udp", "awg", "ssh", "rgt-tcp", "rgt-udp", "rgt-direct", "paqet", "alghadir"} {
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
