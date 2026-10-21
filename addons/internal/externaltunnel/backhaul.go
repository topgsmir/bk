package externaltunnel

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const backhaulArchiveSHA256 = "57936b9f5bfa3f4d283792068eeae78c64e3bf7159362160af415a55c958e614"
const backhaulCoreSHA256 = "6931b4cefad948af639ad97984356ce04d4e61c4d2efdd9877eea455ef6514b1"

func checkBackhaulCore() error {
	b, err := os.ReadFile(filepath.Join(CoreDir, "backhaul-premium"))
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	if hex.EncodeToString(sum[:]) != backhaulCoreSHA256 {
		return fmt.Errorf("Backhaul supplied core checksum differs")
	}
	return nil
}
func installBackhaul(ctx context.Context, out io.Writer) error {
	b, err := bundledBackhaul()
	if err != nil {
		return err
	}
	for _, p := range []string{"ip", "iptables", "ip6tables", "unshare", "mount"} {
		if _, err = exec.LookPath(p); err != nil {
			if err = runCommand(ctx, command("apt-get", "update")); err != nil {
				return err
			}
			if err = runCommand(ctx, command("apt-get", "install", "-y", "iproute2", "iptables", "util-linux", "mount")); err != nil {
				return err
			}
			break
		}
	}
	if err = installCoreBytes(b, Asset{SHA256: backhaulArchiveSHA256, Member: "backhaul-premium"}, filepath.Join(CoreDir, "backhaul-premium")); err != nil {
		return err
	}
	fmt.Fprintln(out, "Pinned supplied Backhaul core installed; original bytes preserved.")
	return nil
}

