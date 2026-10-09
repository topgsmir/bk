package externaltunnel

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
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
	"strings"
	"time"
)

const daggerCoreSHA256 = "fccb58880f6dfd53e74349b5780f662b3b2c9f1fea6c21d9e3c738b830457eb7"

func checkDaggerCore() error {
	body, err := os.ReadFile(filepath.Join(CoreDir, "dagger-rs"))
	if err != nil {
		return err
	}
	sum := sha256.Sum256(body)
	if hex.EncodeToString(sum[:]) != daggerCoreSHA256 {
		return fmt.Errorf("Dagger installed core checksum differs from the pinned core")
	}
	return nil
}

// A setup link carries the shared root secret. Domain-separated Noise identities
// and TLS material are derived locally; no operational private-key file is copied.
func daggerKey(s Spec, side string) string {
	b, _ := base64.StdEncoding.DecodeString(s.Key("dagger/" + side))
	return hex.EncodeToString(b)
}
func daggerPublic(s Spec, side string) string {
	b, _ := base64.StdEncoding.DecodeString(s.AWGPublic("dagger/" + side))
	return hex.EncodeToString(b)
}
func daggerTLS(s Spec, dir string) error {
	seed, _ := base64.StdEncoding.DecodeString(s.Key("dagger/tls"))
	private := ed25519.NewKeyFromSeed(seed)
	hash := sha256.Sum256(seed)
	template := &x509.Certificate{SerialNumber: new(big.Int).SetBytes(hash[:16]), Subject: pkix.Name{CommonName: "bk-dagger.local"}, DNSNames: []string{"bk-dagger.local"}, NotBefore: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), NotAfter: time.Date(2120, 1, 1, 0, 0, 0, 0, time.UTC), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IsCA: false, BasicConstraintsValid: true}
	cert, err := x509.CreateCertificate(nil, template, template, private.Public(), private)
	if err != nil {
		return err
	}
	key, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(dir, "tls.pem"), pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}), 0600); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "tls-key.pem"), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}), 0600)
}
func daggerConfig(s Spec, dir string) (map[string]any, error) {
	carrier, profile, upload, reverse := DaggerOptions(s.Kind)
	side, other, role := s.Side, "kharej", "server"
	if side == "kharej" {
		other, role = "iran", "client"
	}
	keyfile := filepath.Join(dir, "noise.key")
	if e := os.WriteFile(keyfile, []byte(daggerKey(s, side)), 0600); e != nil {
		return nil, e
	}
	config := map[string]any{"mode": role, "reverse": reverse, "private_key_file": keyfile, "heartbeat_sec": 2, "dead_timeout_sec": 15, "max_connections": 16, "max_streams": 256}
	addr := net.JoinHostPort(s.LocalIP, fmt.Sprint(s.Port))
	dial := net.JoinHostPort(s.PeerIP, fmt.Sprint(s.Port))
	if carrier == "dc6" {
		if net.ParseIP(s.LocalIPv6) == nil || net.ParseIP(s.PeerIPv6) == nil {
			return nil, fmt.Errorf("Dagger DC6 requires both assigned IPv6 addresses")
		}
		addr, dial = net.JoinHostPort(s.LocalIPv6, fmt.Sprint(s.Port)), net.JoinHostPort(s.PeerIPv6, fmt.Sprint(s.Port))
	}
	endpoint := map[string]any{"transport": carrier, "connection_pool": 1}
	if side == "iran" {
		endpoint["addr"] = addr
		if reverse {
			endpoint["addr"] = dial
		}
		endpoint["maps"] = []any{map[string]any{"type": s.Protocol, "bind": s.Listen, "target": s.Target}}
		config["peer_public_keys"] = []string{daggerPublic(s, other)}
		config["listeners"] = []any{endpoint}
	} else {
		endpoint["addr"] = dial
		if reverse {
			endpoint["addr"] = addr
		}
		endpoint["server_public_key"] = daggerPublic(s, other)
		endpoint["retry_interval"] = 1
		endpoint["dial_timeout"] = 15
		config["paths"] = []any{endpoint}
		config["allowed_targets"] = []string{s.Target}
	}
	if upload != "" {
		endpoint["xhttp"] = map[string]any{"mode": upload, "up_concurrency": 4}
		endpoint["http_path"] = "/bk-dagger"
	}
	if carrier == "https" || carrier == "wss" || carrier == "xhttps" {
		if e := daggerTLS(s, dir); e != nil {
			return nil, e
		}
		listening := (side == "iran") != reverse
		if listening {
			endpoint["cert_file"] = filepath.Join(dir, "tls.pem")
			endpoint["key_file"] = filepath.Join(dir, "tls-key.pem")
		} else {
			endpoint["ca_file"] = filepath.Join(dir, "tls.pem")
			endpoint["server_name"] = "bk-dagger.local"
		}
	}
	if profile != "" {
		clientPort := s.SourcePort
		if clientPort == 0 || clientPort == s.Port {
			clientPort = s.Port + 1
			if clientPort > 65535 {
				clientPort = s.Port - 1
			}
		}
		// Distinct identities exclude the kernel's reflected ICMP echo reply.
		raw := map[string]any{"profile": profile, "local_ip": s.LocalIP, "peer_ip": s.PeerIP, "l4_port": s.Port, "source_port": clientPort}
		if profile == "dcpi" {
			raw["profile"] = "tcp"
			raw["dcpi_mode"] = true
		}
		if s.WAN != "" {
			raw["interface"] = s.WAN
		}
		if s.RouterMAC != "" {
			raw["peer_mac"] = s.RouterMAC
		}
		endpoint["raw"] = raw
	}
	if carrier == "quantum" || carrier == "quantum-gaming" || carrier == "quantum+" {
		endpoint["quantum"] = map[string]any{"mtu": s.MTU, "knock": false}
	}
	if carrier == "tun" {
		hash := sha256.Sum256([]byte(s.Name + "/dagger"))
		iface := "bk3" + hex.EncodeToString(hash[:5])
		if _, e := net.InterfaceByName(iface); e == nil {
			return nil, fmt.Errorf("Dagger interface %s already exists", iface)
		}
		local, peer := s.IranIP, s.KharejIP
		if s.Side == "kharej" {
			local, peer = peer, local
		}
		config["tun"] = map[string]any{"name": iface, "address": strings.Split(local, "/")[0], "peer_address": strings.Split(peer, "/")[0], "mtu": s.MTU, "forwarding_port": s.SourcePort}
	}
	return config, nil
}
func runDagger(ctx context.Context, s Spec, dir string) error {
	config, e := daggerConfig(s, dir)
	if e != nil {
		return e
	}
	b, e := json.Marshal(config)
	if e != nil {
		return e
	}
	path := filepath.Join(dir, "dagger.json")
	if e = os.WriteFile(path, b, 0600); e != nil {
		return e
	}
	core := filepath.Join(CoreDir, "dagger-rs")
	check := exec.CommandContext(ctx, core, "--config", path, "--check")
	if b, e := check.CombinedOutput(); e != nil {
		return fmt.Errorf("Dagger configuration: %w: %s", e, strings.TrimSpace(string(b)))
	}
	cmd := exec.CommandContext(ctx, core, "--config", path)
	cmd.Env = append(os.Environ(), "TOKIO_WORKER_THREADS=2", "RUST_LOG=debug")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	// Graceful cancellation is required for the core to remove its own TUN device.
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 5 * time.Second
	e = cmd.Run()
	if ctx.Err() != nil {
		return nil
	}
	return e
}

