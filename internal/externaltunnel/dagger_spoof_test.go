package externaltunnel

import (
	"strings"
	"testing"
)

func TestDaggerSpoofSettingsSurvivePairingAndSelectCorrectSide(t *testing.T) {
	s := New("spoof-pair", "d3-quantum-tcp-spoof", "iran")
	s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
	s.IranSpoofSource, s.KharejSpoofSource = "198.18.42.1", "198.18.42.2"
	s.IranSpoofDestination, s.KharejSpoofDestination = "198.19.42.2", "198.19.42.1"
	link, err := s.Link()
	if err != nil {
		t.Fatal(err)
	}
	peer, err := ParseLink(link)
	if err != nil {
		t.Fatal(err)
	}
	if peer.Side != "kharej" || peer.IranSpoofSource != s.IranSpoofSource || peer.KharejSpoofDestination != s.KharejSpoofDestination {
		t.Fatal("paired spoof fields changed")
	}
	for _, side := range []Spec{s, peer} {
		cfg, err := daggerConfig(side, t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		field := "listeners"
		src, dst := s.IranSpoofSource, s.IranSpoofDestination
		if side.Side == "kharej" {
			field = "paths"
			src, dst = s.KharejSpoofSource, s.KharejSpoofDestination
		}
		raw := cfg[field].([]any)[0].(map[string]any)["raw"].(map[string]any)
		if raw["spoof_src_ip"] != src || raw["spoof_dst_ip"] != dst || raw["peer_ip"] != side.PeerIP || raw["local_ip"] != side.LocalIP {
			t.Fatal("spoof settings replaced routing endpoints or wrong peer fields")
		}
	}
	for _, mode := range []string{"source", "destination"} {
		test := s
		if mode == "source" {
			test.IranSpoofDestination = ""
			test.KharejSpoofDestination = ""
		} else {
			test.IranSpoofSource = ""
			test.KharejSpoofSource = ""
		}
		if err := test.Validate(); err != nil {
			t.Fatal(mode, err)
		}
	}
}
func TestDaggerSpoofRejectsUnavailableAndInvalidCombinations(t *testing.T) {
	s := New("spoof-pair", "d3-tun-udp-spoof", "iran")
	s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
	if err := s.Validate(); err == nil {
		t.Fatal("empty spoof mode accepted")
	}
	for _, ip := range []string{"not-an-ip", "::1", "0.0.0.0", "127.0.0.1", "224.0.0.1", "255.255.255.255"} {
		s.IranSpoofSource = ip
		if err := s.Validate(); err == nil {
			t.Fatal(ip)
		}
	}
	s.IranSpoofSource = "198.18.42.1"
	for _, kind := range []string{"d3-tun-dcpi", "d3-quantum-dcpi", "d3-qplus", "d3-tcp", "b3-tcp"} {
		s.Kind = kind
		if err := s.Validate(); err == nil || !strings.Contains(err.Error(), "spoof") {
			t.Fatalf("%s accepted spoofing: %v", kind, err)
		}
	}
}
