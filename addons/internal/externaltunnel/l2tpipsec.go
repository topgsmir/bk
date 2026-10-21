package externaltunnel

import (
	"context"
	"encoding/base64"
	"encoding/hex"
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

func L2TPIPsec(kind string) bool { return kind == "e3-l2tp-ipsec" }
func installL2TPIPsec(ctx context.Context, out io.Writer) error {
	packages := []string{"iproute2", "iptables", "ppp", "dpkg", "strongswan-libcharon", "libstrongswan-standard-plugins", "libcharon-extra-plugins", "strongswan-swanctl"}
	if err := runCommand(ctx, command("apt-get", "update")); err != nil {
		return err
	}
	if err := runCommand(ctx, command("apt-get", append([]string{"install", "-y", "--no-install-recommends"}, packages...)...)); err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "bk-ipsec-packages-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	for _, name := range []string{"strongswan-charon", "xl2tpd", "strongswan-swanctl"} {
		cmd := exec.CommandContext(ctx, "apt-get", "download", name)
		cmd.Dir = dir
		cmd.Stdout, cmd.Stderr = out, out
		if err = cmd.Run(); err != nil {
			return err
		}
		matches, _ := filepath.Glob(filepath.Join(dir, name+"_*.deb"))
		if len(matches) != 1 {
			return fmt.Errorf("expected one stock %s package", name)
		}
		if err = runCommand(ctx, command("dpkg-deb", "-x", matches[0], filepath.Join(dir, "extracted"))); err != nil {
			return err
		}
	}
	if err = os.MkdirAll(CoreDir, 0700); err != nil {
		return err
	}
	for name, source := range map[string]string{"charon": "usr/lib/ipsec/charon", "xl2tpd": "usr/sbin/xl2tpd", "swanctl": "usr/sbin/swanctl"} {
		body, e := os.ReadFile(filepath.Join(dir, "extracted", source))
		if e != nil {
			return e
		}
		f, e := os.CreateTemp(CoreDir, ".ipsec-")
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
	fmt.Fprintln(out, "Stock L2TP/IPsec engines installed without deploying a global VPN daemon.")
	return nil
}
func l2tpIPsecPreflight(s Spec) error {
	for _, port := range []int{500, 4500, 1701} {
		c, err := net.ListenPacket("udp4", net.JoinHostPort("0.0.0.0", strconv.Itoa(port)))
		if err != nil {
			return fmt.Errorf("L2TP/IPsec needs unused UDP %d; existing VPN services are preserved: %w", port, err)
		}
		c.Close()
	}
	// Refuse a conflicting owner before installing any kernel state. Never log
	// XFRM state JSON: it contains live encryption keys.
	body, err := exec.Command("ip", "xfrm", "state", "list", "nokeys", "reqid", strconv.Itoa(s.ID)).Output()
	if err != nil {
		return fmt.Errorf("IPsec XFRM is unavailable: %w", err)
	}
	if strings.TrimSpace(string(body)) != "" {
		return fmt.Errorf("IPsec reqid %d already exists", s.ID)
	}
	body, err = exec.Command("ip", "xfrm", "policy", "list").Output()
	if err != nil {
		return err
	}
	if len(ownedIPsecPolicies(string(body), s.ID)) > 0 {
		return fmt.Errorf("IPsec policy reqid %d already exists", s.ID)
	}
	return nil
}
func l2tpIPsecConfig(s Spec, dir string) string {
	key, _ := base64.StdEncoding.DecodeString(s.Key("eylan/ipsec/psk"))
	other := "kharej"
	if s.Side == "kharej" {
		other = "iran"
	}
	return fmt.Sprintf(`connections {
 bk {
  version = 2
  local_addrs = %s
  remote_addrs = %s
  proposals = aes256gcm16-prfsha256-ecp256
  mobike = no
  encap = yes
  dpd_delay = 10s
  local { auth = psk
   id = bk-%s
  }
  remote { auth = psk
   id = bk-%s
  }
  children {
   bk-l2tp {
    mode = transport
    local_ts = %s[udp/1701]
    remote_ts = %s[udp/1701]
    esp_proposals = aes256gcm16
    reqid = %d
    start_action = none
    dpd_action = clear
   }
  }
 }
}
secrets {
 ike-bk { id-1 = bk-iran
  id-2 = bk-kharej
  secret = 0x%s
 }
}
`, s.LocalIP, s.PeerIP, s.Side, other, s.LocalIP, s.PeerIP, s.ID, hex.EncodeToString(key))
}
func l2tpPPPConfig(s Spec) string {
	text := fmt.Sprintf("ifname %s\nmtu %d\nmru %d\nnodefaultroute\nnoipdefault\nnoipv6\nnoauth\nnoccp\nnoaccomp\nnopcomp\nlcp-echo-interval 5\nlcp-echo-failure 3\nname bk\n", s.Interface(), s.MTU, s.MTU)
	if s.Side == "iran" {
		text += "require-chap\nrefuse-pap\nrefuse-mschap\nrefuse-mschap-v2\nrefuse-eap\n"
	} else {
		text += "user bk\npassword " + s.Secret + "\n"
	}
	return text
}
func l2tpDaemonConfig(s Spec, dir string) string {
	text := fmt.Sprintf("[global]\nport = 1701\nlisten-addr = %s\nipsec saref = no\naccess control = yes\n", s.LocalIP)
	if s.Side == "iran" {
		text += fmt.Sprintf("[lns bk]\nlocal ip = %s\nip range = %s-%s\nlac = %s\nrequire chap = yes\nrefuse pap = yes\nrequire authentication = yes\nname = bk\npppoptfile = %s\nlength bit = yes\n", strings.Split(s.IranIP, "/")[0], strings.Split(s.KharejIP, "/")[0], strings.Split(s.KharejIP, "/")[0], s.PeerIP, filepath.Join(dir, "ppp.options"))
	} else {
		text += fmt.Sprintf("[lac bk]\nlns = %s\nname = bk\nrequire authentication = no\nrefuse pap = yes\npppoptfile = %s\nautodial = yes\nredial = yes\nredial timeout = 2\nlength bit = yes\n", s.PeerIP, filepath.Join(dir, "ppp.options"))
	}
	return text
}

const l2tpPrivateMounts = `set -eu
mount --bind "$1/run" /run
mount --bind "$1/ppp" /etc/ppp
mount --bind /proc/sys /proc/sys
mount -o remount,bind,ro /proc/sys
shift
exec "$@"`

func runL2TPIPsec(ctx context.Context, s Spec, dir string) error {
	if err := l2tpIPsecPreflight(s); err != nil {
		return err
	}
	for _, name := range []string{"run", "ppp"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0700); err != nil {
			return err
		}
	}
	for _, name := range []string{"ip-up", "ip-down", "ipv6-up", "ipv6-down"} {
		if err := os.WriteFile(filepath.Join(dir, "ppp", name), []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
			return err
		}
	}
	strong := fmt.Sprintf("charon {\n load_modular = no\n load = random nonce openssl kdf kernel-netlink socket-default vici\n threads = 16\n plugins { vici { socket = unix://%s\n }\n kernel-netlink { install_routes = no\n }\n }\n}\nswanctl { load = random nonce openssl\n}\n", filepath.Join(dir, "charon.vici"))
	for name, body := range map[string]string{"strongswan.conf": strong, "swanctl.conf": l2tpIPsecConfig(s, dir), "xl2tpd.conf": l2tpDaemonConfig(s, dir), "ppp.options": l2tpPPPConfig(s), "ppp/options": "", "ppp/chap-secrets": "bk * " + s.Secret + " *\n"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0600); err != nil {
			return err
		}
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var workers sync.WaitGroup
	done := make(chan error, 3)
	launch := func(fn func() error) { workers.Add(1); go func() { defer workers.Done(); done <- fn() }() }
	// Accept L2TP only if it passed this tunnel's ESP policy. The rule is
	// exclusively owned and removed exactly; administrator rules are preserved.
	rule := command("iptables", "-w", "-I", "INPUT", "1", "-p", "udp", "-d", s.LocalIP, "--dport", "1701", "-m", "policy", "--dir", "in", "--pol", "none", "-m", "comment", "--comment", "bk-e3-"+s.Name, "-j", "DROP")
	if err := runCommand(ctx, rule); err != nil {
		return err
	}
	defer cleanup(command("iptables", append([]string{"-w", "-D", "INPUT"}, rule.Args[4:]...)...))
	defer func() {
		if _, err := net.InterfaceByName(s.Interface()); err == nil {
			cleanup(command("ip", "link", "del", s.Interface()))
		}
		cleanupL2TPIPsec(s)
	}()
	env := []string{"STRONGSWAN_CONF=" + filepath.Join(dir, "strongswan.conf")}
	launch(func() error {
		return childRedacted(ctx, "env", []string{s.Secret}, append(env, "unshare", "--mount", "--propagation", "private", "/bin/sh", "-c", l2tpPrivateMounts, "bk-ipsec", dir, filepath.Join(CoreDir, "charon"))...)
	})
	uri := "unix://" + filepath.Join(dir, "charon.vici")
	ready, stop := context.WithTimeout(ctx, 20*time.Second)
	defer stop()
	for {
		if _, err := os.Stat(filepath.Join(dir, "charon.vici")); err == nil {
			break
		}
		select {
		case err := <-done:
			cancel()
			workers.Wait()
			return fmt.Errorf("IPsec engine stopped: %w", err)
		case <-ready.Done():
			cancel()
			workers.Wait()
			return fmt.Errorf("IPsec control socket not ready: %w", ready.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
	swan := func(args ...string) error {
		return childCommandRedacted(ctx, "env", []string{s.Secret}, append(env, append([]string{filepath.Join(CoreDir, "swanctl")}, args...)...)...)
	}
	if err := swan("--load-all", "--file", filepath.Join(dir, "swanctl.conf"), "--uri", uri, "--noprompt"); err != nil {
		cancel()
		workers.Wait()
		return err
	}
	if s.Side == "kharej" {
		if err := swan("--initiate", "--child", "bk-l2tp", "--uri", uri, "--timeout", "20"); err != nil {
			cancel()
			workers.Wait()
			return err
		}
	}
	launch(func() error {
		return childRedacted(ctx, "unshare", []string{s.Secret}, "--mount", "--propagation", "private", "/bin/sh", "-c", l2tpPrivateMounts, "bk-l2tp", dir, filepath.Join(CoreDir, "xl2tpd"), "-D", "-c", filepath.Join(dir, "xl2tpd.conf"), "-p", filepath.Join(dir, "run", "xl2tpd.pid"), "-C", filepath.Join(dir, "run", "l2tp-control"))
	})
	if err := waitVPNAddress(ctx, s.Interface(), strings.Split(s.TunnelIP(), "/")[0]); err != nil {
		cancel()
		workers.Wait()
		return err
	}
	launch(func() error { return runForward(ctx, eylanForward(s)) })
	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
	}
	cancel()
	workers.Wait()
	return err
}

// A short control command may exit successfully; unlike a supervised tunnel
// engine, that is expected. Secrets are never passed to its command line.
func childCommandRedacted(ctx context.Context, bin string, secrets []string, args ...string) error {
	ctx, stop := context.WithTimeout(ctx, 30*time.Second)
	defer stop()
	cmd := exec.CommandContext(ctx, bin, args...)
	body, err := cmd.CombinedOutput()
	if err != nil {
		message := string(body)
		for _, secret := range secrets {
			message = strings.ReplaceAll(message, secret, "[redacted]")
		}
		return fmt.Errorf("%s: %w: %s", filepath.Base(bin), err, message)
	}
	return nil
}
func cleanupL2TPIPsec(s Spec) {
	// The unique reqid was verified absent before startup. Select it explicitly;
	// do not flush shared state and never request plaintext encryption keys.
	cleanup(command("ip", "xfrm", "state", "deleteall", "reqid", strconv.Itoa(s.ID)))
	body, err := exec.Command("ip", "xfrm", "policy", "list").Output()
	if err != nil {
		return
	}
	for _, p := range ownedIPsecPolicies(string(body), s.ID) {
		cleanup(command("ip", "xfrm", "policy", "delete", "dir", p.dir, "index", p.index))
	}
}

type ipsecPolicyOwner struct{ dir, index string }

func ownedIPsecPolicies(body string, id int) []ipsecPolicyOwner {
	var result []ipsecPolicyOwner
	var block []string
	consume := func() {
		var reqid, dir, index string
		for i := 0; i+1 < len(block); i++ {
			switch block[i] {
			case "reqid":
				reqid = block[i+1]
			case "dir":
				dir = block[i+1]
			case "index":
				index = block[i+1]
			}
		}
		if reqid != strconv.Itoa(id) || (dir != "in" && dir != "out" && dir != "fwd") {
			return
		}
		if _, err := strconv.ParseUint(index, 0, 32); err != nil {
			return
		}
		result = append(result, ipsecPolicyOwner{dir, index})
	}
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "src ") {
			consume()
			block = nil
		}
		block = append(block, strings.Fields(line)...)
	}
	consume()
	return result
}
