package externaltunnel

import (
	"net"
	"strconv"
	"testing"
)

func TestSolarpassPairedNativeConfiguration(t *testing.T) {
	for _, kind := range solarpassKinds() {
		for _, protocol := range []string{"tcp", "udp"} {
			if UDPOnly(kind.ID) && protocol == "tcp" {
				continue
			}
			s := New("solar-case", kind.ID, "iran")
			s.LocalIP, s.PeerIP = "192.0.2.10", "192.0.2.11"
			s.Protocol = protocol
			if err := s.Validate(); err != nil {
				t.Fatal(kind.ID, protocol, err)
			}
			for _, spec := range []Spec{s, s.Mirror()} {
				cfg, err := solarpassConfig(spec)
				if err != nil {
					t.Fatal(err)
				}
				if cfg["token"] != s.Secret {
					t.Fatal("peer token differs")
				}
				if rules := cfg["ports"].([]any); len(rules) > 0 {
					rule := rules[0].(map[string]any)
					target := s.Target
					if protocol == "udp" {
						target = net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort+3))
					}
					if rule["target"] != target {
						t.Fatal(kind.ID, "wrong backend/decoder", rule)
					}
				}
			}
		}
	}
	s := New("blocked", "s3-spoof", "iran")
	s.LocalIP, s.PeerIP = "192.0.2.10", "192.0.2.11"
	s.SourcePort = s.Port
	if err := s.Validate(); err == nil {
		t.Fatal("overlapping native port blocks accepted")
	}
	s.SourcePort = 65533
	if err := s.Validate(); err == nil {
		t.Fatal("overflowing native port block accepted")
	}
	s.SourcePort = 17020
	s.Protocol = "tcp"
	if err := s.Validate(); err == nil {
		t.Fatal("unsupported native Spoof TCP accepted")
	}
}