// The supplied core refuses private-only machines before starting any transport.
// Report this real requirement instead of modifying or inventing an address.
func checkBackhaulLocalIP(s Spec) error {
	ip := net.ParseIP(s.LocalIP)
	if !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return fmt.Errorf("supplied Backhaul requires a public IPv4 assigned to a local NIC (a NAT public address alone is insufficient)")
	}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, address := range addresses {
		if current, _, err := net.ParseCIDR(address.String()); err == nil && current.Equal(ip) {
			return nil
		}
	}
	return fmt.Errorf("Backhaul local IPv4 %s is not assigned to this server", s.LocalIP)
}
func backhaulTLS(s Spec, dir string) (*tls.Config, error) {
	seed, _ := base64.StdEncoding.DecodeString(s.Key("backhaul/tls"))
	curve := elliptic.P256()
	scalar := new(big.Int).SetBytes(seed)
	scalar.Mod(scalar, new(big.Int).Sub(curve.Params().N, big.NewInt(1)))
	scalar.Add(scalar, big.NewInt(1))
	x, y := curve.ScalarBaseMult(scalar.Bytes())
	key := &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y}, D: scalar}
	sum := sha256.Sum256(seed)
	tmpl := &x509.Certificate{SerialNumber: new(big.Int).SetBytes(sum[:16]), Subject: pkix.Name{CommonName: "bk-backhaul.local"}, DNSNames: []string{"bk-backhaul.local"}, NotBefore: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), NotAfter: time.Date(2120, 1, 1, 0, 0, 0, 0, time.UTC), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}}
	caSeed, _ := base64.StdEncoding.DecodeString(s.Key("backhaul/tls-ca"))
	caKey := ed25519.NewKeyFromSeed(caSeed)
	caSum := sha256.Sum256(caSeed)
	ca := &x509.Certificate{SerialNumber: new(big.Int).SetBytes(caSum[:16]), Subject: pkix.Name{CommonName: "bk-backhaul-ca"}, NotBefore: tmpl.NotBefore, NotAfter: tmpl.NotAfter, IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature}
	caDER, err := x509.CreateCertificate(rand.Reader, ca, ca, caKey.Public(), caKey)
	if err != nil {
		return nil, err
	}
	tmpl.IsCA = false
	tmpl.KeyUsage = x509.KeyUsageDigitalSignature
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	certPEM := append(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), caPEM...)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	for name, body := range map[string][]byte{"tls.pem": certPEM, "tls-key.pem": keyPEM} {
		if err := os.WriteFile(filepath.Join(dir, name), body, 0600); err != nil {
			return nil, err
		}
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(caPEM)
	return &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: roots, ClientCAs: roots, ClientAuth: tls.RequireAndVerifyClientCert, ServerName: "bk-backhaul.local", MinVersion: tls.VersionTLS13}, nil
}
func backhaulConfig(s Spec, dir string) string {
	carrier := strings.TrimPrefix(s.Kind, "b3-")
	tun := strings.HasPrefix(carrier, "tun-")
	ipx := strings.HasPrefix(carrier, "tun-ipx-")
	native := carrier
	if tun {
		native = "tun"
	}
	var b strings.Builder
	if !ipx {
		if s.Side == "iran" {
			fmt.Fprintf(&b, "[listener]\nbind_addr = %q\n", net.JoinHostPort(s.LocalIP, strconv.Itoa(s.Port)))
		} else {
			fmt.Fprintf(&b, "[dialer]\nremote_addr = %q\ndial_timeout = 5\nretry_interval = 1\n", net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)))
		}
	}
	fmt.Fprintf(&b, "[transport]\ntype = %q\nnodelay = true\nconnection_pool = 2\nheartbeat_interval = 2\nheartbeat_timeout = 10\n", native)
	if carrier == "tcp" && s.Protocol == "udp" && s.Side == "iran" {
		fmt.Fprintln(&b, "accept_udp = true\n[accept_udp]\nring_size = 64\nframe_size = 2048\npeer_idle_timeout_s = 120\nwrite_timeout_ms = 3")
	}
	if ipx {
		fmt.Fprintf(&b, "[security]\nenable_encryption = true\nalgorithm = \"aes-256-gcm\"\npsk = %q\nkdf_iterations = 100000\n", s.Key("backhaul/ipx"))
	} else {
		fmt.Fprintf(&b, "[security]\ntoken = %q\n", s.Secret)
	}
	fmt.Fprintln(&b, "[tuning]\nauto_tuning = false\nworkers = 2\nchannel_size = 1024\nbatch_size = 64\nso_sndbuf = 4194304\n[logging]\nlog_level = \"debug\"")
	if tun {
		encapsulation := "tcp"
		if ipx {
			encapsulation = "ipx"
		}
		fmt.Fprintf(&b, "[tun]\nencapsulation = %q\nname = %q\nlocal_addr = %q\nremote_addr = %q\nhealth_port = %d\nmtu = %d\n", encapsulation, s.Interface(), s.TunnelIP(), s.PeerTunnelIP()+"/"+strings.Split(s.TunnelIP(), "/")[1], s.SourcePort+2, s.MTU)
	}
	if ipx {
		role := "server"
		if s.Side == "kharej" {
			role = "client"
		}
		fmt.Fprintf(&b, "[ipx]\nmode = %q\nprofile = %q\nlisten_ip = %q\ndst_ip = %q\ninterface = %q\nicmp_type = 0\nicmp_code = 0\n", role, strings.TrimPrefix(carrier, "tun-ipx-"), s.LocalIP, s.PeerIP, s.WAN)
	}
	if carrier == "wss" || carrier == "wssmux" || carrier == "anytls" {
		fmt.Fprintln(&b, "[tls]\nsni = \"bk-backhaul.local\"")
		if s.Side == "iran" {
			fmt.Fprintf(&b, "tls_cert = %q\ntls_key = %q\n", filepath.Join(dir, "tls.pem"), filepath.Join(dir, "tls-key.pem"))
		}
	}
	if strings.HasSuffix(carrier, "mux") {
		fmt.Fprintln(&b, "[mux]\nmux_version = 2\nmux_framesize = 32768\nmux_recievebuffer = 4194304\nmux_streambuffer = 262144\nmux_concurrency = 2")
	}
	if !tun && s.Side == "iran" {
		fmt.Fprintf(&b, "[ports]\nmapping = [%q]\n", fmt.Sprintf("%d=%d", s.SourcePort, s.SourcePort+3))
	}
	return b.String()
}

