// Package externaltunnel builds isolated, named tunnels inspired by the supplied
// GRE, L2TPv3, AWG and SSH scripts and the RGT/Paqet managers.
package externaltunnel

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/curve25519"
)

type Kind struct{ ID, Title, Source, Requirement string }

var Kinds = []Kind{
	{"gre", "GRE IPv4", "new/gre4.sh", "IPv4 protocol 47"},
	{"l2tp-ip", "L2TPv3 IP", "new/setup-l2tpv3.sh", "IPv4 protocol 115; l2tp_eth kernel module"},
	{"l2tp-udp", "L2TPv3 UDP", "new/setup-l2tpv3.sh", "UDP; l2tp_eth kernel module"},
	{"awg", "AmneziaWG", "new/awg-relay.sh", "kernel or pinned userspace AWG; UDP"},
	{"ssh", "SSH pool", "new/ssh.sh", "SSH key authentication and verified host key; TCP"},
	{"rgt-tcp", "RGT reverse TCP", "https://github.com/black-sec/RGT", "pinned RGT core; TCP"},
	{"rgt-udp", "RGT reverse UDP", "https://github.com/black-sec/RGT", "pinned RGT core; UDP data, TCP control"},
	{"rgt-direct", "RGT direct VXLAN", "https://github.com/black-sec/RGT", "VXLAN kernel support; UDP"},
	{"paqet", "Paqet KCP", "https://github.com/Ramin-Setoodehnia/Paqet-Tunnel", "pinned Paqet core; libpcap; raw TCP"},
	{"alghadir", "Alghadir full stack", "new/alghadir.sh (reconnected full stack)", "GRE/IPsec/KCP/udp2raw/obfs4; Linux network namespaces; TCP"},
}

// Spec is local state. Only paired protocol settings go into a setup/test link;
// SSH keys and this machine's physical interface never travel to the peer.
type Spec struct {
	Backend     string `json:"backend,omitempty"`
	Version     int    `json:"v"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Side        string `json:"side"`
	LocalIP     string `json:"local_ip"`
	PeerIP      string `json:"peer_ip"`
	IranIP      string `json:"iran_ip"`
	KharejIP    string `json:"kharej_ip"`
	Port        int    `json:"port"`
	SourcePort  int    `json:"source_port,omitempty"`
	ID          int    `json:"id"`
	MTU         int    `json:"mtu"`
	Secret      string `json:"secret"`
	Listen      string `json:"listen"`
	Target      string `json:"target"`
	Protocol    string `json:"protocol"`
	Connections int    `json:"connections"`
	Mode        string `json:"mode,omitempty"`
	SSHUser     string `json:"ssh_user,omitempty"`
	SSHKey      string `json:"ssh_key,omitempty"`
	KnownHosts  string `json:"known_hosts,omitempty"`
	WAN         string `json:"wan,omitempty"`
	RouterMAC   string `json:"router_mac,omitempty"`
}

var safeName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,39}$`)

