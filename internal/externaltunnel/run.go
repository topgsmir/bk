package externaltunnel

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/topgsmir/bk/internal/tunnel/l3"
)

var CoreDir = "/etc/bk/external/cores"

func Check(s Spec) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if runtime.GOOS != "linux" {
		return fmt.Errorf("additional tunnels require Linux")
	}
	programs := []string{}
	if Layer3(s.Kind) {
		programs = append(programs, "ip")
	}
	switch s.Kind {
	case "awg":
		programs = append(programs, awgTool())
	case "paqet":
		programs = append(programs, "ip", "iptables", filepath.Join(CoreDir, "paqet"))
	case "rgt-tcp", "rgt-udp":
		programs = append(programs, filepath.Join(CoreDir, "rgt"))
	case "ssh":
		if s.Side == "iran" {
			if _, e := os.Stat(s.SSHKey); e != nil {
				return fmt.Errorf("SSH key missing; use Additional Tunnels -> Prepare SSH Authentication")
			}
			if _, e := os.Stat(s.KnownHosts); e != nil {
				return fmt.Errorf("SSH known hosts missing; verify the peer fingerprint first")
			}
		}
	case "alghadir":
		programs = append(programs, "ip", "iptables", "obfs4proxy", filepath.Join(CoreDir, "udp2raw"))
	}
	for _, p := range programs {
		if _, err := exec.LookPath(p); err != nil {
			return fmt.Errorf("missing %s; use Additional Tunnels -> Install Dependencies", p)
		}
	}
	if Layer3(s.Kind) && s.Kind != "alghadir" {
		_, subnet, _ := net.ParseCIDR(s.TunnelIP())
		interfaces, e := net.Interfaces()
		if e != nil {
			return e
		}
		for _, iface := range interfaces {
			addresses, e := iface.Addrs()
			if e != nil {
				return e
			}
			for _, address := range addresses {
				_, existing, e := net.ParseCIDR(address.String())
				if e == nil && existing.IP.To4() != nil && (subnet.Contains(existing.IP) || existing.Contains(subnet.IP)) {
					return fmt.Errorf("tunnel subnet overlaps interface %s; choose an unused subnet", iface.Name)
				}
			}
		}
	}
	return nil
}
func runCommand(ctx context.Context, c Command) error {
	cmd := exec.CommandContext(ctx, c.Program, c.Args...)
	out, e := cmd.CombinedOutput()
	if e != nil {
		return fmt.Errorf("%s: %w: %s", c.Program, e, strings.TrimSpace(string(out)))
	}
	return nil
}
func Run(ctx context.Context, s Spec) error {
	if err := Check(s); err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "bk-ext-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if s.Kind == "alghadir" {
		return runAlghadir(ctx, s, dir)
	}
	if s.Kind == "ssh" {
		if s.Side == "kharej" {
			<-ctx.Done()
			return nil
		}
		return runSSH(ctx, s)
	}
	if Reverse(s.Kind) {
		path := filepath.Join(dir, "rgt.toml")
		if err := os.WriteFile(path, []byte(s.RGTConfig()), 0600); err != nil {
			return err
		}
		return child(ctx, filepath.Join(CoreDir, "rgt"), path)
	}
	if s.Kind == "paqet" {
		if err := s.detectRoute(); err != nil {
			return err
		}
		port := s.SourcePort
		if s.Side == "kharej" {
			port = s.Port
		}
		if port == 0 {
			l, e := net.Listen("tcp", net.JoinHostPort(s.LocalIP, "0"))
			if e != nil {
				return e
			}
			port = l.Addr().(*net.TCPAddr).Port
			l.Close()
			s.SourcePort = port
		}
		rules := paqetRules(s, port)
		var added []Command
		defer func() {
			for i := len(added) - 1; i >= 0; i-- {
				c := added[i]
				for j, a := range c.Args {
					if a == "-A" {
						c.Args[j] = "-D"
						break
					}
				}
				cleanup(c)
			}
		}()
		for _, c := range rules {
			if err := runCommand(ctx, c); err != nil {
				return err
			}
			added = append(added, c)
		}
		path := filepath.Join(dir, "paqet.yaml")
		if err := os.WriteFile(path, []byte(s.PaqetConfig()), 0600); err != nil {
			return err
		}
		return child(ctx, filepath.Join(CoreDir, "paqet"), "run", "-c", path)
	}
	up, down := s.KernelCommands()
	// Never adopt or delete an interface created by another program.
	if _, e := net.InterfaceByName(s.Interface()); e == nil {
		return fmt.Errorf("interface %s already exists; stop the tunnel before starting it again", s.Interface())
	}
	created := false
	var awgStopped <-chan struct{}
	defer func() {
		if created {
			if s.Kind == "awg" {
				if _, e := net.InterfaceByName(s.Interface()); e != nil {
					return
				}
			}
			for _, c := range down {
				cleanup(c)
			}
		}
	}()
	for i, c := range up {
		if e := runCommand(ctx, c); e != nil {
			if s.Kind == "awg" && i == 0 {
				stop, stopped, userspaceErr := startAWGUserspace(ctx, s.Interface())
				if userspaceErr != nil {
					return fmt.Errorf("kernel AWG: %v; userspace: %w", e, userspaceErr)
				}
				defer stop()
				awgStopped = stopped
			} else {
				return e
			}
		}
		if i == 0 {
			created = true
		}
		if s.Kind == "awg" && i == 0 {
			path := filepath.Join(dir, "awg.conf")
			if e := os.WriteFile(path, []byte(s.AWGConfig()), 0600); e != nil {
				return e
			}
			if e := runCommand(ctx, command(awgTool(), "setconf", s.Interface(), path)); e != nil {
				return e
			}
		}
	}
	if awgStopped != nil {
		return superviseAWGForward(ctx, s, awgStopped)
	}
	return runForward(ctx, s)
}
func cleanup(c Command) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if e := runCommand(ctx, c); e != nil {
		fmt.Fprintln(os.Stderr, "cleanup:", e)
	}
}
func child(ctx context.Context, bin string, args ...string) error {
	return childRedacted(ctx, bin, nil, args...)
}
func childRedacted(ctx context.Context, bin string, secrets []string, args ...string) error {
	cmd := exec.Command(bin, args...)
	if len(secrets) > 0 {
		r, w := io.Pipe()
		cmd.Stdout, cmd.Stderr = w, w
		drained := make(chan struct{})
		go func() { redactLines(os.Stderr, r, secrets); close(drained) }()
		defer func() { w.Close(); <-drained; r.Close() }()
	} else {
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case e := <-done:
		if e == nil && ctx.Err() == nil {
			return fmt.Errorf("%s stopped unexpectedly", filepath.Base(bin))
		}
		return e
	case <-ctx.Done():
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
			<-done
		}
		return nil
	}
}
func runForward(ctx context.Context, s Spec) error {
	target := s.Target
	if s.Side == "iran" {
		_, port, _ := net.SplitHostPort(target)
		target = net.JoinHostPort(s.PeerTunnelIP(), port)
	} else {
		// Kharej accepts only on its tunnel address and forwards to the real local backend.
		_, port, _ := net.SplitHostPort(target)
		s.Listen = net.JoinHostPort(strings.Split(s.TunnelIP(), "/")[0], port)
		if s.Backend != "" {
			target = s.Backend
		}
	}
	f, e := l3.NewForwarder(l3.Config{Ports: []string{s.Listen + "=" + target}, AcceptUDP: s.Protocol == "udp"}, nil)
	if e != nil {
		return e
	}
	return f.Run(ctx)
}
func paqetRules(s Spec, port int) []Command {
	comment := "bk-ext-" + s.Name
	return []Command{
		command("iptables", "-w", "-t", "raw", "-A", "PREROUTING", "-d", s.LocalIP, "-s", s.PeerIP, "-p", "tcp", "--dport", strconv.Itoa(port), "-m", "comment", "--comment", comment, "-j", "NOTRACK"),
		command("iptables", "-w", "-t", "raw", "-A", "OUTPUT", "-s", s.LocalIP, "-d", s.PeerIP, "-p", "tcp", "--sport", strconv.Itoa(port), "-m", "comment", "--comment", comment, "-j", "NOTRACK"),
		command("iptables", "-w", "-t", "mangle", "-A", "OUTPUT", "-s", s.LocalIP, "-d", s.PeerIP, "-p", "tcp", "--sport", strconv.Itoa(port), "--tcp-flags", "RST", "RST", "-m", "comment", "--comment", comment, "-j", "DROP"),
	}
}
func (s *Spec) detectRoute() error {
	cmd := exec.Command("ip", "-j", "route", "get", s.PeerIP)
	b, e := cmd.Output()
	if e != nil {
		return fmt.Errorf("route to peer: %w", e)
	}
	var routes []struct{ Dev, Gateway, Prefsrc string }
	if e = json.Unmarshal(b, &routes); e != nil || len(routes) == 0 {
		return fmt.Errorf("cannot determine route to peer")
	}
	r := routes[0]
	if s.WAN == "" {
		s.WAN = r.Dev
	}
	if s.LocalIP == "0.0.0.0" {
		s.LocalIP = r.Prefsrc
	}
	if s.RouterMAC == "" {
		hop := r.Gateway
		if hop == "" {
			hop = s.PeerIP
		}
		// A bounded probe populates the neighbour cache; it need not answer ICMP.
		probe, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = exec.CommandContext(probe, "ping", "-c", "1", "-W", "1", hop).Run()
		cancel()
		b, e = exec.Command("ip", "-j", "neigh", "show", "to", hop, "dev", s.WAN).Output()
		if e != nil {
			return e
		}
		var neigh []struct{ Lladdr string }
		if e = json.Unmarshal(b, &neigh); e != nil || len(neigh) == 0 || neigh[0].Lladdr == "" {
			return fmt.Errorf("cannot resolve next-hop MAC for %s on %s; enter WAN and router MAC in advanced setup", hop, s.WAN)
		}
		s.RouterMAC = neigh[0].Lladdr
	}
	return s.Validate()
}

