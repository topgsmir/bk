//go:build linux

package externaltunnel

import (
	"bufio"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
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
	"time"

	"github.com/topgsmir/bk/internal/tunnel/l3"
	"github.com/xtaci/kcp-go/v5"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/sys/unix"
)

func udp2rawAsset(arch string) (Asset, error) {
	member := map[string]string{"amd64": "udp2raw_amd64", "arm64": "udp2raw_arm", "arm": "udp2raw_arm", "386": "udp2raw_x86"}[arch]
	// The upstream ARM binary is 32-bit; do not assume an ARM64 host has its loader.
	if arch == "arm64" {
		return Asset{}, fmt.Errorf("udp2raw has no pinned native ARM64 binary; Alghadir currently supports amd64/386/arm")
	}
	if member == "" {
		return Asset{}, fmt.Errorf("no pinned udp2raw binary for %s", arch)
	}
	return Asset{"https://github.com/wangyu-/udp2raw/releases/download/20230206.0/udp2raw_binaries.tar.gz", "503cf5781aa97e50b4954c6bc4622c3ea6be02f6a35def4bb3b3eaf95bd2c7e8", member}, nil
}
func obfsState(s Spec) ([]byte, string) {
	seed := sha256.Sum256([]byte("bk-obfs4-node/" + s.Secret))
	priv := sha256.Sum256([]byte("bk-obfs4-key/" + s.Secret))
	pub, e := curve25519.X25519(priv[:], curve25519.Basepoint)
	if e != nil {
		panic(e)
	}
	drbg := sha256.Sum256([]byte("bk-obfs4-drbg/" + s.Secret))
	state, _ := json.Marshal(map[string]any{"node-id": hex.EncodeToString(seed[:20]), "private-key": hex.EncodeToString(priv[:]), "public-key": hex.EncodeToString(pub), "drbg-seed": hex.EncodeToString(drbg[:24]), "iat-mode": 0})
	cert := append(append([]byte{}, seed[:20]...), pub...)
	return state, strings.TrimSuffix(base64.StdEncoding.EncodeToString(cert), "==")
}

// sessionKey authenticates both ends and creates fresh IPsec/udp2raw/KCP keys
// for each outer connection. A restarted process never reuses an ESP AEAD nonce.
func sessionKey(c net.Conn, s Spec) (string, error) {
	_ = c.SetDeadline(time.Now().Add(15 * time.Second))
	defer c.SetDeadline(time.Time{})
	own := make([]byte, 32)
	if _, e := rand.Read(own); e != nil {
		return "", e
	}
	other := make([]byte, 32)
	if s.Side == "iran" {
		if _, e := c.Write(own); e != nil {
			return "", e
		}
		if _, e := io.ReadFull(c, other); e != nil {
			return "", e
		}
	} else {
		if _, e := io.ReadFull(c, other); e != nil {
			return "", e
		}
		if _, e := c.Write(own); e != nil {
			return "", e
		}
	}
	a, b := own, other
	if s.Side != "iran" {
		a, b = other, own
	}
	m := hmac.New(sha256.New, []byte(s.Secret))
	m.Write([]byte("bk-alghadir-session-v1"))
	m.Write(a)
	m.Write(b)
	key := m.Sum(nil)
	proof := func(side string) []byte { m := hmac.New(sha256.New, key); m.Write([]byte(side)); return m.Sum(nil) }
	received := make([]byte, 32)
	if s.Side == "iran" {
		if _, e := c.Write(proof(s.Side)); e != nil {
			return "", e
		}
		if _, e := io.ReadFull(c, received); e != nil {
			return "", e
		}
	} else {
		if _, e := io.ReadFull(c, received); e != nil {
			return "", e
		}
		if _, e := c.Write(proof(s.Side)); e != nil {
			return "", e
		}
	}
	peer := "iran"
	if s.Side == "iran" {
		peer = "kharej"
	}
	if subtle.ConstantTimeCompare(received, proof(peer)) != 1 {
		return "", fmt.Errorf("Alghadir paired secret does not match")
	}
	return hex.EncodeToString(key), nil
}

