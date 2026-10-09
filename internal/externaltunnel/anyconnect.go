package externaltunnel

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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
)

func AnyConnect(kind string) bool { return strings.HasPrefix(kind, "e3-anyconnect-") }
func installAnyConnect(ctx context.Context, out io.Writer) error {
	// Extract the stock server executables without deploying or starting a
	// distribution-wide ocserv daemon or touching /etc/ocserv.
	libraries := []string{"iproute2", "openconnect", "dpkg", "libev4", "libhttp-parser2.9", "libmaxminddb0", "libnl-3-200", "libnl-route-3-200", "liboath0", "libprotobuf-c1", "libradcli4", "libseccomp2", "libtalloc2"}
	if err := runCommand(ctx, command("apt-get", "update")); err != nil {
		return err
	}
	if err := runCommand(ctx, command("apt-get", append([]string{"install", "-y", "--no-install-recommends"}, libraries...)...)); err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "bk-ocserv-package-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	cmd := exec.CommandContext(ctx, "apt-get", "download", "ocserv")
	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr = out, out
	if err = cmd.Run(); err != nil {
		return err
	}
	packages, _ := filepath.Glob(filepath.Join(dir, "ocserv_*.deb"))
	if len(packages) != 1 {
		return fmt.Errorf("expected one stock ocserv package")
	}
	extracted := filepath.Join(dir, "extracted")
	if err = runCommand(ctx, command("dpkg-deb", "-x", packages[0], extracted)); err != nil {
		return err
	}
	if err = os.MkdirAll(CoreDir, 0700); err != nil {
		return err
	}
	for name, source := range map[string]string{"ocserv": "usr/sbin/ocserv", "ocpasswd": "usr/bin/ocpasswd", "ocserv-worker": "usr/sbin/ocserv-worker"} {
		body, e := os.ReadFile(filepath.Join(extracted, source))
		if e != nil {
			return e
		}
		f, e := os.CreateTemp(CoreDir, ".ocserv-")
		if e != nil {
			return e
		}
		path := f.Name()
		if _, e = f.Write(body); e == nil {
			e = f.Chmod(0700)
		}
		closeErr := f.Close()
		if e == nil {
			e = closeErr
		}
		if e == nil {
			e = os.Rename(path, filepath.Join(CoreDir, name))
		}
		os.Remove(path)
		if e != nil {
			return e
		}
	}
	fmt.Fprintln(out, "Stock AnyConnect server installed as isolated bk executables.")
	return nil
}
func anyConnectConfig(s Spec, dir string) string {
	_, subnet, _ := net.ParseCIDR(s.IranIP)
	ones, _ := subnet.Mask.Size()
	return fmt.Sprintf("auth = \"plain[passwd=%s]\"\nlisten-host = %s\ntcp-port = %d\nudp-port = %d\nrun-as-user = root\nrun-as-group = root\nsocket-file = %s\npid-file = %s\nserver-cert = %s\nserver-key = %s\nmax-clients = 1\nmax-same-clients = 1\nkeepalive = 5\ndpd = 10\nmobile-dpd = 10\ntry-mtu-discovery = false\nmtu = %d\ncompression = false\ndevice = %s\nipv4-network = %s/%d\nroute = %s/%d\nrestrict-user-to-routes = false\ncisco-client-compat = true\ncookie-timeout = 60\nauth-timeout = 15\n", filepath.Join(dir, "passwd"), s.LocalIP, s.Port, s.Port, filepath.Join(dir, "ocserv.sock"), filepath.Join(dir, "ocserv.pid"), filepath.Join(dir, "ovpn-cert.pem"), filepath.Join(dir, "ovpn-key.pem"), s.MTU, s.Interface(), subnet.IP, ones, subnet.IP, ones)
}
func runAnyConnect(ctx context.Context, s Spec, dir string) error {
	// A single /30 supplies exactly the paired Iran .1 and Kharej .2 addresses.
	iran, subnet, _ := net.ParseCIDR(s.IranIP)
	peer, _, _ := net.ParseCIDR(s.KharejIP)
	ones, bits := subnet.Mask.Size()
	first := append(net.IP(nil), subnet.IP.To4()...)
	first[3]++
	second := append(net.IP(nil), subnet.IP.To4()...)
	second[3] += 2
	if ones != 30 || bits != 32 || !iran.Equal(first) || !peer.Equal(second) {
		return fmt.Errorf("AnyConnect requires paired .1/.2 addresses in one unused /30")
	}
	if err := eylanOpenVPNIdentity(s, dir); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan error, 2)
	var workers sync.WaitGroup
	var launch func() error
	iface := s.Interface()
	logPath := filepath.Join(dir, "native.log")
	if s.Side == "iran" {
		cmd := exec.CommandContext(ctx, filepath.Join(CoreDir, "ocpasswd"), "-c", filepath.Join(dir, "passwd"), "-g", "bk", "bk")
		cmd.Stdin = strings.NewReader(s.Secret + "\n" + s.Secret + "\n")
		if body, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("AnyConnect credential generation: %w: %s", err, body)
		}
		path := filepath.Join(dir, "ocserv.conf")
		if err := os.WriteFile(path, []byte(anyConnectConfig(s, dir)), 0600); err != nil {
			return err
		}
		if err := runCommand(ctx, command(filepath.Join(CoreDir, "ocserv"), "--test-config", "--config", path)); err != nil {
			return err
		}
		iface += "0"
		if _, err := net.InterfaceByName(iface); err == nil {
			return fmt.Errorf("AnyConnect interface %s already exists", iface)
		}
		launch = func() error {
			return childInputLogged(ctx, logPath, "", []string{s.Secret}, filepath.Join(CoreDir, "ocserv"), "--foreground", "--config", path, "--debug", "1")
		}
	} else {
		// Derive the server certificate from the same paired identity. Verify the
		// server's public-key pin; never use an insecure certificate exception.
		serverSpec := s.Mirror()
		serverDir := filepath.Join(dir, "server-identity")
		if err := os.Mkdir(serverDir, 0700); err != nil {
			return err
		}
		if err := eylanOpenVPNIdentity(serverSpec, serverDir); err != nil {
			return err
		}
		body, err := os.ReadFile(filepath.Join(serverDir, "ovpn-cert.pem"))
		if err != nil {
			return err
		}
		block, _ := pem.Decode(body)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
		pin := "pin-sha256:" + base64.StdEncoding.EncodeToString(sum[:])
		script := fmt.Sprintf("#!/bin/sh\nset -eu\ncase \"${reason:-}\" in\nconnect|reconnect) [ \"${TUNDEV:-}\" = %q ] || exit 1; ip addr replace %s dev %s; mtu=\"${INTERNAL_IP4_MTU:-%d}\"; case \"$mtu\" in *[!0-9]*|\"\") exit 1;; esac; [ \"$mtu\" -ge 576 ] && [ \"$mtu\" -le %d ] || exit 1; ip link set dev %s mtu \"$mtu\" up;;\n*) ;;\nesac\n", s.Interface(), s.TunnelIP(), s.Interface(), s.MTU, s.MTU, s.Interface())
		path := filepath.Join(dir, "configure-interface.sh")
		if err = os.WriteFile(path, []byte(script), 0700); err != nil {
			return err
		}
		args := []string{"--protocol=anyconnect", "--non-inter", "--passwd-on-stdin", "--user=bk", "--interface=" + iface, "--script=" + path, "--servercert=" + pin, "--no-system-trust", "--cafile=" + filepath.Join(serverDir, "ovpn-ca.pem"), "--no-proxy", "--disable-ipv6", "--compression=none", "--verbose", "--mtu=" + strconv.Itoa(s.MTU)}
		if s.Kind == "e3-anyconnect-tls" {
			args = append(args, "--no-dtls")
		}
		args = append(args, "https://"+net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)))
		launch = func() error {
			return childInputLogged(ctx, logPath, s.Secret+"\n", []string{s.Secret}, "openconnect", args...)
		}
	}
	if s.Side == "kharej" {
		ready, stop := context.WithTimeout(ctx, 30*time.Second)
		err := waitSolarpassCarrier(ready, net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)))
		stop()
		if err != nil {
			return fmt.Errorf("AnyConnect peer listener not ready: %w", err)
		}
	}
	workers.Add(1)
	go func() { defer workers.Done(); done <- launch() }()
	defer func() {
		if _, err := net.InterfaceByName(iface); err == nil {
			cleanup(command("ip", "link", "del", iface))
		}
	}()
	if err := waitVPNAddress(ctx, iface, strings.Split(s.TunnelIP(), "/")[0]); err != nil {
		cancel()
		workers.Wait()
		return err
	}
	if s.Kind == "e3-anyconnect-dtls" && s.Side == "kharej" {
		ready, stop := context.WithTimeout(ctx, 20*time.Second)
		for {
			body, _ := os.ReadFile(logPath)
			if strings.Contains(string(body), "Established DTLS connection") {
				for _, line := range strings.Split(string(body), "\n") {
					if strings.Contains(line, "MTU") || strings.Contains(line, "Established DTLS") {
						fmt.Fprintln(os.Stderr, line)
					}
				}
				break
			}
			select {
			case <-ready.Done():
				stop()
				cancel()
				workers.Wait()
				return fmt.Errorf("AnyConnect DTLS was not established; no TLS fallback counted as success")
			case <-time.After(100 * time.Millisecond):
			}
		}
		stop()
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
func waitVPNAddress(ctx context.Context, iface, address string) error {
	ready, stop := context.WithTimeout(ctx, 25*time.Second)
	defer stop()
	for {
		if device, err := net.InterfaceByName(iface); err == nil {
			addresses, _ := device.Addrs()
			for _, a := range addresses {
				ip, _, _ := net.ParseCIDR(a.String())
				if ip != nil && ip.String() == address {
					return nil
				}
			}
		}
		select {
		case <-ready.Done():
			return fmt.Errorf("VPN interface %s did not become ready: %w", iface, ready.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}
func childInputLogged(ctx context.Context, logPath, input string, secrets []string, bin string, args ...string) error {
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = strings.NewReader(input)
	r, w := io.Pipe()
	cmd.Stdout, cmd.Stderr = w, w
	drained := make(chan struct{})
	go func() { redactLines(io.MultiWriter(file, os.Stderr), r, secrets); close(drained) }()
	defer func() { w.Close(); <-drained; r.Close() }()
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 5 * time.Second
	err = cmd.Run()
	if ctx.Err() != nil {
		return nil
	}
	if err == nil {
		return fmt.Errorf("%s stopped unexpectedly", filepath.Base(bin))
	}
	return err
}
