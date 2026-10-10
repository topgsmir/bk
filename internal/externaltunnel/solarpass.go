package externaltunnel

import (
	"context"
	"crypto/sha256"
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
)

const solarpassArchiveSHA256 = "536188d6c65768aecf65b4131772a519d0f7af7b18daeff7201e8a571e3cbaeb"

func checkSolarpassCore() error {
	b, err := os.ReadFile(filepath.Join(CoreDir, "solarpass-tun"))
	if err != nil {
		return fmt.Errorf("Solarpass core: %w", err)
	}
	hash := sha256.Sum256(b)
	if hex.EncodeToString(hash[:]) != solarpassRepairedSHA256 {
		return fmt.Errorf("Solarpass core needs the pinned Spoof repair")
	}
	return nil
}
func installSolarpass(ctx context.Context, out io.Writer) error {
	archive, err := bundledSolarpass()
	if err != nil {
		return err
	}
	hash := sha256.Sum256(archive)
	if hex.EncodeToString(hash[:]) != solarpassArchiveSHA256 {
		return fmt.Errorf("Solarpass archive checksum differs")
	}
	original, err := extractCore(archive, "solarpass-tun")
	if err != nil {
		return err
	}
	repaired, err := repairSolarpass(original)
	if err != nil {
		return err
	}
	for _, program := range []string{"ip", "iptables", "unshare", "mount", "ip6tables"} {
		if _, err = exec.LookPath(program); err != nil {
			if err = runCommand(ctx, command("apt-get", "update")); err != nil {
				return err
			}
			if err = runCommand(ctx, command("apt-get", "install", "-y", "iproute2", "iptables", "util-linux", "mount")); err != nil {
				return err
			}
			break
		}
	}
	if err = os.MkdirAll(CoreDir, 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(CoreDir, ".solarpass-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err = file.Write(repaired); err != nil {
		file.Close()
		return err
	}
	if err = file.Chmod(0755); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	if err = os.Rename(file.Name(), filepath.Join(CoreDir, "solarpass-tun")); err != nil {
		return err
	}
	fmt.Fprintln(out, "Pinned Solarpass core installed with verified Spoof return-port repair.")
	return nil
}
func solarpassConfig(s Spec) (map[string]any, error) {
	carrier := strings.TrimPrefix(s.Kind, "s3-")
	target := s.Target
	if s.Protocol == "udp" {
		target = net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort+3))
	}
	endpoint := net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port))
	if s.Side == "kharej" && carrier != "spoof" {
		endpoint = net.JoinHostPort(s.LocalIP, strconv.Itoa(s.Port))
	}
	iran := s.LocalIP
	if s.Side == "kharej" {
		iran = s.PeerIP
	}
	cfg := map[string]any{"name": "bk-" + s.Name, "side": s.Side, "endpoint": endpoint, "transport": carrier, "token": s.Secret, "iran_endpoint": net.JoinHostPort(iran, strconv.Itoa(s.SourcePort+1)), "log_level": "debug", "ports": []any{}}
	if s.Side == "iran" || carrier == "spoof" || carrier == "gamepass" {
		listen := s.SourcePort
		if s.Side == "kharej" && carrier == "spoof" {
			listen++
		}
		cfg["ports"] = []any{map[string]any{"listen_port": listen, "target": target, "tcp": s.Protocol == "tcp", "udp": s.Protocol == "udp"}}
	}
	if carrier == "spoof" {
		filter, err := solarpassPeerFilter(s.PeerIP)
		if err != nil {
			return nil, err
		}
		cfg["spoof"] = map[string]any{"peer_spoof_ip": filter}
	}
	return cfg, nil
}

// Solarpass runs its own global sysctl optimizer unconditionally. A private
// mount namespace keeps /proc/sys read-only for this child without changing any
// host setting or any original bk engine.
const solarpassMountGuard = `set -eu
mount --bind /proc/sys /proc/sys
mount -o remount,bind,ro /proc/sys
exec "$@"`

