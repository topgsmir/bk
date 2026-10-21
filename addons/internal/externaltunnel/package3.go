package externaltunnel

import "strings"

// Package 3 has its own catalog; the existing additional methods keep their catalog.
var Package3Kinds = append(append(append(daggerKinds(), solarpassKinds()...), backhaulKinds()...), eylanKinds()...)

func daggerKinds() []Kind {
	var kinds []Kind
	add := func(id, title string) {
		kinds = append(kinds, Kind{id, title, "new2/dagger: dagger-rs 0.2.1, MIT, ir_spoof", "pinned Dagger core; authenticated TCP and UDP forwarding"})
	}
	for _, c := range []string{"tcp", "kcp", "http", "https", "ws", "wss", "dc6", "qplus"} {
		add("d3-"+c, "Dagger "+strings.ToUpper(c))
	}
	for _, c := range []string{"xhttp", "xhttps"} {
		for _, mode := range []string{"stream", "packet", "auto"} {
			for _, reverse := range []bool{false, true} {
				id := "d3-" + c + "-" + mode
				title := "Dagger " + strings.ToUpper(c) + " " + mode
				if reverse {
					id += "-reverse"
					title += " reversed direction"
				}
				add(id, title)
			}
		}
	}
	for _, c := range []string{"quantum", "gaming", "tun"} {
		for _, profile := range []string{"tcp", "udp", "icmp", "gre", "ipip", "bip", "raw", "dcpi"} {
			add("d3-"+c+"-"+profile, "Dagger "+strings.ToUpper(c)+" / "+strings.ToUpper(profile))
		}
	}

	for _, carrier := range []string{"quantum", "gaming", "tun"} {
		for _, profile := range []string{"tcp", "udp", "icmp", "gre", "ipip", "bip", "raw"} {
			add("d3-"+carrier+"-"+profile+"-spoof", "Dagger "+strings.ToUpper(carrier)+" / "+strings.ToUpper(profile)+" / IP Spoof")
		}
	}
	return kinds
}
func Dagger(kind string) bool { return strings.HasPrefix(kind, "d3-") }
func DaggerOptions(kind string) (carrier, profile, xhttp string, reverse bool) {
	parts := strings.Split(strings.TrimPrefix(kind, "d3-"), "-")
	carrier = parts[0]
	if carrier == "qplus" {
		carrier = "quantum+"
	}
	if carrier == "gaming" {
		carrier = "quantum-gaming"
	}
	if carrier == "xhttp" || carrier == "xhttps" {
		xhttp = parts[1]
		if xhttp != "auto" {
			xhttp += "-up"
		}
		reverse = len(parts) > 2
	}
	if carrier == "quantum" || carrier == "quantum-gaming" || carrier == "tun" {
		profile = parts[1]
	}
	return
}

// Packet maps on these cores independently carry TCP and UDP.
func BothProtocols(kind string) bool {
	return Layer3(kind) || kind == "paqet" || Dagger(kind) || Eylan(kind) || kind == "b3-tcp" || (Solarpass(kind) && !UDPOnly(kind))
}

func Solarpass(kind string) bool { return strings.HasPrefix(kind, "s3-") }
func UDPOnly(kind string) bool {
	return kind == "rgt-udp" || kind == "s3-spoof" || kind == "s3-gamepass"
}
func solarpassKinds() []Kind {
	var result []Kind
	for _, transport := range []string{"tcp", "tcpmux", "ws", "wss", "wsmux", "wssmux", "quic", "quicmux", "hysteria", "spoof", "gamepass", "ultimatepass", "speedpass", "ultimategamepass", "speedpassmux", "arenapass"} {
		result = append(result, Kind{"s3-" + transport, "Solarpass " + strings.ToUpper(transport), "new2/solarpass-tun: supplied 2.3.0 core", "provided amd64 core; private kernel tuning namespace"})
	}
	return result
}

func Backhaul(kind string) bool { return strings.HasPrefix(kind, "b3-") }
func backhaulKinds() []Kind {
	var result []Kind
	for _, carrier := range []string{"tcp", "tcpmux", "xtcpmux", "ws", "wss", "wsmux", "wssmux", "xwsmux", "anytls", "tun-tcp", "tun-ipx-icmp", "tun-ipx-ipip", "tun-ipx-udp", "tun-ipx-tcp", "tun-ipx-gre", "tun-ipx-bip"} {
		result = append(result, Kind{"b3-" + carrier, "Backhaul " + strings.ToUpper(carrier), "new2/backhaul: supplied 2.0.3 core", "assigned public IPv4; pinned amd64 core; paired authenticated forwarding"})
	}
	return result
}

func Eylan(kind string) bool { return strings.HasPrefix(kind, "e3-") }
func eylanKinds() []Kind {
	result := []Kind{
		{"e3-l2tp-ipsec", "Eylan / L2TP IPsec", "new2/EylanPanel; stock strongSwan/xl2tpd", "paired IKEv2/AES-GCM; unused UDP 500/4500/1701; PPP kernel"},
		{"e3-anyconnect-tls", "Eylan / AnyConnect TLS", "new2/EylanPanel; stock ocserv/OpenConnect", "verified certificate and login; TLS tunnel"},
		{"e3-anyconnect-dtls", "Eylan / AnyConnect DTLS", "new2/EylanPanel; stock ocserv/OpenConnect", "verified certificate and login; DTLS handshake required"},
		{"e3-wireguard", "Eylan / WireGuard", "new2/EylanPanel: standard WireGuard protocol", "included WireGuard userspace engine; UDP"},
		{"e3-openvpn-tcp", "Eylan / OpenVPN TCP", "new2/EylanPanel: standard OpenVPN protocol", "OpenVPN 2.6+; paired TLS; TCP"},
		{"e3-openvpn-udp", "Eylan / OpenVPN UDP", "new2/EylanPanel: standard OpenVPN protocol", "OpenVPN 2.6+; paired TLS; UDP"},
	}
	for _, family := range []string{"vless", "vmess", "trojan", "shadowsocks", "hysteria2"} {
		result = append(result, Kind{"e3-sing-" + family, "Eylan / " + strings.ToUpper(family), "new2/EylanPanel; stock sing-box by SagerNet", "pinned sing-box; paired authentication; TCP and UDP"})
	}
	return result
}