// The native IPX TCP helper adds/deletes an unlabelled RST rule. Intercept only
// that exact command in this child's private PATH and label it with its owner.
// No system sudo or iptables executable is replaced.
func backhaulRSTHelper(s Spec, dir string) (string, Command, error) {
	tool, err := exec.LookPath("iptables")
	if err != nil {
		return "", Command{}, err
	}
	helperDir := filepath.Join(dir, "helpers")
	if err = os.Mkdir(helperDir, 0700); err != nil {
		return "", Command{}, err
	}
	rule := []string{"OUTPUT", "-p", "tcp", "-d", s.PeerIP, "--tcp-flags", "RST", "RST", "-j", "DROP"}
	comment := "bk-b3-" + s.Name
	script := fmt.Sprintf("#!/bin/sh\nset -eu\ncase \"$*\" in\n%q|%q) ;;\n*) echo 'Backhaul helper refused unexpected command' >&2; exit 126;;\nesac\nshift\nexec %q -w \"$@\" -m comment --comment %q\n", "iptables -A "+strings.Join(rule, " "), "iptables -D "+strings.Join(rule, " "), tool, comment)
	if err = os.WriteFile(filepath.Join(helperDir, "sudo"), []byte(script), 0700); err != nil {
		return "", Command{}, err
	}
	down := command(tool, append(append([]string{"-w", "-D"}, rule...), "-m", "comment", "--comment", comment)...)
	return helperDir, down, nil
}
func backhaulFirewall(s Spec) []Command {
	if Layer3(s.Kind) {
		return nil
	}
	source, port := "127.0.0.1", s.SourcePort
	if s.Side == "kharej" {
		port += 3
	}
	return []Command{command("iptables", "-w", "-I", "INPUT", "1", "-p", s.Protocol, "--dport", strconv.Itoa(port), "!", "-s", source, "-m", "comment", "--comment", "bk-b3-"+s.Name, "-j", "DROP")}
}
func runBackhaul(ctx context.Context, s Spec, dir string) error {
	if Layer3(s.Kind) {
		if _, err := net.InterfaceByName(s.Interface()); err == nil {
			return fmt.Errorf("Backhaul interface %s already exists", s.Interface())
		}
	}
	if strings.HasPrefix(s.Kind, "b3-tun-ipx-") && s.WAN == "" {
		body, err := exec.CommandContext(ctx, "ip", "-j", "route", "get", s.PeerIP).Output()
		if err != nil {
			return err
		}
		// Read only the physical route, without probing or changing the host routes.
		var routes []struct{ Dev string }
		if err = json.Unmarshal(body, &routes); err != nil || len(routes) == 0 {
			return fmt.Errorf("Backhaul route interface unavailable")
		}
		s.WAN = routes[0].Dev
	}
	config, err := backhaulTLS(s, dir)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "backhaul.toml")
	if err = os.WriteFile(path, []byte(backhaulConfig(s, dir)), 0600); err != nil {
		return err
	}
	var added []Command
	defer func() {
		for i := len(added) - 1; i >= 0; i-- {
			c := added[i]
			c.Args = append([]string(nil), c.Args...)
			c.Args[1] = "-D"
			c.Args = append(c.Args[:3], c.Args[4:]...)
			cleanup(c)
		}
	}()
	for _, c := range backhaulFirewall(s) {
		if err = runCommand(ctx, c); err != nil {
			return err
		}
		added = append(added, c)
	}
	pathEnv := os.Getenv("PATH")
	if s.Kind == "b3-tun-ipx-tcp" {
		helpers, down, e := backhaulRSTHelper(s, dir)
		if e != nil {
			return e
		}
		pathEnv = helpers + ":" + pathEnv
		// Native graceful cleanup handles the common path. Remove any duplicate
		// rules left by a native restart; each rule includes this tunnel's owner.
		defer func() {
			check := down
			check.Args = append([]string(nil), down.Args...)
			check.Args[1] = "-C"
			for i := 0; i < 64; i++ {
				if runCommand(context.Background(), check) != nil {
					break
				}
				cleanup(down)
			}
		}()
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var workers sync.WaitGroup
	done := make(chan error, 3)
	start := func(fn func() error) { workers.Add(1); go func() { defer workers.Done(); done <- fn() }() }
	guard := solarpassMountGuard
	var guardArgs []string
	if s.Kind == "b3-tun-tcp" {
		// The native TUN constructor insists on writing ip_forward=1 even though
		// this adapter terminates traffic locally and performs no host routing.
		// Overlay just that path with an owned file in the private mount namespace.
		flag := filepath.Join(dir, "ip-forward")
		if err = os.WriteFile(flag, []byte("1\n"), 0600); err != nil {
			cancel()
			return err
		}
		guard = strings.Replace(guard, `exec "$@"`, `mount --bind "$1" /proc/sys/net/ipv4/ip_forward
shift
exec "$@"`, 1)
		guardArgs = append(guardArgs, flag)
	}
	start(func() error {
		args := []string{"GOMAXPROCS=2", "PATH=" + pathEnv, "unshare", "--mount", "--propagation", "private", "/bin/sh", "-c", guard, "bk-backhaul"}
		args = append(args, guardArgs...)
		args = append(args, filepath.Join(CoreDir, "backhaul-premium"), "-c", path)
		return childRedacted(ctx, "env", []string{s.Secret, s.Key("backhaul/ipx")}, args...)
	})

	if Layer3(s.Kind) {
		defer func() {
			if _, err := net.InterfaceByName(s.Interface()); err == nil {
				cleanup(command("ip", "link", "del", s.Interface()))
			}
		}()
		if err = waitBackhaulInterface(ctx, s); err != nil {
			cancel()
			workers.Wait()
			return err
		}
	}

	if s.Protocol == "tcp" {
		address := net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort))
		if Layer3(s.Kind) {
			address = net.JoinHostPort(s.PeerTunnelIP(), strconv.Itoa(s.SourcePort))
		}
		if s.Side == "iran" {
			start(func() error {
				return serveTCP(ctx, s.Listen, func(ctx context.Context) (net.Conn, error) {
					d := tls.Dialer{NetDialer: &net.Dialer{Timeout: 5 * time.Second}, Config: config}
					return d.DialContext(ctx, "tcp", address)
				})
			})
		} else {
			listen := net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort+3))
			if Layer3(s.Kind) {
				listen = net.JoinHostPort(strings.Split(s.TunnelIP(), "/")[0], strconv.Itoa(s.SourcePort))
			}
			start(func() error {
				return serveTCP(ctx, listen, func(ctx context.Context) (net.Conn, error) {
					return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", s.Target)
				}, config)
			})
		}
	} else {
		if Layer3(s.Kind) {
			forward := s
			forward.Target = net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort))
			forward.Backend = net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort+3))
			if s.Side == "iran" {
				forward.Listen = forward.Target
			}
			start(func() error { return runForward(ctx, forward) })
		}
		start(func() error { return runPackage3UDPBridge(ctx, s, "backhaul/udp/envelope") })
	}
	select {
	case err = <-done:
	case <-ctx.Done():
	}
	cancel()
	workers.Wait()
	return err
}
func waitBackhaulInterface(ctx context.Context, s Spec) error {
	ready, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	for {
		if iface, err := net.InterfaceByName(s.Interface()); err == nil {
			addresses, _ := iface.Addrs()
			for _, a := range addresses {
				ip, _, _ := net.ParseCIDR(a.String())
				if ip != nil && ip.String() == strings.Split(s.TunnelIP(), "/")[0] {
					return nil
				}
			}
		}
		select {
		case <-ready.Done():
			return fmt.Errorf("Backhaul TUN did not become ready: %w", ready.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}
