package externaltunnel

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func SingBox(kind string) bool { return strings.HasPrefix(kind, "e3-sing-") }
func singBoxAsset(arch string) (Asset, error) {
	switch arch {
	case "386":
		return Asset{URL: "https://github.com/SagerNet/sing-box/releases/download/v1.14.3/sing-box-1.14.3-linux-386.tar.gz", SHA256: "92dad4663615d2cb30fb8feba48787e8529b1940004341712b0d0d4548cb42ff", Member: "sing-box"}, nil
	case "amd64":
		return Asset{URL: "https://github.com/SagerNet/sing-box/releases/download/v1.14.3/sing-box-1.14.3-linux-amd64.tar.gz", SHA256: "e2bdf179c15a3652955dc44867e33ca0965c9c9b91b34267dcc5a5d639a5feee", Member: "sing-box"}, nil
	case "arm64":
		return Asset{URL: "https://github.com/SagerNet/sing-box/releases/download/v1.14.3/sing-box-1.14.3-linux-arm64.tar.gz", SHA256: "29dae24d76bea3bc62d07dd3e60cfc12003de7debbd6b768df9be9c68360174e", Member: "sing-box"}, nil
	case "arm":
		return Asset{URL: "https://github.com/SagerNet/sing-box/releases/download/v1.14.3/sing-box-1.14.3-linux-armv6.tar.gz", SHA256: "63d91abdde9b9de2d1ef9d607d17779e4d9b4c23bef0ffa06bbb8bacb82311f2", Member: "sing-box"}, nil
	case "mipsle":
		return Asset{URL: "https://github.com/SagerNet/sing-box/releases/download/v1.14.3/sing-box-1.14.3-linux-mipsle.tar.gz", SHA256: "49e0f4ac3ee9f8c1cbba97419247d270517c5b2a54624c51d7e68d547ee42ab9", Member: "sing-box"}, nil
	}
	return Asset{}, fmt.Errorf("no pinned sing-box 1.14.3 release for %s", arch)
}
func installSingBox(ctx context.Context, out io.Writer) error {
	a, err := singBoxAsset(runtime.GOARCH)
	if err != nil {
		return err
	}
	path := filepath.Join(CoreDir, "sing-box")
	if err = installAsset(ctx, a, path, out); err != nil {
		return err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(body)
	return os.WriteFile(filepath.Join(CoreDir, "sing-box.sha256"), []byte(a.SHA256+" "+hex.EncodeToString(sum[:])+"\n"), 0600)
}
func checkSingBoxCore() error {
	a, err := singBoxAsset(runtime.GOARCH)
	if err != nil {
		return err
	}
	manifest, err := os.ReadFile(filepath.Join(CoreDir, "sing-box.sha256"))
	if err != nil {
		return err
	}
	fields := strings.Fields(string(manifest))
	if len(fields) != 2 || fields[0] != a.SHA256 {
		return fmt.Errorf("sing-box pinned version metadata differs")
	}
	body, err := os.ReadFile(filepath.Join(CoreDir, "sing-box"))
	if err != nil {
		return err
	}
	sum := sha256.Sum256(body)
	if hex.EncodeToString(sum[:]) != fields[1] {
		return fmt.Errorf("sing-box installed core checksum differs")
	}
	return nil
}
func singBoxConfig(s Spec, dir string) (map[string]any, error) {
	family := strings.TrimPrefix(s.Kind, "e3-sing-")
	tlsSpec := s
	seed, _ := base64.StdEncoding.DecodeString(s.Key("eylan/sing-box/tls-root"))
	tlsSpec.Secret = hex.EncodeToString(seed)
	if family != "shadowsocks" {
		if _, err := backhaulTLS(tlsSpec, dir); err != nil {
			return nil, err
		}
	}
	inTLS := map[string]any{"enabled": true, "certificate_path": filepath.Join(dir, "tls.pem"), "key_path": filepath.Join(dir, "tls-key.pem"), "min_version": "1.3"}
	outTLS := map[string]any{"enabled": true, "server_name": "bk-backhaul.local", "certificate_path": filepath.Join(dir, "tls.pem"), "min_version": "1.3"}
	raw, _ := base64.StdEncoding.DecodeString(s.Key("eylan/sing-box/uuid"))
	raw[6] = (raw[6] & 15) | 64
	raw[8] = (raw[8] & 63) | 128
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", raw[:4], raw[4:6], raw[6:8], raw[8:10], raw[10:16])
	carrier := map[string]any{"type": family, "tag": "carrier"}
	if s.Side == "iran" {
		carrier["server"], carrier["server_port"] = s.PeerIP, s.Port
	} else {
		carrier["listen"], carrier["listen_port"] = s.LocalIP, s.Port
	}
	switch family {
	case "vless", "vmess":
		if s.Side == "iran" {
			carrier["uuid"] = uuid
			if family == "vmess" {
				carrier["security"] = "auto"
			}
		} else {
			carrier["users"] = []any{map[string]any{"name": "bk", "uuid": uuid}}
		}
	case "trojan", "hysteria2":
		if s.Side == "iran" {
			carrier["password"] = s.Secret
		} else {
			carrier["users"] = []any{map[string]any{"name": "bk", "password": s.Secret}}
		}
	case "shadowsocks":
		carrier["method"], carrier["password"] = "2022-blake3-aes-256-gcm", s.Key("eylan/sing-box/shadowsocks")
	default:
		return nil, fmt.Errorf("unknown sing-box carrier %s", family)
	}
	if family != "shadowsocks" {
		if s.Side == "iran" {
			carrier["tls"] = outTLS
		} else {
			carrier["tls"] = inTLS
		}
	}
	if family == "hysteria2" {
		carrier["up_mbps"], carrier["down_mbps"] = 100, 100
	}
	host, portText, err := net.SplitHostPort(s.Target)
	if err != nil {
		return nil, err
	}
	port, _ := strconv.Atoi(portText)
	cfg := map[string]any{"log": map[string]any{"level": "info"}}
	if s.Side == "iran" {
		bind, listenPortText, _ := net.SplitHostPort(s.Listen)
		listenPort, _ := strconv.Atoi(listenPortText)
		cfg["inbounds"] = []any{map[string]any{"type": "direct", "tag": "application", "listen": bind, "listen_port": listenPort, "network": s.Protocol, "override_address": host, "override_port": port}}
		cfg["outbounds"] = []any{carrier}
		cfg["route"] = map[string]any{"final": "carrier"}
	} else {
		cfg["inbounds"] = []any{carrier}
		cfg["outbounds"] = []any{map[string]any{"type": "direct", "tag": "backend"}}
		bits := 32
		if net.ParseIP(host).To4() == nil {
			bits = 128
		}
		cfg["route"] = map[string]any{"rules": []any{map[string]any{"ip_cidr": []string{fmt.Sprintf("%s/%d", host, bits)}, "port": []int{port}, "action": "route", "outbound": "backend"}, map[string]any{"action": "reject"}}}
	}
	return cfg, nil
}
func runSingBox(ctx context.Context, s Spec, dir string) error {
	cfg, err := singBoxConfig(s, dir)
	if err != nil {
		return err
	}
	body, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "sing-box.json")
	if err = os.WriteFile(path, body, 0600); err != nil {
		return err
	}
	core := filepath.Join(CoreDir, "sing-box")
	if err = runCommand(ctx, command(core, "check", "-c", path)); err != nil {
		return fmt.Errorf("sing-box configuration: %w", err)
	}
	return childRedacted(ctx, core, []string{s.Secret, s.Key("eylan/sing-box/shadowsocks")}, "run", "-c", path)
}