// obfsConnection supervises a standard Tor managed obfs4 transport. No Tor
// circuit is involved: the two VPS talk directly through the obfs4 stream.
func obfsConnection(ctx context.Context, s Spec, dir string) (net.Conn, func(), error) {
	state, cert := obfsState(s)
	stateDir := filepath.Join(dir, "obfs")
	if e := os.Mkdir(stateDir, 0700); e != nil {
		return nil, nil, e
	}
	if e := os.WriteFile(filepath.Join(stateDir, "obfs4_state.json"), state, 0600); e != nil {
		return nil, nil, e
	}
	cmd := exec.CommandContext(ctx, "obfs4proxy")
	cmd.Env = append(os.Environ(), "TOR_PT_MANAGED_TRANSPORT_VER=1", "TOR_PT_STATE_LOCATION="+stateDir)
	var listener net.Listener
	var e error
	if s.Side == "kharej" {
		listener, e = net.Listen("tcp", "127.0.0.1:0")
		if e != nil {
			return nil, nil, e
		}
		cmd.Env = append(cmd.Env, "TOR_PT_SERVER_TRANSPORTS=obfs4", "TOR_PT_SERVER_BINDADDR=obfs4-"+net.JoinHostPort(s.LocalIP, strconv.Itoa(s.Port)), "TOR_PT_ORPORT="+listener.Addr().String())
	} else {
		cmd.Env = append(cmd.Env, "TOR_PT_CLIENT_TRANSPORTS=obfs4")
	}
	pipe, e := cmd.StdoutPipe()
	if e != nil {
		if listener != nil {
			listener.Close()
		}
		return nil, nil, e
	}
	cmd.Stderr = os.Stderr
	if e = cmd.Start(); e != nil {
		if listener != nil {
			listener.Close()
		}
		return nil, nil, e
	}
	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()
	stop := func() {
		if listener != nil {
			listener.Close()
		}
		_ = cmd.Process.Kill()
		<-wait
	}
	ready := make(chan string, 1)
	go func() {
		scan := bufio.NewScanner(pipe)
		for scan.Scan() {
			line := scan.Text()
			if strings.HasPrefix(line, "CMETHOD obfs4 socks5 ") {
				select {
				case ready <- strings.TrimPrefix(line, "CMETHOD obfs4 socks5 "):
				default:
				}
			}
			if strings.HasPrefix(line, "SMETHOD obfs4 ") {
				select {
				case ready <- "server":
				default:
				}
			}
		}
	}()
	var address string
	select {
	case address = <-ready:
	case <-ctx.Done():
		stop()
		return nil, nil, ctx.Err()
	case <-time.After(15 * time.Second):
		stop()
		return nil, nil, fmt.Errorf("obfs4 did not become ready")
	}
	if s.Side == "kharej" {
		go func() { <-ctx.Done(); listener.Close() }()
		c, e := listener.Accept()
		if e != nil {
			stop()
			return nil, nil, e
		}
		return c, stop, nil
	}
	var c net.Conn
	for ctx.Err() == nil {
		c, e = obfsSOCKS(ctx, address, net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)), cert)
		if e == nil {
			return c, stop, nil
		}
		fmt.Fprintln(os.Stderr, "Alghadir obfs4 connection:", e)
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
	}
	stop()
	return nil, nil, ctx.Err()
}
func obfsSOCKS(ctx context.Context, proxy, target, cert string) (net.Conn, error) {
	d := net.Dialer{Timeout: 5 * time.Second}
	c, e := d.DialContext(ctx, "tcp", proxy)
	if e != nil {
		return nil, e
	}
	ok := false
	defer func() {
		if !ok {
			c.Close()
		}
	}()
	_ = c.SetDeadline(time.Now().Add(10 * time.Second))
	if _, e = c.Write([]byte{5, 1, 2}); e != nil {
		return nil, e
	}
	r := make([]byte, 2)
	if _, e = io.ReadFull(c, r); e != nil {
		return nil, e
	}
	if r[0] != 5 || r[1] != 2 {
		return nil, fmt.Errorf("obfs4 SOCKS authentication rejected")
	}
	user := "cert=" + cert + ";iat-mode=0"
	auth := append([]byte{1, byte(len(user))}, []byte(user)...)
	auth = append(auth, 1, 0)
	if _, e = c.Write(auth); e != nil {
		return nil, e
	}
	if _, e = io.ReadFull(c, r); e != nil || r[1] != 0 {
		return nil, fmt.Errorf("obfs4 SOCKS arguments rejected: %v", e)
	}
	host, port, _ := net.SplitHostPort(target)
	n, _ := strconv.Atoi(port)
	ip := net.ParseIP(host).To4()
	request := append([]byte{5, 1, 0, 1}, ip...)
	request = append(request, byte(n>>8), byte(n))
	if _, e = c.Write(request); e != nil {
		return nil, e
	}
	head := make([]byte, 4)
	if _, e = io.ReadFull(c, head); e != nil {
		return nil, e
	}
	if head[1] != 0 {
		return nil, fmt.Errorf("obfs4 connection failed (%d)", head[1])
	}
	size := 0
	switch head[3] {
	case 1:
		size = 6
	case 4:
		size = 18
	case 3:
		var length [1]byte
		if _, e = io.ReadFull(c, length[:]); e != nil {
			return nil, e
		}
		size = int(length[0]) + 2
	default:
		return nil, fmt.Errorf("invalid SOCKS address")
	}
	if _, e = io.CopyN(io.Discard, c, int64(size)); e != nil {
		return nil, e
	}
	_ = c.SetDeadline(time.Time{})
	ok = true
	return c, nil
}
func runAlghadir(ctx context.Context, s Spec, dir string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	outer, stop, e := obfsConnection(ctx, s, dir)
	if e != nil {
		return e
	}
	defer stop()
	defer outer.Close()
	fmt.Fprintln(os.Stdout, "Alghadir: obfs4 connected")
	session, e := sessionKey(outer, s)
	if e != nil {
		return e
	}
	ns := s.Interface() + "n"
	hostDev := s.Interface() + "h"
	nsDev := s.Interface() + "v"
	// A /30 unique to this local side. Neither namespace has a default route.
	h := sha256.Sum256([]byte(s.Name + s.Side))
	block := fmt.Sprintf("198.18.%d.%d", h[0], (int(h[1])/4)*4)
	parts := strings.Split(block, ".")
	base, _ := strconv.Atoi(parts[3])
	prefix := strings.Join(parts[:3], ".") + "."
	hostIP, nsIP := prefix+strconv.Itoa(base+1), prefix+strconv.Itoa(base+2)
	if e = runCommand(ctx, command("ip", "netns", "add", ns)); e != nil {
		return e
	}
	defer cleanup(command("ip", "netns", "del", ns))
	steps := []Command{command("ip", "link", "add", hostDev, "type", "veth", "peer", "name", nsDev), command("ip", "link", "set", nsDev, "netns", ns), command("ip", "addr", "add", hostIP+"/30", "dev", hostDev), command("ip", "link", "set", hostDev, "up"), command("ip", "-n", ns, "addr", "add", nsIP+"/30", "dev", nsDev), command("ip", "-n", ns, "link", "set", nsDev, "up"), command("ip", "-n", ns, "link", "set", "lo", "up")}
	veth := false
	defer func() {
		if veth {
			cleanup(command("ip", "link", "del", hostDev))
		}
	}()
	for i, c := range steps {
		if e = runCommand(ctx, c); e != nil {
			return e
		}
		if i == 0 {
			veth = true
		}
	}
	ipc := filepath.Join(dir, "raw.sock")
	l, e := net.Listen("unix", ipc)
	if e != nil {
		return e
	}
	defer l.Close()
	go func() { <-ctx.Done(); l.Close() }()
	// Parent listeners keep local/public sockets out of the isolated namespace.
	local := s
	local.Secret = session
	frontPort := 22001
	backendPort := 22002
	var forward *l3.Forwarder
	if s.Side == "iran" {
		local.Listen = net.JoinHostPort(nsIP, strconv.Itoa(frontPort))
		forward, e = l3.NewForwarder(l3.Config{Ports: []string{s.Listen + "=" + local.Listen}, AcceptUDP: s.Protocol == "udp"}, nil)
	} else {
		local.Backend = net.JoinHostPort(hostIP, strconv.Itoa(backendPort))
		forward, e = l3.NewForwarder(l3.Config{Ports: []string{local.Backend + "=" + s.Target}, AcceptUDP: s.Protocol == "udp"}, nil)
	}
	if e != nil {
		return e
	}
	doneForward := make(chan error, 1)
	go func() { doneForward <- forward.Run(ctx); cancel() }()
	defer func() { cancel(); <-doneForward }()
	cfg := filepath.Join(dir, "child.json")
	if e = Save(cfg, local); e != nil {
		return e
	}
	bin, e := os.Executable()
	if e != nil {
		return e
	}
	cmd := exec.Command("ip", "netns", "exec", ns, bin, "external", "alghadir-child", "-c", cfg, "--bridge", ipc)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if e = cmd.Start(); e != nil {
		return e
	}
	childDone := make(chan error, 1)
	go func() { e := cmd.Wait(); childDone <- e; l.Close() }()
	defer func() {
		_ = cmd.Process.Signal(unix.SIGTERM)
		select {
		case <-childDone:
		case <-time.After(6 * time.Second):
			_ = cmd.Process.Kill()
			<-childDone
		}
	}()
	inside, e := l.Accept()
	if e != nil {
		return e
	}
	defer inside.Close()
	relay(ctx, outer, inside)
	if ctx.Err() == nil {
		return fmt.Errorf("Alghadir outer stream disconnected")
	}
	return nil
}
func openPacketTUN(name, ip string, mtu int) (*os.File, error) {
	fd, e := unix.Open("/dev/net/tun", unix.O_RDWR|unix.O_CLOEXEC, 0)
	if e != nil {
		return nil, e
	}
	req, e := unix.NewIfreq(name)
	if e != nil {
		unix.Close(fd)
		return nil, e
	}
	req.SetUint16(unix.IFF_TUN | unix.IFF_NO_PI)
	if e = unix.IoctlIfreq(fd, unix.TUNSETIFF, req); e != nil {
		unix.Close(fd)
		return nil, e
	}
	if e = unix.SetNonblock(fd, true); e != nil {
		unix.Close(fd)
		return nil, e
	}
	f := os.NewFile(uintptr(fd), name)
	for _, c := range []Command{command("ip", "addr", "add", ip, "dev", name), command("ip", "link", "set", name, "mtu", strconv.Itoa(mtu), "up")} {
		if e := runCommand(context.Background(), c); e != nil {
			f.Close()
			return nil, e
		}
	}
	return f, nil
}
func packetStream(ctx context.Context, tun *os.File, stream net.Conn) error {
	done := make(chan error, 2)
	go func() {
		b := make([]byte, 65535)
		head := make([]byte, 2)
		for {
			n, e := tun.Read(b)
			if e != nil {
				done <- e
				return
			}
			if n < 20 {
				continue
			}
			binary.BigEndian.PutUint16(head, uint16(n))
			if _, e = stream.Write(append(head, b[:n]...)); e != nil {
				done <- e
				return
			}
		}
	}()
	go func() {
		head := make([]byte, 2)
		b := make([]byte, 65535)
		for {
			if _, e := io.ReadFull(stream, head); e != nil {
				done <- e
				return
			}
			n := int(binary.BigEndian.Uint16(head))
			if n < 20 {
				done <- fmt.Errorf("invalid IP frame")
				return
			}
			if _, e := io.ReadFull(stream, b[:n]); e != nil {
				done <- e
				return
			}
			if _, e := tun.Write(b[:n]); e != nil {
				done <- e
				return
			}
		}
	}()
	var err error
	finished := 0
	select {
	case <-ctx.Done():
	case err = <-done:
		finished = 1
	}
	tun.Close()
	stream.Close()
	for finished < 2 {
		<-done
		finished++
	}
	return err
}
func ipsecGRECommands(s Spec) []Command {
	local, peer := "10.254.0.1", "10.254.0.2"
	if s.Side == "kharej" {
		local, peer = peer, local
	}
	key := func(label string) string {
		a := sha256.Sum256([]byte("bk-ipsec/" + label + s.Secret))
		salt := sha256.Sum256([]byte("bk-ipsec-salt/" + label + s.Secret))
		return "0x" + hex.EncodeToString(a[:]) + hex.EncodeToString(salt[:4])
	}
	outSide, inSide := "iran", "kharej"
	if s.Side == "kharej" {
		outSide, inSide = inSide, outSide
	}
	spi := func(side string) string {
		if side == "iran" {
			return "0x10001"
		}
		return "0x10002"
	}
	commands := []Command{}
	for _, d := range []struct{ src, dst, side string }{{local, peer, outSide}, {peer, local, inSide}} {
		commands = append(commands, command("ip", "xfrm", "state", "add", "src", d.src, "dst", d.dst, "proto", "esp", "spi", spi(d.side), "reqid", "1", "mode", "transport", "replay-window", "64", "aead", "rfc4106(gcm(aes))", key(d.side), "128", "limit", "packet-hard", "4000000000"))
	}
	for _, d := range []struct{ src, dst, dir string }{{local, peer, "out"}, {peer, local, "in"}} {
		commands = append(commands, command("ip", "xfrm", "policy", "add", "dir", d.dir, "src", d.src+"/32", "dst", d.dst+"/32", "proto", "gre", "tmpl", "src", d.src, "dst", d.dst, "proto", "esp", "reqid", "1", "mode", "transport"))
	}
	commands = append(commands, command("ip", "tunnel", "add", s.Interface(), "mode", "gre", "local", local, "remote", peer, "key", strconv.Itoa(s.ID)), command("ip", "addr", "add", s.TunnelIP(), "dev", s.Interface()), command("ip", "link", "set", s.Interface(), "mtu", strconv.Itoa(s.MTU), "up"))
	return commands
}

