//go:build linux

package manage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// This test uses two real Linux network namespaces, the production binary and
// actual upstream cores. It is opt-in locally and mandatory in the addon CI.
func TestAdditionalTunnelsAcrossNamespaces(t *testing.T) {
	if os.Getenv("BK_LIVE_ADDONS") != "1" {
		t.Skip("set BK_LIVE_ADDONS=1 in the isolated privileged lab")
	}
	if os.Geteuid() != 0 {
		t.Fatal("the live lab must run as root")
	}
	bin := os.Getenv("BK_ADDON_BINARY")
	if bin == "" {
		t.Fatal("BK_ADDON_BINARY required")
	}
	previousBinary := connTestBinary
	defer func() { connTestBinary = previousBinary }()
	connTestBinary = func() (string, error) { return bin, nil }
	prevSoak := connTestSoak
	connTestSoak = 5
	defer func() { connTestSoak = prevSoak }()
	ns := "bk-addon-test-peer"
	run := func(args ...string) {
		t.Helper()
		b, e := exec.Command("ip", args...).CombinedOutput()
		if e != nil {
			t.Fatalf("ip %v: %v %s", args, e, b)
		}
	}
	run("netns", "add", ns)
	defer exec.Command("ip", "netns", "del", ns).Run()
	run("link", "add", "bk-test-host", "type", "veth", "peer", "name", "bk-test-peer")
	defer exec.Command("ip", "link", "del", "bk-test-host").Run()
	run("link", "set", "bk-test-peer", "netns", ns)
	run("addr", "add", "10.201.0.1/30", "dev", "bk-test-host")
	run("link", "set", "bk-test-host", "up")
	run("-n", ns, "addr", "add", "10.201.0.2/30", "dev", "bk-test-peer")
	run("-n", ns, "link", "set", "bk-test-peer", "up")
	run("-n", ns, "link", "set", "lo", "up")
	kinds := strings.Split(os.Getenv("BK_ADDON_KINDS"), ",")
	if len(kinds) == 1 && kinds[0] == "" {
		kinds = []string{"gre", "l2tp-ip", "l2tp-udp", "rgt-direct", "rgt-tcp", "rgt-udp", "paqet", "awg", "alghadir", "ssh"}
	}
	var sshSpec *externaltunnel.Spec
	for _, kind := range kinds {
		if kind == "ssh" {
			var stop func()
			sshSpec, stop = startAdditionalLiveSSH(t, ns)
			defer stop()
			break
		}
	}
	s, e := startExternalTestOptions("10.201.0.1", "10.201.0.2", kinds, sshSpec)
	if e != nil {
		t.Fatal(e)
	}
	defer s.close()
	exe, _ := os.Executable()
	child := exec.Command("ip", "netns", "exec", ns, exe, "-test.run=^TestAdditionalTunnelPeerHelper$", "-test.v")
	child.Env = append(os.Environ(), "BK_ADDITIONAL_PEER="+externalTestLink(s.plan))
	var out bytes.Buffer
	child.Stdout, child.Stderr = &out, &out
	if e = child.Start(); e != nil {
		t.Fatal(e)
	}
	defer child.Process.Kill()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	select {
	case <-s.coord.joined:
	case <-ctx.Done():
		t.Fatal("peer did not join")
	}
	rows := s.run(ctx, nil)
	if e = child.Wait(); e != nil {
		t.Errorf("peer: %v\n%s", e, out.String())
	}
	t.Log("\n" + ConnTestTable(rows))
	t.Log("peer output:\n" + out.String())
	for _, engine := range s.engines {
		if body, e := os.ReadFile(engine.log); e == nil {
			t.Logf("%s log:\n%s", engine.name, body)
		}
	}
	for _, row := range rows {
		if row.Status != ctOK {
			t.Errorf("%s %s: %s", row.Transport, row.Status, row.Detail)
		}
	}
	// Cleanup is checked explicitly, including kernel state and namespaces.
	s.close()
	b, _ := exec.Command("ip", "-j", "link", "show").Output()
	if bytes.Contains(b, []byte("bkx")) {
		t.Errorf("additional interfaces leaked: %s", b)
	}
	b, _ = exec.Command("ip", "netns", "list").Output()
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "bkx") {
			t.Errorf("Alghadir namespace leaked: %s", line)
		}
	}
	// Temporary configs must never be promoted to the permanent store.
	if files, _ := filepath.Glob("/etc/bk/external/xt-*.json"); len(files) > 0 {
		t.Errorf("temporary configs leaked: %v", files)
	}
}
func TestAdditionalTunnelPeerHelper(t *testing.T) {
	link := os.Getenv("BK_ADDITIONAL_PEER")
	if link == "" {
		t.Skip("started by the two-server live test")
	}
	connTestBinary = func() (string, error) { return os.Getenv("BK_ADDON_BINARY"), nil }
	rows, _, e := runExternalKharej(context.Background(), link, io.Discard, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, r := range rows {
		if r.Status != ctOK {
			fmt.Fprintf(os.Stdout, "%s %s %s\n", r.Transport, r.Status, r.Detail)
		}
	}
}

func TestAdditionalRGTOverLoopback(t *testing.T) {
	if os.Getenv("BK_LIVE_RGT") != "1" {
		t.Skip("requires the pinned RGT core")
	}
	bin := os.Getenv("BK_ADDON_BINARY")
	previous := connTestBinary
	previousSoak := connTestSoak
	connTestBinary = func() (string, error) { return bin, nil }
	connTestSoak = 4
	defer func() { connTestBinary, connTestSoak = previous, previousSoak }()
	s, e := startExternalTest("127.0.0.1", "127.0.0.2", []string{"rgt-tcp", "rgt-udp"})
	if e != nil {
		t.Fatal(e)
	}
	defer s.close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	done := make(chan error, 1)
	go func() { _, _, e := runExternalKharej(ctx, externalTestLink(s.plan), io.Discard, nil); done <- e }()
	select {
	case <-s.coord.joined:
	case <-ctx.Done():
		t.Fatal("peer never joined")
	}
	rows := s.run(ctx, nil)
	if e := <-done; e != nil {
		t.Fatal(e)
	}
	t.Log("\n" + ConnTestTable(rows))
	for _, r := range rows {
		if r.Status != ctOK {
			t.Errorf("%s %s %s", r.Transport, r.Status, r.Detail)
		}
	}
}

// A dedicated OpenSSH daemon inside the peer namespace verifies the production
// SSH path, including the real backend socket and its host-key check.
func startAdditionalLiveSSH(t *testing.T, ns string) (*externaltunnel.Spec, func()) {
	t.Helper()
	dir := t.TempDir()
	key, hostKey := filepath.Join(dir, "id"), filepath.Join(dir, "host")
	for _, path := range []string{key, hostKey} {
		if b, e := exec.Command("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", path).CombinedOutput(); e != nil {
			t.Fatalf("SSH fixture: %v %s", e, b)
		}
	}
	pub, _ := os.ReadFile(hostKey + ".pub")
	fields := strings.Fields(string(pub))
	known := filepath.Join(dir, "known_hosts")
	if e := os.WriteFile(known, []byte("10.201.0.2 "+strings.Join(fields[:2], " ")+"\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.MkdirAll("/run/sshd", 0755); e != nil {
		t.Fatal(e)
	}
	cfg := filepath.Join(dir, "sshd_config")
	body := fmt.Sprintf("Port 22\nListenAddress 10.201.0.2\nHostKey %s\nPidFile %s\nAuthorizedKeysFile %s.pub\nPermitRootLogin prohibit-password\nPasswordAuthentication no\nKbdInteractiveAuthentication no\nUsePAM yes\nAllowTcpForwarding yes\n", hostKey, filepath.Join(dir, "pid"), key)
	if e := os.WriteFile(cfg, []byte(body), 0600); e != nil {
		t.Fatal(e)
	}
	cmd := exec.Command("ip", "netns", "exec", ns, "/usr/sbin/sshd", "-D", "-e", "-f", cfg)
	var log bytes.Buffer
	cmd.Stdout, cmd.Stderr = &log, &log
	if e := cmd.Start(); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	stop := func() { cmd.Process.Kill(); <-done; t.Log("OpenSSH fixture: " + log.String()) }
	ready := false
	for i := 0; i < 100; i++ {
		c, e := net.DialTimeout("tcp", "10.201.0.2:22", 100*time.Millisecond)
		if e == nil {
			c.Close()
			ready = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !ready {
		stop()
		t.Fatalf("SSH daemon not ready: %s", log.String())
	}
	spec := externaltunnel.New("ssh-live", "ssh", "iran")
	spec.Port, spec.SSHKey, spec.KnownHosts = 22, key, known
	return &spec, stop
}
