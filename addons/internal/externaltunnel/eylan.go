package externaltunnel

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	wgconn "golang.zx2c4.com/wireguard/conn"
	wgdevice "golang.zx2c4.com/wireguard/device"
	wgtun "golang.zx2c4.com/wireguard/tun"
)

// Eylan's supplied panel controls ordinary VPN protocols. These adapters use
// their stock engines with isolated named state; no panel license code is run.
func installEylan(ctx context.Context, kind string, out io.Writer) error {
	if L2TPIPsec(kind) {
		return installL2TPIPsec(ctx, out)
	}
	if AnyConnect(kind) {
		return installAnyConnect(ctx, out)
	}
	if SingBox(kind) {
		return installSingBox(ctx, out)
	}
	packages := []string{"iproute2"}
	if strings.HasPrefix(kind, "e3-openvpn-") {
		packages = append(packages, "openvpn")
	}
	if err := runCommand(ctx, command("apt-get", "update")); err != nil {
		return err
	}
	if err := runCommand(ctx, command("apt-get", append([]string{"install", "-y"}, packages...)...)); err != nil {
		return err
	}
	fmt.Fprintln(out, "Stock VPN dependencies installed for", kind)
	return nil
}
func eylanForward(s Spec) Spec {
	forward := s
	forward.Target = net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort))
	forward.Backend = s.Target
	return forward
}
func runEylan(ctx context.Context, s Spec, dir string) error {
	if _, err := net.InterfaceByName(s.Interface()); err == nil {
		return fmt.Errorf("VPN interface %s already exists", s.Interface())
	}
	if SingBox(s.Kind) {
		return runSingBox(ctx, s, dir)
	}
	if L2TPIPsec(s.Kind) {
		return runL2TPIPsec(ctx, s, dir)
	}
	if AnyConnect(s.Kind) {
		return runAnyConnect(ctx, s, dir)
	}
	switch s.Kind {
	case "e3-wireguard":
		return runEylanWireGuard(ctx, s)
	case "e3-openvpn-tcp", "e3-openvpn-udp":
		return runEylanOpenVPN(ctx, s, dir)
	}
	return fmt.Errorf("unsupported Eylan protocol %s", s.Kind)
}
func runEylanWireGuard(ctx context.Context, s Spec) error {
	tunnel, err := wgtun.CreateTUN(s.Interface(), s.MTU)
	if err != nil {
		return err
	}
	engine := wgdevice.NewDevice(tunnel, wgconn.NewDefaultBind(), wgdevice.NewLogger(wgdevice.LogLevelError, "bk-wireguard: "))
	defer engine.Close()
	other := "kharej"
	if s.Side == "kharej" {
		other = "iran"
	}
	hexKey := func(key string) string { b, _ := base64.StdEncoding.DecodeString(key); return hex.EncodeToString(b) }
	config := fmt.Sprintf("private_key=%s\nlisten_port=%d\nreplace_peers=true\npublic_key=%s\npreshared_key=%s\nendpoint=%s\nallowed_ip=%s/32\npersistent_keepalive_interval=15\n", hexKey(s.Key("eylan/wireguard/"+s.Side)), s.Port, hexKey(s.AWGPublic("eylan/wireguard/"+other)), hexKey(s.Key("eylan/wireguard/psk")), net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)), s.PeerTunnelIP())
	if err = engine.IpcSet(config); err != nil {
		return err
	}
	if err = engine.Up(); err != nil {
		return err
	}
	for _, c := range []Command{command("ip", "addr", "add", s.TunnelIP(), "dev", s.Interface()), command("ip", "link", "set", "dev", s.Interface(), "mtu", strconv.Itoa(s.MTU), "up")} {
		if err = runCommand(ctx, c); err != nil {
			return err
		}
	}
	forwardCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		select {
		case <-engine.Wait():
			cancel()
		case <-forwardCtx.Done():
		}
	}()
	return runForward(forwardCtx, eylanForward(s))
}
func eylanOpenVPNIdentity(s Spec, dir string) error {
	seed, _ := base64.StdEncoding.DecodeString(s.Key("eylan/openvpn/ca"))
	caKey := ed25519.NewKeyFromSeed(seed)
	sum := sha256.Sum256(seed)
	start, end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2120, 1, 1, 0, 0, 0, 0, time.UTC)
	ca := &x509.Certificate{SerialNumber: new(big.Int).SetBytes(sum[:16]), Subject: pkix.Name{CommonName: "bk-openvpn-ca"}, NotBefore: start, NotAfter: end, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature, IsCA: true, BasicConstraintsValid: true}
	caDER, err := x509.CreateCertificate(rand.Reader, ca, ca, caKey.Public(), caKey)
	if err != nil {
		return err
	}
	seed, _ = base64.StdEncoding.DecodeString(s.Key("eylan/openvpn/" + s.Side))
	key := ed25519.NewKeyFromSeed(seed)
	sum = sha256.Sum256(seed)
	eku := x509.ExtKeyUsageClientAuth
	if s.Side == "iran" {
		eku = x509.ExtKeyUsageServerAuth
	}
	leaf := &x509.Certificate{SerialNumber: new(big.Int).SetBytes(sum[:16]), Subject: pkix.Name{CommonName: "bk-openvpn-" + s.Side}, NotBefore: start, NotAfter: end, KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{eku}, BasicConstraintsValid: true}
	if AnyConnect(s.Kind) {
		leaf.IPAddresses = []net.IP{net.ParseIP(s.LocalIP)}
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leaf, ca, key.Public(), caKey)
	if err != nil {
		return err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	var crypt strings.Builder
	crypt.WriteString("-----BEGIN OpenVPN Static key V1-----\n")
	for i := 0; i < 8; i++ {
		b, _ := base64.StdEncoding.DecodeString(s.Key(fmt.Sprintf("eylan/openvpn/tls-crypt/%d", i)))
		crypt.WriteString(hex.EncodeToString(b[:16]) + "\n" + hex.EncodeToString(b[16:]) + "\n")
	}
	crypt.WriteString("-----END OpenVPN Static key V1-----\n")
	for file, body := range map[string][]byte{"ovpn-ca.pem": pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}), "ovpn-cert.pem": pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leafDER}), "ovpn-key.pem": pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}), "ovpn-crypt.key": []byte(crypt.String())} {
		if err := os.WriteFile(filepath.Join(dir, file), body, 0600); err != nil {
			return err
		}
	}
	return nil
}
func eylanOpenVPNConfig(s Spec, dir string) string {
	carrier := "udp"
	if s.Kind == "e3-openvpn-tcp" {
		carrier = "tcp-server"
		if s.Side == "kharej" {
			carrier = "tcp-client"
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "dev %s\ndev-type tun\nproto %s\nlocal %s\nifconfig %s %s\ntopology p2p\ntun-mtu %d\n", s.Interface(), carrier, s.LocalIP, strings.Split(s.TunnelIP(), "/")[0], s.PeerTunnelIP(), s.MTU)
	if s.Side == "iran" {
		fmt.Fprintf(&b, "port %d\ntls-server\ndh none\nremote-cert-tls client\nverify-x509-name bk-openvpn-kharej name\n", s.Port)
	} else {
		fmt.Fprintf(&b, "remote %s %d\nlport %d\ntls-client\nremote-cert-tls server\nverify-x509-name bk-openvpn-iran name\n", s.PeerIP, s.Port, s.SourcePort+1)
	}
	for option, file := range map[string]string{"ca": "ovpn-ca.pem", "cert": "ovpn-cert.pem", "key": "ovpn-key.pem", "tls-crypt": "ovpn-crypt.key"} {
		fmt.Fprintf(&b, "%s %q\n", option, filepath.Join(dir, file))
	}
	fmt.Fprintln(&b, "data-ciphers AES-256-GCM\ndata-ciphers-fallback AES-256-GCM\ntls-version-min 1.3\nallow-compression no\nscript-security 1\nping 2\nping-restart 15\nconnect-retry 1 2\nresolv-retry infinite\nverb 3")
	return b.String()
}
func runEylanOpenVPN(ctx context.Context, s Spec, dir string) error {
	if err := eylanOpenVPNIdentity(s, dir); err != nil {
		return err
	}
	path := filepath.Join(dir, "openvpn.conf")
	if err := os.WriteFile(path, []byte(eylanOpenVPNConfig(s, dir)), 0600); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan error, 2)
	var workers sync.WaitGroup
	workers.Add(1)
	go func() {
		defer workers.Done()
		done <- childRedacted(ctx, "openvpn", []string{s.Secret}, "--config", path)
	}()
	// Register ownership cleanup even if the native constructor stops halfway.
	defer func() {
		if _, err := net.InterfaceByName(s.Interface()); err == nil {
			cleanup(command("ip", "link", "del", s.Interface()))
		}
	}()
	if err := waitBackhaulInterface(ctx, s); err != nil {
		cancel()
		workers.Wait()
		return err
	}
	workers.Add(1)
	go func() { defer workers.Done(); done <- runForward(ctx, eylanForward(s)) }()
	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
	}
	cancel()
	workers.Wait()
	return err
}