// AlghadirChild runs only inside its private network namespace. Both packet
// boundaries are real TUN devices: encrypted GRE enters KCP; fake-TCP packets
// from udp2raw enter the outer obfs4 stream via a local Unix socket.
func AlghadirChild(ctx context.Context, s Spec, ipc string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	local, peer := "10.254.0.1", "10.254.0.2"
	rawLocal, rawPeer := "10.254.1.1", "10.254.1.2"
	if s.Side == "kharej" {
		local, peer = peer, local
		rawLocal, rawPeer = rawPeer, rawLocal
	}
	enc, e := openPacketTUN("bkenc", local+"/30", s.MTU+100)
	if e != nil {
		return e
	}
	defer enc.Close()
	raw, e := openPacketTUN("bkraw", rawLocal+"/30", 1500)
	if e != nil {
		return e
	}
	defer raw.Close()
	_ = peer // direct connected route through bkenc
	for _, c := range ipsecGRECommands(s) {
		if e = runCommand(ctx, c); e != nil {
			return e
		}
	}
	bridge, e := net.DialTimeout("unix", ipc, 5*time.Second)
	if e != nil {
		return e
	}
	defer bridge.Close()
	var tasks sync.WaitGroup
	failures := make(chan error, 4)
	start := func(name string, run func() error) {
		tasks.Add(1)
		go func() {
			defer tasks.Done()
			e := run()
			if e == nil && ctx.Err() == nil {
				e = fmt.Errorf("stopped unexpectedly")
			}
			if e != nil {
				failures <- fmt.Errorf("Alghadir %s: %w", name, e)
			}
			cancel()
		}()
	}
	defer func() { cancel(); raw.Close(); bridge.Close(); tasks.Wait() }()
	start("raw packet bridge", func() error { return packetStream(ctx, raw, bridge) })
	args := []string{"-c", "-l", "127.0.0.1:51500", "-r", rawPeer + ":52000", "-k", s.Secret, "--raw-mode", "faketcp", "-a", "--disable-color"}
	if s.Side == "kharej" {
		args = []string{"-s", "-l", rawLocal + ":52000", "-r", "127.0.0.1:51501", "-k", s.Secret, "--raw-mode", "faketcp", "-a", "--disable-color"}
	}
	start("udp2raw", func() error {
		return childRedacted(ctx, filepath.Join(CoreDir, "udp2raw"), []string{s.Secret}, args...)
	})
	key := sha256.Sum256([]byte("bk-alghadir-kcp/" + s.Secret))
	block, e := kcp.NewAESBlockCrypt(key[:])
	if e != nil {
		return e
	}
	var session *kcp.UDPSession
	if s.Side == "kharej" {
		l, e := kcp.ListenWithOptions("127.0.0.1:51501", block, 10, 3)
		if e != nil {
			return e
		}
		defer l.Close()
		go func() { <-ctx.Done(); l.Close() }()
		session, e = l.AcceptKCP()
		if e != nil {
			return e
		}
	} else {
		session, e = kcp.DialWithOptions("127.0.0.1:51500", block, 10, 3)
		if e != nil {
			return e
		}
	}
	defer session.Close()
	session.SetStreamMode(true)
	session.SetNoDelay(1, 10, 2, 1)
	session.SetWindowSize(256, 256)
	session.SetMtu(1200)
	start("forwarder", func() error { return runForward(ctx, s) })
	fmt.Fprintln(os.Stdout, "Alghadir: GRE -> IPsec ESP -> KCP -> udp2raw fake TCP -> obfs4 ready")
	start("encrypted GRE over KCP", func() error { return packetStream(ctx, enc, session) })
	select {
	case e := <-failures:
		return e
	case <-ctx.Done():
		select {
		case e := <-failures:
			return e
		default:
			return nil
		}
	}
}