func solarpassFirewall(s Spec) []Command {
	var rules []Command
	add := func(port int, protocol, source string) {
		rules = append(rules, command("iptables", "-w", "-I", "INPUT", "1", "-p", protocol, "--dport", strconv.Itoa(port), "!", "-s", source, "-m", "comment", "--comment", "bk-s3-"+s.Name, "-j", "DROP"))
	}
	if s.Side == "iran" {
		for _, proto := range []string{"tcp", "udp"} {
			add(s.SourcePort, proto, "127.0.0.1")
		}
		if s.Kind == "s3-spoof" {
			add(s.SourcePort+1, "udp", s.PeerIP)
		}
		add(s.SourcePort+2, "tcp", s.PeerIP)
	} else {
		add(s.Port+1, "tcp", s.PeerIP)
	}
	syncPort := s.Port + 1
	if s.Side == "iran" {
		syncPort = s.SourcePort + 2
	}
	rules = append(rules, command("ip6tables", "-w", "-I", "INPUT", "1", "-p", "tcp", "--dport", strconv.Itoa(syncPort), "-m", "comment", "--comment", "bk-s3-"+s.Name, "-j", "DROP"))
	return rules
}
func runSolarpass(ctx context.Context, s Spec, dir string) error {
	// These two native carriers abort their first connection attempts instead
	// of retrying when Iran starts before the Kharej listener during a test wave.
	if s.Side == "iran" && (s.Kind == "s3-gamepass" || s.Kind == "s3-arenapass") {
		ready, cancel := context.WithTimeout(ctx, 45*time.Second)
		err := waitSolarpassCarrier(ready, net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)))
		cancel()
		if err != nil {
			return err
		}
	}
	cfg, err := solarpassConfig(s)
	if err != nil {
		return err
	}
	body, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	config := filepath.Join(dir, "solarpass.json")
	if err = os.WriteFile(config, body, 0600); err != nil {
		return err
	}
	var added []Command
	defer func() {
		for i := len(added) - 1; i >= 0; i-- {
			rule := added[i]
			args := append([]string(nil), rule.Args...)
			for j, a := range args {
				if a == "-I" {
					args[j] = "-D"
					args = append(args[:j+2], args[j+3:]...)
					break
				}
			}
			rule.Args = args
			cleanup(rule)
		}
	}()
	for _, rule := range solarpassFirewall(s) {
		if err = runCommand(ctx, rule); err != nil {
			return err
		}
		added = append(added, rule)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan error, 2)
	var workers sync.WaitGroup
	start := func(fn func() error) { workers.Add(1); go func() { defer workers.Done(); done <- fn() }() }
	start(func() error {
		return childRedacted(ctx, "env", []string{s.Secret}, "TOKIO_WORKER_THREADS=2", "unshare", "--mount", "--propagation", "private", "/bin/sh", "-c", solarpassMountGuard, "bk-solarpass", filepath.Join(CoreDir, "solarpass-tun"), "run", "--config", config)
	})
	if s.Side == "iran" {
		address := net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort))
		if s.Protocol == "tcp" {
			start(func() error {
				return serveTCP(ctx, s.Listen, func(ctx context.Context) (net.Conn, error) {
					return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", address)
				})
			})
		} else {
			start(func() error { return runSolarpassUDPBridge(ctx, s) })
		}
	}
	if s.Side == "kharej" && s.Protocol == "udp" {
		start(func() error { return runSolarpassUDPBridge(ctx, s) })
	}
	select {
	case err = <-done:
	case <-ctx.Done():
	}
	cancel()
	workers.Wait()
	return err
}

func waitSolarpassCarrier(ctx context.Context, address string) error {
	for {
		conn, err := (&net.Dialer{Timeout: time.Second}).DialContext(ctx, "tcp", address)
		if err == nil {
			conn.Close()
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("Solarpass initial carrier not ready: %w", ctx.Err())
		case <-time.After(150 * time.Millisecond):
		}
	}
}
