//go:build linux

package manage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	if value := os.Getenv("BK_ADDON_SOAK"); value != "" {
		n, e := strconv.Atoi(value)
		if e != nil || n < 1 || n > 300 {
			t.Fatal("BK_ADDON_SOAK must be 1..300")
		}
		connTestSoak = n
	}
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
	if os.Getenv("BK_PACKAGE3_TEST") == "1" {
		run("-6", "addr", "add", "fd42:3::1/64", "dev", "bk-test-host", "nodad")
		run("-n", ns, "-6", "addr", "add", "fd42:3::2/64", "dev", "bk-test-peer", "nodad")
	}
	run("-n", ns, "addr", "add", "10.201.0.2/30", "dev", "bk-test-peer")
	run("-n", ns, "link", "set", "bk-test-peer", "up")
	run("-n", ns, "link", "set", "lo", "up")
	kinds := strings.Split(os.Getenv("BK_ADDON_KINDS"), ",")
	if len(kinds) == 1 && kinds[0] == "" {
		kinds = []string{"gre", "l2tp-ip", "l2tp-udp", "rgt-direct", "rgt-tcp", "rgt-udp", "paqet", "awg", "alghadir", "ssh", "ssh-reverse"}
	}
	host, peer := "10.201.0.1", "10.201.0.2"
	for _, kind := range kinds {
		if externaltunnel.Backhaul(kind) {
			host, peer = "192.0.2.10", "192.0.2.11"
			run("addr", "add", host+"/32", "dev", "bk-test-host")
			run("-n", ns, "addr", "add", peer+"/32", "dev", "bk-test-peer")
			run("route", "add", peer+"/32", "dev", "bk-test-host", "src", host)
			run("-n", ns, "route", "add", host+"/32", "dev", "bk-test-peer", "src", peer)
			break
		}
	}
	for _, kind := range kinds {
		if externaltunnel.Solarpass(kind) || externaltunnel.SingBox(kind) {
			run("route", "add", "default", "via", "10.201.0.2", "dev", "bk-test-host", "src", host)
			run("-n", ns, "route", "add", "default", "via", "10.201.0.1", "dev", "bk-test-peer", "src", peer)
			break
		}
	}
	var sshSpec *externaltunnel.Spec
	for _, kind := range kinds {
		if kind == "ssh" {
			var stop func()
			sshSpec, stop = startAdditionalLiveSSH(t, ns, "10.201.0.2", true)
			defer stop()
			break
		}
	}
	var reverseSpec *externaltunnel.Spec
	for _, kind := range kinds {
		if kind == "ssh-reverse" {
			var stop func()
			reverseSpec, stop = startAdditionalLiveSSH(t, "", "10.201.0.1", true)
			defer stop()
			break
		}
	}
	ipv6 := [2]string{}
	if os.Getenv("BK_PACKAGE3_TEST") == "1" {
		ipv6 = [2]string{"fd42:3::1", "fd42:3::2"}
	}
	s, e := startExternalTestReadyIPv6(host, peer, kinds, sshSpec, reverseSpec, nil, ipv6)
	if e != nil {
		t.Fatal(e)
	}
	defer s.close()
	exe, _ := os.Executable()
	child := exec.Command("ip", "netns", "exec", ns, exe, "-test.run=^TestAdditionalTunnelPeerHelper$", "-test.v")
	child.Env = append(os.Environ(), "BK_ADDITIONAL_PEER="+externalTestLink(s.plan))
	if reverseSpec != nil {
		child.Env = append(child.Env, "BK_TEST_REVERSE_KEY="+reverseSpec.SSHKey, "BK_TEST_REVERSE_HOSTS="+reverseSpec.KnownHosts)
	}
	var out bytes.Buffer
	child.Stdout, child.Stderr = &out, &out
	if e = child.Start(); e != nil {
		t.Fatal(e)
	}
	defer child.Process.Kill()
	duration := 3 * time.Minute
	if s.plan.BatchSize > 0 {
		duration = 60 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
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
	if bytes.Contains(b, []byte("bkx")) || bytes.Contains(b, []byte("bk3")) {
		t.Errorf("additional interfaces leaked: %s", b)
	}
	b, _ = exec.Command("ip", "-n", ns, "-j", "link", "show").Output()
	if bytes.Contains(b, []byte("bkx")) || bytes.Contains(b, []byte("bk3")) {
		t.Errorf("peer interfaces leaked: %s", b)
	}
	b, _ = exec.Command("ip", "netns", "list").Output()
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "bkx") {
			t.Errorf("Alghadir namespace leaked: %s", line)
		}
	}
	for _, prefix := range [][]string{nil, {"netns", "exec", ns}} {
		args := append(append([]string(nil), prefix...), "iptables", "-t", "raw", "-S")
		program := "ip"
		if len(prefix) == 0 {
			program, args = "iptables", args[1:]
		}
		for _, table := range []string{"raw", "mangle", "filter"} {
			for i, a := range args {
				if a == "-t" {
					args[i+1] = table
					break
				}
			}
			body, e := exec.Command(program, args...).CombinedOutput()
			if e != nil {
				t.Errorf("cleanup firewall inspection: %v %s", e, body)
			}
			if bytes.Contains(body, []byte("bk-ext-xt-")) || bytes.Contains(body, []byte("bk-s3-xt-")) || bytes.Contains(body, []byte("bk-b3-xt-")) || bytes.Contains(body, []byte("bk-e3-xt-")) {
				t.Errorf("temporary firewall rule leaked: %s", body)
			}
		}
	}

	for _, prefix := range [][]string{nil, {"netns", "exec", ns}} {
		for _, table := range []string{"raw", "mangle", "filter"} {
			program := "ip6tables"
			args := []string{"-t", table, "-S"}
			if len(prefix) > 0 {
				program = "ip"
				args = append(append(append([]string{}, prefix...), "ip6tables"), args...)
			}
			body, err := exec.Command(program, args...).CombinedOutput()
			if err != nil {
				t.Errorf("IPv6 cleanup inspection failed: %v", err)
			}
			for _, owner := range []string{"bk-ext-xt-", "bk-s3-xt-", "bk-b3-xt-", "bk-e3-xt-"} {
				if bytes.Contains(body, []byte(owner)) {
					t.Errorf("owned IPv6 rule leaked: %s", body)
				}
			}
		}
		for _, resource := range []string{"state", "policy"} {
			program := "ip"
			args := append([]string{}, prefix...)
			if len(prefix) > 0 {
				args = append(args, "ip")
			}
			args = append(args, "xfrm", resource, "list")
			if resource == "state" {
				args = append(args, "nokeys")
			}
			body, err := exec.Command(program, args...).Output()
			if err != nil {
				t.Errorf("IPsec cleanup inspection failed: %v", err)
			}
			if len(bytes.TrimSpace(body)) > 0 {
				t.Error("IPsec kernel state was not fully removed")
			}
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
	rows, _, e := runExternalKharejReady(context.Background(), link, io.Discard, nil, func(ctx context.Context, cases []externaltunnel.Spec, out io.Writer) map[string]string {
		for i, spec := range cases {
			if spec.Kind == "ssh-reverse" {
				cases[i].SSHKey = os.Getenv("BK_TEST_REVERSE_KEY")
				cases[i].KnownHosts = os.Getenv("BK_TEST_REVERSE_HOSTS")
			}
		}
		return nil
	})
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
func startAdditionalLiveSSH(t *testing.T, ns, host string, allowForward bool) (*externaltunnel.Spec, func()) {
	t.Helper()
	dir, e := os.MkdirTemp("/root", "bk-addon-ssh-")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	key, hostKey := filepath.Join(dir, "id"), filepath.Join(dir, "host")
	for _, path := range []string{key, hostKey} {
		if b, e := exec.Command("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", path).CombinedOutput(); e != nil {
			t.Fatalf("SSH fixture: %v %s", e, b)
		}
	}
	pub, _ := os.ReadFile(hostKey + ".pub")
	fields := strings.Fields(string(pub))
	probe, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	port := probe.Addr().(*net.TCPAddr).Port
	probe.Close()
	known := filepath.Join(dir, "known_hosts")
	if e := os.WriteFile(known, []byte(fmt.Sprintf("[%s]:%d ", host, port)+strings.Join(fields[:2], " ")+"\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.MkdirAll("/run/sshd", 0755); e != nil {
		t.Fatal(e)
	}
	cfg := filepath.Join(dir, "sshd_config")
	body := fmt.Sprintf("Port %d\nListenAddress %s\nHostKey %s\nPidFile %s\nAuthorizedKeysFile %s.pub\nPermitRootLogin prohibit-password\nPasswordAuthentication no\nKbdInteractiveAuthentication no\nUsePAM yes\nAllowTcpForwarding yes\n", port, host, hostKey, filepath.Join(dir, "pid"), key)
	if !allowForward {
		body = strings.Replace(body, "AllowTcpForwarding yes", "AllowTcpForwarding no", 1)
	}
	if e := os.WriteFile(cfg, []byte(body), 0600); e != nil {
		t.Fatal(e)
	}
	cmd := exec.Command("/usr/sbin/sshd", "-D", "-e", "-f", cfg)
	if ns != "" {
		cmd = exec.Command("ip", "netns", "exec", ns, "/usr/sbin/sshd", "-D", "-e", "-f", cfg)
	}
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
		c, e := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 100*time.Millisecond)
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
	spec.Port, spec.SSHKey, spec.KnownHosts = port, key, known
	return &spec, stop
}

func TestAdditionalStopKillsEveryUnresponsiveProcessGroup(t *testing.T) {
	var engines []*ctEngine
	for i := 0; i < 2; i++ {
		cmd := exec.Command("sh", "-c", "trap '' TERM; echo ready; sleep 60")
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		pipe, e := cmd.StdoutPipe()
		if e != nil {
			t.Fatal(e)
		}
		if e = cmd.Start(); e != nil {
			t.Fatal(e)
		}
		t.Cleanup(func() { syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) })
		if _, e = bufio.NewReader(pipe).ReadString('\n'); e != nil {
			t.Fatal(e)
		}
		engine := &ctEngine{cmd: cmd, done: make(chan struct{})}
		go func() { cmd.Wait(); close(engine.done) }()
		engines = append(engines, engine)
	}
	done := make(chan struct{})
	go func() { stopAdditionalEnginesWithin(engines, 20*time.Millisecond); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cleanup hung after killing the first process")
	}
	for _, engine := range engines {
		select {
		case <-engine.done:
		default:
			t.Fatal("process remained alive")
		}
	}
}

// OpenSSH can permit a login while denying -L/-R. Verify that this becomes a
// forwarding failure with a reason, never a skipped or successful traffic test.
func TestAdditionalSSHLoginDoesNotImplyForwardingPermission(t *testing.T) {
	if os.Getenv("BK_LIVE_ADDONS") != "1" {
		t.Skip("requires isolated OpenSSH lab")
	}
	s, stop := startAdditionalLiveSSH(t, "", "127.0.0.1", false)
	defer stop()
	s.Kind, s.Side, s.LocalIP, s.PeerIP = "ssh-reverse", "kharej", "127.0.0.1", "127.0.0.1"
	s.SourcePort = ctPickPort(map[int]bool{}, true)
	s.Connections = 1
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if e := externaltunnel.CheckSSHAuthentication(ctx, *s); e != nil {
		t.Fatal("fixture login failed:", e)
	}
	if e := externaltunnel.Run(ctx, *s); e == nil || !strings.Contains(e.Error(), "AllowTcpForwarding") {
		t.Fatalf("remote forwarding denial was not diagnosed: %v", e)
	}
}

// Force the SSH control socket closed while the two managers stay alive. Real
// OpenSSH must release/recreate -R, then a fresh payload must cross the same port.
func TestAdditionalSSHReverseReconnects(t *testing.T) {
	if os.Getenv("BK_LIVE_ADDONS") != "1" {
		t.Skip("requires isolated OpenSSH lab")
	}
	sshSpec, stop := startAdditionalLiveSSH(t, "", "127.0.0.1", true)
	defer stop()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	proxy, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer proxy.Close()
	var mu sync.Mutex
	var active []net.Conn
	var copies sync.WaitGroup
	closeConnections := func() {
		mu.Lock()
		defer mu.Unlock()
		for _, c := range active {
			c.Close()
		}
		active = nil
	}
	acceptDone := make(chan struct{})
	defer func() { cancel(); proxy.Close(); <-acceptDone; closeConnections(); copies.Wait() }()
	go func() {
		defer close(acceptDone)
		for {
			a, e := proxy.Accept()
			if e != nil {
				return
			}
			b, e := (&net.Dialer{Timeout: 3 * time.Second}).DialContext(ctx, "tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(sshSpec.Port)))
			if e != nil {
				a.Close()
				return
			}
			mu.Lock()
			active = append(active, a, b)
			mu.Unlock()
			copies.Add(1)
			go func(a, b net.Conn) {
				defer copies.Done()
				done := make(chan struct{}, 2)
				go func() { io.Copy(a, b); done <- struct{}{} }()
				go func() { io.Copy(b, a); done <- struct{}{} }()
				<-done
				a.Close()
				b.Close()
				<-done
			}(a, b)
		}
	}()
	known, e := os.ReadFile(sshSpec.KnownHosts)
	if e != nil {
		t.Fatal(e)
	}
	proxyPort := proxy.Addr().(*net.TCPAddr).Port
	known = []byte(strings.Replace(string(known), ":"+strconv.Itoa(sshSpec.Port)+" ", ":"+strconv.Itoa(proxyPort)+" ", 1))
	if e = os.WriteFile(sshSpec.KnownHosts, known, 0600); e != nil {
		t.Fatal(e)
	}
	backend, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer backend.Close()
	go ctServeTCPEcho(backend)
	used := map[int]bool{}
	s := externaltunnel.New("reconnect", "ssh-reverse", "iran")
	s.LocalIP, s.PeerIP = "127.0.0.1", "127.0.0.1"
	s.Port = proxyPort
	s.SourcePort = ctPickPort(used, true)
	s.Listen = net.JoinHostPort("127.0.0.1", strconv.Itoa(ctPickPort(used, true)))
	s.Target = backend.Addr().String()
	peer := s.Mirror()
	peer.SSHKey, peer.KnownHosts = sshSpec.SSHKey, sshSpec.KnownHosts
	iranDone, peerDone := make(chan error, 1), make(chan error, 1)
	go func() { iranDone <- externaltunnel.Run(ctx, s) }()
	go func() { peerDone <- externaltunnel.Run(ctx, peer) }()
	defer func() {
		cancel()
		select {
		case e := <-iranDone:
			if e != nil {
				t.Error(e)
			}
		case <-time.After(3 * time.Second):
			t.Error("Iran manager did not stop")
		}
		select {
		case e := <-peerDone:
			if e != nil {
				t.Error(e)
			}
		case <-time.After(3 * time.Second):
			t.Error("reverse manager did not stop")
		}
	}()
	exchange := func(payload []byte) {
		t.Helper()
		until := time.Now().Add(8 * time.Second)
		var last error
		for time.Now().Before(until) {
			c, e := net.DialTimeout("tcp", s.Listen, 200*time.Millisecond)
			if e == nil {
				c.SetDeadline(time.Now().Add(time.Second))
				done := make(chan error, 1)
				go func() { _, e := c.Write(payload); done <- e }()
				got := make([]byte, len(payload))
				_, e = io.ReadFull(c, got)
				c.Close()
				<-done
				if e == nil && bytes.Equal(got, payload) {
					return
				}
			}
			last = e
			time.Sleep(100 * time.Millisecond)
		}
		t.Fatalf("reverse data did not recover: %v", last)
	}
	exchange(bytes.Repeat([]byte("before-disconnect"), 4096))
	closeConnections()
	exchange(bytes.Repeat([]byte("after-disconnect"), 4096))
}