// Retain all license texts delivered with the pinned core. Only regular notice
// files below known archive roots can be installed; archive paths are not trusted.
func installDaggerNotices(archive []byte) error {
	g, e := gzip.NewReader(bytes.NewReader(archive))
	if e != nil {
		return e
	}
	defer g.Close()
	t := tar.NewReader(g)
	for {
		h, e := t.Next()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		name := strings.TrimPrefix(h.Name, "dagger-rs-linux/")
		if name != "LICENSE" && !strings.HasPrefix(name, "licenses/") {
			continue
		}
		clean := filepath.Clean(name)
		if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "\\") || h.Size > 2<<20 {
			return fmt.Errorf("invalid Dagger notice path or size")
		}
		path := filepath.Join(CoreDir, "dagger-notices", clean)
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			return e
		}
		body, e := io.ReadAll(io.LimitReader(t, 2<<20))
		if e != nil {
			return e
		}
		tmp, e := os.CreateTemp(filepath.Dir(path), ".notice-")
		if e != nil {
			return e
		}
		_, e = tmp.Write(body)
		if e == nil {
			e = tmp.Chmod(0644)
		}
		closeErr := tmp.Close()
		if e == nil {
			e = closeErr
		}
		if e == nil {
			e = os.Rename(tmp.Name(), path)
		}
		os.Remove(tmp.Name())
		if e != nil {
			return e
		}
	}
	return nil
}
