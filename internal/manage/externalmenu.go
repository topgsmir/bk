package manage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/topgsmir/bk/internal/app"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
)

var externalDir = app.ConfigDir + "/external"

func ConnectionTestMenu() {
	switch tui.ChooseOpt("Connection Test", []tui.Option{
		{Title: "Original bk transports", Desc: "the existing full reverse/direct test"},
		{Title: "Additional tunnels", Desc: "GRE, L2TPv3, AWG, SSH, RGT, Paqet, Alghadir"},
	}) {
	case 0:
		ConnectionTest()
	case 1:
		ExternalConnectionTest()
	}
}
func AdditionalTunnels() {
	for {
		tui.Clear()
		tui.Title("Additional Tunnels")
		pick := tui.ChooseOpt("Action", []tui.Option{{Title: "Connection test", Desc: "temporary real tunnels, same echo/soak/bulk test"}, {Title: "Setup Iran", Desc: "create an additional tunnel"}, {Title: "Setup Kharej", Desc: "create the other end"}, {Title: "Apply setup link", Desc: "paired settings from the other server"}, {Title: "Manage", Desc: "start, stop, logs, delete"}, {Title: "Install dependencies", Desc: "pinned cores and required packages"}, {Title: "Prepare SSH authentication", Desc: "dedicated key; verify the host fingerprint"}})
		switch pick {
		case 0:
			ExternalConnectionTest()
		case 1:
			setupExternal("iran")
		case 2:
			setupExternal("kharej")
		case 3:
			raw := tui.Prompt("bk://e. Link: ")
			s, e := externaltunnel.ParseLink(raw)
			if e != nil {
				tui.Error(e.Error())
				tui.PressEnter()
				continue
			}
			finishExternalSetup(s)
		case 4:
			manageExternal()
		case 5:
			k := chooseExternalKind()
			if k == "" {
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			e := externaltunnel.InstallDependencies(ctx, k, os.Stdout)
			cancel()
			report(e, "Dependencies Installed")
		case 6:
			prepareExternalSSH()
		default:
			return
		}
	}
}
func chooseExternalKind() string {
	opts := make([]tui.Option, len(externaltunnel.Kinds))
	for i, k := range externaltunnel.Kinds {
		opts[i] = tui.Option{Title: k.Title, Desc: k.Requirement}
	}
	i := tui.ChooseOpt("Tunnel Type", opts)
	if i < 0 {
		return ""
	}
	return externaltunnel.Kinds[i].ID
}
func externalNumber(label string, n int) int {
	v, e := strconv.Atoi(tui.PromptDefault(label, strconv.Itoa(n)))
	if e != nil {
		return -1
	}
	return v
}
func setupExternal(side string) {
	k := chooseExternalKind()
	if k == "" {
		return
	}
	s := externaltunnel.New(strings.TrimSpace(tui.Prompt("Tunnel Name: ")), k, side)
	s.LocalIP = strings.TrimSpace(tui.PromptDefault("This Server's Local IPv4 (Assigned To Its NIC)", linkHost()))
	s.PeerIP = strings.TrimSpace(tui.Prompt("Other Server's IPv4: "))
	if k != "ssh" {
		s.Port = externalNumber("Tunnel Port", s.Port)
	}
	if k == "ssh" {
		s.Port = externalNumber("SSH Port", 22)
	}
	if externaltunnel.Layer3(k) {
		s.IranIP = tui.PromptDefault("Iran Tunnel IPv4/CIDR", s.IranIP)
		s.KharejIP = tui.PromptDefault("Kharej Tunnel IPv4/CIDR", s.KharejIP)
		s.ID = externalNumber("Paired Tunnel ID (Unique)", s.ID)
		s.MTU = externalNumber("MTU", s.MTU)
	}
	if k != "ssh" {
		s.Secret = tui.PromptDefault("Paired Secret (64 Hex Characters, Same On Both Servers)", s.Secret)
	}
	s.Listen = tui.PromptDefault("Iran Listen IP:Port", s.Listen)
	s.Target = tui.PromptDefault("Kharej Target IP:Port", s.Target)
	if k == "rgt-udp" {
		s.Protocol = "udp"
	} else if k != "ssh" && k != "rgt-tcp" {
		if i := tui.ChooseOpt("Forwarded Traffic", []tui.Option{{Title: "TCP"}, {Title: "UDP (Also Keeps TCP Available For Layer-3)"}}); i == 1 {
			s.Protocol = "udp"
		} else if i < 0 {
			return
		}
	}
	if k == "ssh" {
		s.Connections = externalNumber("Parallel Connections (1..64)", s.Connections)
	}
	if k == "ssh" && side == "iran" {
		s.SSHUser = tui.PromptDefault("Kharej SSH User", s.SSHUser)
		s.SSHKey = tui.PromptDefault("SSH Private Key Path", s.SSHKey)
		s.KnownHosts = tui.PromptDefault("Verified Known Hosts File", s.KnownHosts)
	}
	if k == "paqet" {
		s.Mode = tui.PromptDefault("KCP Mode (normal/fast/fast2/fast3)", s.Mode)
		s.MTU = externalNumber("Paqet MTU (<=1500)", s.MTU)
		s.WAN = tui.PromptDefault("Physical Interface (Blank = Detect)", "")
		s.RouterMAC = tui.PromptDefault("Next-Hop MAC (Blank = Detect)", "")
	}
	finishExternalSetup(s)
}
func finishExternalSetup(s externaltunnel.Spec) {
	if e := externaltunnel.Check(s); e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	cfg := filepath.Join(externalDir, s.Name+".json")
	unit := filepath.Join(app.ServiceDir, externaltunnel.ServiceName(s.Name))
	if _, e := os.Lstat(unit); e == nil {
		tui.Error("That additional tunnel already exists.")
		tui.PressEnter()
		return
	}
	// Apply links may contain a peer-specific name. Validation rejects traversal.
	if e := externaltunnel.Save(cfg, s); e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	if e := os.WriteFile(unit, []byte(externaltunnel.Unit(s, cfg, app.BinPath)), 0644); e != nil {
		os.Remove(cfg)
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	run := func(args ...string) error {
		b, e := exec.CommandContext(ctx, "systemctl", args...).CombinedOutput()
		if e != nil {
			return fmt.Errorf("%w: %s", e, b)
		}
		return nil
	}
	if e := run("daemon-reload"); e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	if e := run("enable", "--now", externaltunnel.ServiceName(s.Name)); e != nil {
		tui.Error("Saved, but startup failed: " + e.Error())
		tui.PressEnter()
		return
	}
	time.Sleep(500 * time.Millisecond)
	if e := run("is-active", "--quiet", externaltunnel.ServiceName(s.Name)); e != nil {
		tui.Error("The process stopped. Use Additional Tunnels → Manage → Log.")
	} else {
		tui.Success("Service Started. Run Option 0 On Both Servers To Verify Real Traffic.")
	}
	if link, e := s.Link(); e == nil {
		tui.Info("Setup Link For The Other Server (Contains The Paired Secret):")
		fmt.Println(link)
	}
	tui.PressEnter()
}
func manageExternal() {
	for {
		files, e := filepath.Glob(filepath.Join(externalDir, "*.json"))
		if e != nil || len(files) == 0 {
			tui.Warn("No Additional Tunnels Yet.")
			tui.PressEnter()
			return
		}
		var specs []externaltunnel.Spec
		var opts []tui.Option
		for _, f := range files {
			if s, e := externaltunnel.Load(f); e == nil {
				specs = append(specs, s)
				opts = append(opts, tui.Option{Title: s.Name, Desc: s.Kind + " " + s.Side + " " + plainState(externaltunnel.ServiceName(s.Name))})
			}
		}
		i := tui.ChooseOpt("Additional Tunnel", opts)
		if i < 0 {
			return
		}
		s := specs[i]
		service := externaltunnel.ServiceName(s.Name)
		a := tui.ChooseOpt("Action", []tui.Option{{Title: "Start"}, {Title: "Stop"}, {Title: "Restart"}, {Title: "Log"}, {Title: "Setup Link"}, {Title: "Delete"}})
		switch a {
		case 0:
			report(StartService(service), "Started")
		case 1:
			report(StopService(service), "Stopped")
		case 2:
			report(RestartService(service), "Restarted")
		case 3:
			FollowLog(service)
		case 4:
			link, e := s.Link()
			if e == nil {
				fmt.Println(link)
			} else {
				tui.Error(e.Error())
			}
			tui.PressEnter()
		case 5:
			if !tui.Confirm("Delete Additional Tunnel "+s.Name, false) {
				continue
			}
			if e := StopService(service); e != nil {
				report(e, "Stopped")
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			e := exec.CommandContext(ctx, "systemctl", "disable", service).Run()
			cancel()
			if e != nil {
				report(e, "Disabled")
				continue
			}
			os.Remove(filepath.Join(app.ServiceDir, service))
			os.Remove(filepath.Join(externalDir, s.Name+".json"))
			_ = exec.Command("systemctl", "daemon-reload").Run()
			tui.Success("Deleted Only This Additional Tunnel.")
			tui.PressEnter()
		default:
			return
		}
	}
}
func prepareExternalSSH() {
	host := strings.TrimSpace(tui.Prompt("Kharej IPv4: "))
	port := externalNumber("SSH Port", 22)
	user := tui.PromptDefault("SSH User", "root")
	s := externaltunnel.New("ssh-auth", "ssh", "iran")
	s.LocalIP = "127.0.0.1"
	s.PeerIP = host
	s.Port = port
	s.SSHUser = user
	if e := s.Validate(); e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	os.MkdirAll(filepath.Dir(s.SSHKey), 0700)
	if _, e := os.Lstat(s.SSHKey); os.IsNotExist(e) {
		cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-N", "", "-f", s.SSHKey)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if e = cmd.Run(); e != nil {
			report(e, "Key Created")
			return
		}
	}
	tui.Info("Verify The SSH Host Fingerprint Against The Kharej Server Before Accepting It.")
	// ssh-copy-id owns the interactive password/host-key prompts; never trust
	// ssh-keyscan output automatically or disable host-key verification.
	cmd := exec.Command("ssh-copy-id", "-p", strconv.Itoa(port), "-i", s.SSHKey+".pub", user+"@"+host)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	report(cmd.Run(), "SSH Key Prepared")
}
