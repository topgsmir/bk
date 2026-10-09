package externaltunnel

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixture(kind string) Spec {
	s := New("sample", kind, "iran")
	s.LocalIP = "192.0.2.1"
	s.PeerIP = "192.0.2.2"
	return s
}
func TestPairedSettingsAndKeys(t *testing.T) {
	for _, k := range Kinds {
		t.Run(k.ID, func(t *testing.T) {
			s := fixture(k.ID)
			link, e := s.Link()
			if e != nil {
				t.Fatal(e)
			}
			peer, e := ParseLink(link)
			if e != nil {
				t.Fatal(e)
			}
			if peer.Side != "kharej" || peer.LocalIP != s.PeerIP || peer.PeerIP != s.LocalIP || peer.Secret != s.Secret || peer.ID != s.ID {
				t.Fatalf("unpaired settings: %+v", peer)
			}
			if peer.Mirror().AWGPublic("iran") != s.AWGPublic("iran") {
				t.Fatal("keys disagree")
			}
			if peer.Interface() == s.Interface() {
				t.Fatal("paired interfaces must not collide in shared userspace socket directories")
			}
		})
	}
}
func TestUntrustedSettingsRejected(t *testing.T) {
	for _, change := range []func(*Spec){func(s *Spec) { s.Name = "../../etc/rc.local" }, func(s *Spec) { s.PeerIP = "1.2.3.999" }, func(s *Spec) { s.LocalIP = "1.2.3.4;reboot" }, func(s *Spec) { s.Listen = "0.0.0.0:70000" }, func(s *Spec) { s.Secret = "$(reboot)" }, func(s *Spec) { s.KharejIP = s.IranIP }, func(s *Spec) { s.ID = 0 }, func(s *Spec) { s.Protocol = "icmp" }, func(s *Spec) { s.WAN = "eth0;reboot" }, func(s *Spec) { s.RouterMAC = "bogus" }} {
		s := fixture("gre")
		change(&s)
		if s.Validate() == nil {
			t.Fatalf("accepted %+v", s)
		}
	}
	for _, raw := range []string{"bad", "bk://e.", "bk://e.e30"} {
		if _, e := ParseLink(raw); e == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
}
func TestL2TPEncapsulationAndScopedRemoval(t *testing.T) {
	for _, kind := range []string{"l2tp-ip", "l2tp-udp"} {
		s := fixture(kind)
		up, down := s.KernelCommands()
		want := "ip"
		if kind == "l2tp-udp" {
			want = "udp"
		}
		if up[0].Args[8] != want {
			t.Fatalf("wrong encapsulation: %+v", up)
		}
		all := ""
		for _, c := range append(up, down...) {
			all += c.Program + " " + strings.Join(c.Args, " ") + "\n"
		}
		for _, forbidden := range []string{"rc.local", "flush", "-F", "crontab", "MASQUERADE", "POSTROUTING"} {
			if strings.Contains(all, forbidden) {
				t.Fatal(all)
			}
		}
		if !strings.Contains(all, "name "+s.Interface()) {
			t.Fatal("interface was not explicitly named")
		}
	}
}
func TestConfigurationsMatchOfficialCoreSchemas(t *testing.T) {
	s := fixture("rgt-tcp")
	if !strings.Contains(s.RGTConfig(), "[server.services.bkforward]") {
		t.Fatal(s.RGTConfig())
	}
	peer := s.Mirror()
	if !strings.Contains(peer.RGTConfig(), "[client.services.bkforward]") {
		t.Fatal(peer.RGTConfig())
	}
	s.Kind = "paqet"
	s.WAN = "eth0"
	s.RouterMAC = "aa:bb:cc:dd:ee:ff"
	if !strings.Contains(s.PaqetConfig(), "server:") || !strings.Contains(s.PaqetConfig(), "forward:") {
		t.Fatal(s.PaqetConfig())
	}
	if !strings.Contains(s.Mirror().PaqetConfig(), "listen:") {
		t.Fatal("missing server listen")
	}
	s.Kind = "awg"
	a, b := s.AWGConfig(), s.Mirror().AWGConfig()
	if !strings.Contains(a, s.AWGPublic("kharej")) || !strings.Contains(b, s.AWGPublic("iran")) {
		t.Fatal("wrong peer public key")
	}
}
func TestConfigSaveNeverOverwritesExistingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sample.json")
	s := fixture("gre")
	if e := Save(p, s); e != nil {
		t.Fatal(e)
	}
	if e := Save(p, s); e == nil {
		t.Fatal("overwrote existing config")
	}
	if got, e := Load(p); e != nil || got.Secret != s.Secret {
		t.Fatal(got, e)
	}
	st, _ := os.Stat(p)
	if st.Mode().Perm() != 0600 {
		t.Fatal(st.Mode())
	}
}
func TestArchiveExtractionDoesNotWriteArchivePaths(t *testing.T) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	f, _ := z.Create("../../rgt")
	f.Write([]byte("core"))
	z.Close()
	got, e := extractCore(b.Bytes(), "rgt")
	if e != nil || string(got) != "core" {
		t.Fatal(string(got), e)
	}
	b.Reset()
	g := gzip.NewWriter(&b)
	tw := tar.NewWriter(g)
	tw.WriteHeader(&tar.Header{Name: "../paqet_linux_amd64", Mode: 0700, Size: 4})
	tw.Write([]byte("core"))
	tw.Close()
	g.Close()
	got, e = extractCore(b.Bytes(), "paqet")
	if e != nil || string(got) != "core" {
		t.Fatal(string(got), e)
	}
}
func TestRelayMovesBytesAndStopsIdleStreams(t *testing.T) {
	a, b := net.Pipe()
	c, d := net.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { relay(ctx, b, c); close(done) }()
	payload := []byte("real tunnel payload")
	go a.Write(payload)
	buf := make([]byte, len(payload))
	if _, e := d.Read(buf); e != nil || !bytes.Equal(buf, payload) {
		t.Fatal(e)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("idle relay leaked")
	}
	a.Close()
	d.Close()
}
func TestUnsupportedArchitectureIsExplicit(t *testing.T) {
	if _, e := CoreAsset("rgt-tcp", "arm64"); e == nil {
		t.Fatal("pretended unsupported RGT runs on arm64")
	}
	for _, k := range []string{"paqet", "rgt-tcp", "alghadir"} {
		a, e := CoreAsset(k, "amd64")
		if e != nil || len(a.SHA256) != 64 || !strings.HasPrefix(a.URL, "https://") {
			t.Fatal(a, e)
		}
	}
}