func NewSecret() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func New(name, kind, side string) Spec {
	connections := 4
	if kind == "paqet" {
		connections = 1
	}
	return Spec{Version: 1, Name: name, Kind: kind, Side: side, Port: 17010, ID: 10000, MTU: 1280, Secret: NewSecret(), IranIP: "10.203.0.1/30", KharejIP: "10.203.0.2/30", Listen: "0.0.0.0:8443", Target: "127.0.0.1:8443", Protocol: "tcp", Connections: connections, Mode: "fast", SSHUser: "root", SSHKey: "/root/.ssh/bk_external", KnownHosts: "/root/.ssh/known_hosts"}
}
func Find(kind string) (Kind, bool) {
	for _, k := range Kinds {
		if k.ID == kind {
			return k, true
		}
	}
	return Kind{}, false
}
func Layer3(kind string) bool {
	switch kind {
	case "gre", "l2tp-ip", "l2tp-udp", "awg", "rgt-direct", "alghadir":
		return true
	}
	return false
}
func Reverse(kind string) bool { return strings.HasPrefix(kind, "rgt-") && kind != "rgt-direct" }
func validAddr(a string) bool {
	h, p, e := net.SplitHostPort(a)
	n, e2 := strconv.Atoi(p)
	return e == nil && e2 == nil && n > 0 && n <= 65535 && (h == "" || net.ParseIP(h) != nil)
}
func (s Spec) Validate() error {
	if s.Version != 1 || !safeName.MatchString(s.Name) {
		return fmt.Errorf("invalid version or tunnel name (letters, digits, _ and - only)")
	}
	if _, ok := Find(s.Kind); !ok {
		return fmt.Errorf("unknown external tunnel %q", s.Kind)
	}
	if s.Side != "iran" && s.Side != "kharej" {
		return fmt.Errorf("side must be iran or kharej")
	}
	for _, v := range []string{s.LocalIP, s.PeerIP} {
		if net.ParseIP(v).To4() == nil {
			return fmt.Errorf("both local and peer IPv4 addresses are required")
		}
	}
	if s.Port < 1 || s.Port > 65535 || s.SourcePort < 0 || s.SourcePort > 65535 || s.ID < 1 || s.ID > 0xffffff || s.MTU < 576 || s.MTU > 9000 {
		return fmt.Errorf("invalid tunnel port, ID or MTU")
	}
	if len(s.Secret) != 64 {
		return fmt.Errorf("paired secret must be 64 hexadecimal characters")
	}
	if _, e := hex.DecodeString(s.Secret); e != nil {
		return fmt.Errorf("invalid secret")
	}
	if s.Connections < 1 || s.Connections > 64 {
		return fmt.Errorf("connections must be 1..64")
	}
	if s.Protocol != "tcp" && s.Protocol != "udp" {
		return fmt.Errorf("protocol must be tcp or udp")
	}
	if s.Kind == "ssh" && s.Protocol != "tcp" {
		return fmt.Errorf("SSH supports TCP forwarding only")
	}
	if !validAddr(s.Listen) || !validAddr(s.Target) || (s.Backend != "" && !validAddr(s.Backend)) {
		return fmt.Errorf("listen and target must be numeric IP:port addresses")
	}
	if Layer3(s.Kind) {
		a, an, e := net.ParseCIDR(s.IranIP)
		b, bn, e2 := net.ParseCIDR(s.KharejIP)
		if e != nil || e2 != nil || a.To4() == nil || b.To4() == nil || a.Equal(b) || !an.Contains(b) || !bn.Contains(a) || an.String() != bn.String() {
			return fmt.Errorf("tunnel IPs must be different IPv4 addresses in the same subnet")
		}
	}
	if s.Kind == "paqet" {
		if s.Connections != 1 {
			return fmt.Errorf("Paqet uses one connection with its isolated fixed-port firewall rules")
		}
		switch s.Mode {
		case "normal", "fast", "fast2", "fast3":
		default:
			return fmt.Errorf("invalid Paqet KCP mode")
		}
		if s.MTU > 1500 {
			return fmt.Errorf("Paqet MTU must be <=1500")
		}
	}
	if s.Kind == "ssh" {
		if !safeName.MatchString(s.SSHUser) || !filepath.IsAbs(s.SSHKey) || !filepath.IsAbs(s.KnownHosts) {
			return fmt.Errorf("SSH needs a valid user and absolute key/known-hosts paths")
		}
	}
	if s.WAN != "" && !regexp.MustCompile(`^[a-zA-Z0-9_.:-]{1,15}$`).MatchString(s.WAN) {
		return fmt.Errorf("invalid network interface")
	}
	if s.RouterMAC != "" {
		if m, e := net.ParseMAC(s.RouterMAC); e != nil || len(m) != 6 {
			return fmt.Errorf("invalid router MAC")
		}
	}
	return nil
}
func (s Spec) Interface() string {
	h := sha256.Sum256([]byte(s.Name + "/" + s.Side))
	return "bkx" + hex.EncodeToString(h[:5])
}
func (s Spec) TunnelIP() string {
	if s.Side == "iran" {
		return s.IranIP
	}
	return s.KharejIP
}
func (s Spec) PeerTunnelIP() string {
	if s.Side == "iran" {
		return strings.Split(s.KharejIP, "/")[0]
	}
	return strings.Split(s.IranIP, "/")[0]
}
func (s Spec) Mirror() Spec {
	s.LocalIP, s.PeerIP = s.PeerIP, s.LocalIP
	if s.Side == "iran" {
		s.Side = "kharej"
	} else {
		s.Side = "iran"
	}
	s.WAN, s.RouterMAC, s.Backend = "", "", ""
	s.SSHKey, s.KnownHosts = "/root/.ssh/bk_external", "/root/.ssh/known_hosts"
	return s
}

const Scheme = "bk://e."

func (s Spec) Link() (string, error) {
	if e := s.Validate(); e != nil {
		return "", e
	}
	s.SSHKey, s.KnownHosts, s.WAN, s.RouterMAC, s.Backend = "", "", "", "", ""
	b, e := json.Marshal(s)
	return Scheme + base64.RawURLEncoding.EncodeToString(b), e
}
func ParseLink(raw string) (Spec, error) {
	var s Spec
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, Scheme) || len(raw) > 16384 {
		return s, fmt.Errorf("expected a bk://e. setup link")
	}
	b, e := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(raw, Scheme))
	if e != nil {
		return s, e
	}
	if e = json.Unmarshal(b, &s); e != nil {
		return s, e
	}
	s = s.Mirror()
	return s, s.Validate()
}
func (s Spec) Key(label string) string {
	h := sha256.Sum256([]byte("bk-external/" + label + "/" + s.Secret))
	return base64.StdEncoding.EncodeToString(h[:])
}
func (s Spec) AWGPublic(side string) string {
	k, _ := base64.StdEncoding.DecodeString(s.Key(side))
	p, e := curve25519.X25519(k, curve25519.Basepoint)
	if e != nil {
		panic(e)
	}
	return base64.StdEncoding.EncodeToString(p)
}
