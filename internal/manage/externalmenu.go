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
		{Title: "Test Tunnel Package 3", Desc: "Dagger and supplied package 3 methods"},
	}) {
	case 0:
		ConnectionTest()
	case 1:
		ExternalConnectionTest()
	case 2:
		Package3ConnectionTest()
	}
}
func AdditionalTunnels() { additionalTunnels("Additional Tunnels", externaltunnel.Kinds) }
func Package3Tunnels()   { additionalTunnels("Build Tunnel Package 3", externaltunnel.Package3Kinds) }
func additionalTunnels(title string, catalog []externaltunnel.Kind) {
	for {
		tui.Clear()
		tui.Title(title)
		actions := []tui.Option{{Title: "Connection test", Desc: "temporary real tunnels, same echo/soak/bulk test"}, {Title: "Setup Iran", Desc: "create an additional tunnel"}, {Title: "Setup Kharej", Desc: "create the other end"}, {Title: "Apply setup link", Desc: "paired settings from the other server"}, {Title: "Manage", Desc: "start, stop, logs, delete"}, {Title: "Install dependencies", Desc: "pinned cores and required packages"}, {Title: "Prepare SSH authentication", Desc: "dedicated key; verify the host fingerprint"}}
		if len(catalog) > 0 && externaltunnel.Eylan(catalog[len(catalog)-1].ID) {
			actions = actions[:6]
		}
		pick := tui.ChooseOpt("Action", actions)
		switch pick {
		case 0:
			externalConnectionTestCatalog(title, catalog)
		case 1:
			setupExternalCatalog("iran", catalog)
		case 2:
			setupExternalCatalog("kharej", catalog)
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
			manageExternalCatalog(catalog)
		case 5:
			k := chooseExternalKindCatalog(catalog)
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
func chooseExternalKind() string { return chooseExternalKindCatalog(externaltunnel.Kinds) }
func chooseExternalKindCatalog(catalog []externaltunnel.Kind) string {
	if package3Catalog(catalog) {
		var ok bool
		catalog, ok = choosePackage3Family(catalog)
		if !ok {
			return ""
		}
	}
	opts := make([]tui.Option, len(catalog))
	for i, k := range catalog {
		opts[i] = tui.Option{Title: k.Title, Desc: k.Requirement}
	}
	i := tui.ChooseOpt("Tunnel Type", opts)
	if i < 0 {
		return ""
	}
	return catalog[i].ID
}
func externalNumber(label string, n int) int {
	v, e := strconv.Atoi(tui.PromptDefault(label, strconv.Itoa(n)))
	if e != nil {
		return -1
	}
	return v
}
func setupExternal(side string) { setupExternalCatalog(side, externaltunnel.Kinds) }
func setupExternalCatalog(side string, catalog []externaltunnel.Kind) {
	k := chooseExternalKindCatalog(catalog)
	if k == "" {
		return
	}
	s := externaltunnel.New(strings.TrimSpace(tui.Prompt("Tunnel Name: ")), k, side)
	s.LocalIP = strings.TrimSpace(tui.PromptDefault("This Server's Local IPv4 (Assigned To Its NIC)", linkHost()))
	s.PeerIP = strings.TrimSpace(tui.Prompt("Other Server's IPv4: "))
	if externaltunnel.L2TPIPsec(k) {
		tui.Info("L2TP/IPsec Uses UDP 500, 4500 And 1701. Existing VPN Services Must Leave These Ports Free.")
	} else if !externaltunnel.SSH(k) {
		s.Port = externalNumber("Tunnel Port", s.Port)
	}
	if externaltunnel.SSH(k) {
		s.Port = externalNumber("SSH Port", 22)
	}
	if externaltunnel.Layer3(k) {
		s.IranIP = tui.PromptDefault("Iran Tunnel IPv4/CIDR", s.IranIP)
		s.KharejIP = tui.PromptDefault("Kharej Tunnel IPv4/CIDR", s.KharejIP)
		s.ID = externalNumber("Paired Tunnel ID (Unique)", s.ID)
		s.MTU = externalNumber("MTU", s.MTU)
	}
	if !externaltunnel.SSH(k) {
		s.Secret = tui.PromptDefault("Paired Secret (64 Hex Characters, Same On Both Servers)", s.Secret)
	}
	s.Listen = tui.PromptDefault("Iran Listen IP:Port", s.Listen)
	s.Target = tui.PromptDefault("Kharej Target IP:Port", s.Target)
	if externaltunnel.UDPOnly(k) {
		s.Protocol = "udp"
	} else if externaltunnel.BothProtocols(k) {
		if i := tui.ChooseOpt("Forwarded Traffic", []tui.Option{{Title: "TCP"}, {Title: "UDP (Also Keeps TCP Available For Layer-3)"}}); i == 1 {
			s.Protocol = "udp"
		} else if i < 0 {
			return
		}
	}
	if k == "ssh-reverse" {
		s.SourcePort = externalNumber("Unused Iran Loopback Forward Port (>=1024)", 17011)
	}
	if k == "ssh" {
		s.Connections = externalNumber("Parallel Connections (1..64)", s.Connections)
	}
	if externaltunnel.SSH(k) {
		s.SSHUser = tui.PromptDefault("SSH Server User (Kharej For Direct, Iran For Reverse)", s.SSHUser)
	}
	if externaltunnel.SSHInitiator(s) {
		s.SSHKey = tui.PromptDefault("SSH Private Key Path", externaltunnel.DefaultSSHKey())
		s.KnownHosts = tui.PromptDefault("Verified Known Hosts File", s.KnownHosts)
	}
	if k == "paqet" {
		s.Mode = tui.PromptDefault("KCP Mode (normal/fast/fast2/fast3)", s.Mode)
		s.MTU = externalNumber("Paqet MTU (<=1500)", s.MTU)
		s.WAN = tui.PromptDefault("Physical Interface (Blank = Detect)", "")
		s.RouterMAC = tui.PromptDefault("Next-Hop MAC (Blank = Detect)", "")
	}
	if externaltunnel.Solarpass(k) || externaltunnel.Backhaul(k) || externaltunnel.Eylan(k) {
		s.SourcePort = externalNumber("Unused Frontend/Return/Sync/Decoder Port Block (Four Ports)", s.SourcePort)
	}
	if externaltunnel.Dagger(k) {
		if k == "d3-dc6" {
			s.LocalIPv6 = strings.TrimSpace(tui.Prompt("This Server IPv6 Assigned To Its NIC: "))
			s.PeerIPv6 = strings.TrimSpace(tui.Prompt("Other Server IPv6: "))
		}
		_, profile, _, _ := externaltunnel.DaggerOptions(k)
		if profile != "" {
			s.WAN = tui.PromptDefault("Physical Interface (Blank = Route Detect)", "")
			s.RouterMAC = tui.PromptDefault("Next-Hop MAC (Blank = ARP Detect)", "")
		}
		if strings.HasPrefix(k, "d3-tun-") {
			s.SourcePort = externalNumber("Inner Forwarding Port (>=1024, Same On Both Ends)", s.SourcePort)
		}
	}
	finishExternalSetup(s)
}
func finishExternalSetup(s externaltunnel.Spec) {
	if e := s.Validate(); e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	ctxPrepare, cancelPrepare := context.WithTimeout(context.Background(), 10*time.Minute)
	issues := externaltunnel.PrepareDependencies(ctxPrepare, []string{s.Kind}, os.Stdout)
	if why := issues[s.Kind]; why != "" {
		cancelPrepare()
		tui.Error(why)
		tui.PressEnter()
		return
	}
	if e := ensureExternalSSH(ctxPrepare, s, os.Stdout); e != nil {
		cancelPrepare()
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	cancelPrepare()

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
		tui.Error("The process stopped. Use Additional Tunnels â†’ Manage â†’ Log.")
	} else {
		tui.Success("Service Started. Run Option 0 On Both Servers To Verify Real Traffic.")
	}
	if link, e := s.Link(); e == nil {
		tui.Info("Setup Link For The Other Server (Contains The Paired Secret):")
		fmt.Println(link)
	}
	tui.PressEnter()
}
func manageExternal() { manageExternalCatalog(externaltunnel.Kinds) }
func manageExternalCatalog(catalog []externaltunnel.Kind) {
	allowed := map[string]bool{}
	for _, k := range catalog {
		allowed[k.ID] = true
	}
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
			if s, e := externaltunnel.Load(f); e == nil && allowed[s.Kind] {
				specs = append(specs, s)
				opts = append(opts, tui.Option{Title: s.Name, Desc: s.Kind + " " + s.Side + " " + plainState(externaltunnel.ServiceName(s.Name))})
			}
		}
		if len(specs) == 0 {
			tui.Warn("No Tunnels In This Package Yet.")
			tui.PressEnter()
			return
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
	choice := tui.ChooseOpt("SSH Connection Direction", []tui.Option{{Title: "Direct: Iran Connects To Kharej"}, {Title: "Reverse: Kharej Connects To Iran"}})
	if choice < 0 {
		return
	}
	kind, side := "ssh", "iran"
	if choice == 1 {
		kind, side = "ssh-reverse", "kharej"
	}
	peer := strings.TrimSpace(tui.Prompt("SSH Server IPv4: "))
	s := promptExternalSSH(kind, side, "127.0.0.1", peer)
	if kind == "ssh-reverse" {
		s.SourcePort = 17011
	}
	if e := s.Validate(); e != nil {
		tui.Error(e.Error())
		tui.PressEnter()
		return
	}
	ctx, cancel := connTestContext()
	defer cancel()
	report(ensureExternalSSH(ctx, s, os.Stdout), "SSH Authentication Prepared")
}

func package3Catalog(catalog []externaltunnel.Kind) bool {
	return len(catalog) > 0 && (externaltunnel.Dagger(catalog[0].ID) || externaltunnel.Solarpass(catalog[0].ID) || externaltunnel.Backhaul(catalog[0].ID) || externaltunnel.Eylan(catalog[0].ID))
}
func package3Families(catalog []externaltunnel.Kind) [][]externaltunnel.Kind {
	groups := make([][]externaltunnel.Kind, 4)
	for _, kind := range catalog {
		index := 3
		switch {
		case externaltunnel.Dagger(kind.ID):
			index = 0
		case externaltunnel.Solarpass(kind.ID):
			index = 1
		case externaltunnel.Backhaul(kind.ID):
			index = 2
		}
		groups[index] = append(groups[index], kind)
	}
	return groups
}
func choosePackage3Family(catalog []externaltunnel.Kind) ([]externaltunnel.Kind, bool) {
	groups := package3Families(catalog)
	index := tui.ChooseOpt("Package 3 Family", []tui.Option{{Title: "Dagger", Desc: "44 transports and profiles"}, {Title: "Solarpass", Desc: "16 carriers including repaired Spoof"}, {Title: "Backhaul", Desc: "stream and TUN transports"}, {Title: "Eylan VPN methods", Desc: "WireGuard, OpenVPN, AnyConnect, L2TP/IPsec, sing-box"}})
	if index < 0 {
		return nil, false
	}
	return groups[index], true
}
