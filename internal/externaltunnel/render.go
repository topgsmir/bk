package externaltunnel

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// Command is executed directly, never through a shell or sourced input.
type Command struct {
	Program string
	Args    []string
}

func command(program string, args ...string) Command { return Command{program, args} }
func (s Spec) KernelCommands() (up, down []Command) {
	iface := s.Interface()
	id := strconv.Itoa(s.ID)
	port := strconv.Itoa(s.Port)
	switch s.Kind {
	case "gre":
		up = append(up, command("ip", "tunnel", "add", iface, "mode", "gre", "local", s.LocalIP, "remote", s.PeerIP, "key", id, "ttl", "64"))
		down = append(down, command("ip", "tunnel", "del", iface))
	case "l2tp-ip", "l2tp-udp":
		args := []string{"l2tp", "add", "tunnel", "tunnel_id", id, "peer_tunnel_id", id, "encap", "ip", "local", s.LocalIP, "remote", s.PeerIP}
		if s.Kind == "l2tp-udp" {
			args[8] = "udp"
			args = append(args, "udp_sport", port, "udp_dport", port)
		}
		up = append(up, command("ip", args...), command("ip", "l2tp", "add", "session", "tunnel_id", id, "session_id", id, "peer_session_id", id, "name", iface))
		down = append(down, command("ip", "l2tp", "del", "session", "tunnel_id", id, "session_id", id), command("ip", "l2tp", "del", "tunnel", "tunnel_id", id))
	case "rgt-direct":
		up = append(up, command("ip", "link", "add", iface, "type", "vxlan", "id", id, "local", s.LocalIP, "remote", s.PeerIP, "dstport", port, "nolearning"))
		down = append(down, command("ip", "link", "del", iface))
	case "awg":
		up = append(up, command("ip", "link", "add", iface, "type", "amneziawg"))
		down = append(down, command("ip", "link", "del", iface))
	}
	if len(up) > 0 {
		up = append(up, command("ip", "addr", "add", s.TunnelIP(), "dev", iface), command("ip", "link", "set", "dev", iface, "mtu", strconv.Itoa(s.MTU), "up"))
	}
	return
}
func (s Spec) AWGConfig() string {
	peer := "iran"
	endpoint := ""
	if s.Side == "iran" {
		peer = "kharej"
		endpoint = fmt.Sprintf("Endpoint = %s:%d\n", s.PeerIP, s.Port)
	}
	return fmt.Sprintf("[Interface]\nPrivateKey = %s\nListenPort = %d\nJc = 4\nJmin = 40\nJmax = 80\nS1 = 23\nS2 = 47\nH1 = 100000001\nH2 = 100000002\nH3 = 100000003\nH4 = 100000004\n\n[Peer]\nPublicKey = %s\nPresharedKey = %s\nAllowedIPs = %s/32\n%sPersistentKeepalive = 25\n", s.Key(s.Side), s.Port, s.AWGPublic(peer), s.Key("psk"), s.PeerTunnelIP(), endpoint)
}
func (s Spec) RGTConfig() string {
	tr := strings.TrimPrefix(s.Kind, "rgt-")
	role, field, addr, heartbeat := "server", "bind_addr", net.JoinHostPort(s.LocalIP, strconv.Itoa(s.Port)), "heartbeat_interval"
	serviceField, serviceAddr := "bind_addr", s.Listen
	if s.Side == "kharej" {
		role, field, addr, heartbeat = "client", "remote_addr", net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)), "heartbeat_timeout"
		serviceField, serviceAddr = "local_addr", s.Target
	}
	return fmt.Sprintf("[%s]\n%s = %q\ndefault_token = %q\n%s = 30\n[%s.transport]\ntype = %q\n[%s.transport.%s]\nnodelay = true\nkeepalive_secs = 20\nkeepalive_interval = 8\n[%s.services.bkforward]\ntype = %q\ntoken = %q\n%s = %q\nnodelay = true\n", role, field, addr, s.Secret, heartbeat, role, "tcp", role, "tcp", role, tr, s.Secret, serviceField, serviceAddr)
}
func (s Spec) PaqetConfig() string {
	role, port := "client", s.SourcePort
	if s.Side == "kharej" {
		role, port = "server", s.Port
	}
	extra := fmt.Sprintf("server:\n  addr: %q\nforward:\n  - listen: %q\n    target: %q\n    protocol: %q\n", net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)), s.Listen, s.Target, s.Protocol)
	if role == "server" {
		extra = fmt.Sprintf("listen:\n  addr: %q\n", net.JoinHostPort(s.LocalIP, strconv.Itoa(s.Port)))
	}
	return fmt.Sprintf("role: %q\nlog:\n  level: \"info\"\n%snetwork:\n  interface: %q\n  ipv4:\n    addr: %q\n    router_mac: %q\n  tcp:\n    local_flag: [\"PA\"]\n    remote_flag: [\"PA\"]\ntransport:\n  protocol: \"kcp\"\n  conn: %d\n  kcp:\n    mode: %q\n    mtu: %d\n    block: \"aes\"\n    key: %q\n", role, extra, s.WAN, net.JoinHostPort(s.LocalIP, strconv.Itoa(port)), s.RouterMAC, s.Connections, s.Mode, s.MTU, s.Secret)
}
func ServiceName(name string) string { return "bk-ext-" + name + ".service" }
func Unit(s Spec, config, bin string) string {
	return fmt.Sprintf("[Unit]\nDescription=bk additional tunnel %s (%s)\nAfter=network-online.target\nWants=network-online.target\n\n[Service]\nType=simple\nUMask=0077\nExecStart=%s external run -c %s\nRestart=on-failure\nRestartSec=5\nTimeoutStopSec=25\nKillMode=control-group\nLimitNOFILE=65536\n\n[Install]\nWantedBy=multi-user.target\n", s.Name, s.Kind, strconv.Quote(filepath.Clean(bin)), strconv.Quote(filepath.Clean(config)))
}