// Relay closes both streams on cancellation and waits for both copies.
func relay(ctx context.Context, a, b net.Conn) {
	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(a, b); done <- struct{}{} }()
	go func() { _, _ = io.Copy(b, a); done <- struct{}{} }()
	finished := 0
	select {
	case <-ctx.Done():
	case <-done:
		finished = 1
	}
	a.Close()
	b.Close()
	for finished < 2 {
		<-done
		finished++
	}
}
func serveTCP(ctx context.Context, listen string, dial func(context.Context) (net.Conn, error)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	l, e := net.Listen("tcp", listen)
	if e != nil {
		return e
	}
	defer l.Close()
	go func() { <-ctx.Done(); l.Close() }()
	var wg sync.WaitGroup
	defer func() { cancel(); wg.Wait() }()
	for {
		a, e := l.Accept()
		if e != nil {
			if ctx.Err() != nil {
				return nil
			}
			return e
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			b, e := dial(ctx)
			if e != nil {
				a.Close()
				return
			}
			relay(ctx, a, b)
		}()
	}
}
func Load(path string) (Spec, error) {
	var s Spec
	b, e := os.ReadFile(path)
	if e != nil {
		return s, e
	}
	e = json.Unmarshal(b, &s)
	if e != nil {
		return s, e
	}
	return s, s.Validate()
}
func Save(path string, s Spec) error {
	if e := s.Validate(); e != nil {
		return e
	}
	b, e := json.MarshalIndent(s, "", "  ")
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		return e
	}
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	_, e = f.Write(b)
	closeErr := f.Close()
	return errors.Join(e, closeErr)
}

// Cores such as udp2raw echo their entire command line at startup. Filter the
// session secret before any child output reaches systemd or test logs.
func redactLines(out io.Writer, input io.Reader, secrets []string) {
	r := bufio.NewReader(input)
	for {
		line, e := r.ReadString('\n')
		for _, secret := range secrets {
			if secret != "" {
				line = strings.ReplaceAll(line, secret, "[redacted]")
			}
		}
		if line != "" {
			_, _ = io.WriteString(out, line)
		}
		if e != nil {
			return
		}
	}
}

// The userspace core owns the interface. Its unexpected exit must take down
// the manager so systemd can restart the entire tunnel with a fresh interface.
func superviseAWGForward(ctx context.Context, s Spec, stopped <-chan struct{}) error {
	forwardCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- runForward(forwardCtx, s) }()
	select {
	case e := <-done:
		return e
	case <-ctx.Done():
		cancel()
		<-done
		return nil
	case <-stopped:
		cancel()
		<-done
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("AWG userspace core stopped; restarting the tunnel is required")
	}
}
