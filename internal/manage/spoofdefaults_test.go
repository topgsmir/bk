package manage

import (
	"reflect"
	"strings"
	"testing"

	"github.com/topgsmir/bk/config"
	"github.com/topgsmir/bk/internal/tui"
)

func TestSpoofUDPWizardRotationIsIndependentOfStealth(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		rotation    bool
	}{
		{"accept defaults", "1\n\n\n\n", true},
		{"explicit off", "1\n\n\nn\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			restore := tui.SetInput(strings.NewReader(tc.input))
			defer restore()
			var sc config.SpoofConfig
			askSpoofShared(&sc)
			if sc.SpoofProfile != "udp" || sc.SpoofShufflePort != tc.rotation || spoofStealthOn(sc) {
				t.Fatalf("wizard did not keep the independent rotation answer: %+v", sc)
			}
		})
	}
}

// Exercise the same render -> link -> kharej creation path as the wizard.
func TestSpoofUDPDefaultsSurvivePairedSetup(t *testing.T) {
	iran := l3Spec{
		Name: "spoof-iran", Side: sideIran, Carrier: "spoof", Encap: "gre",
		Addr: "203.0.113.20:9000", Token: "regression-test-token",
		Iface: "bp0", LocalIP: "10.70.0.1/30", PeerIP: "10.70.0.2",
		Ports: []string{"443"}, MTU: defaultL3MTU,
		Spoof: newSpoofConfig("udp", "", ""),
	}
	iran.Spoof.SpoofSrcIP, iran.Spoof.SpoofPeerIP = "198.51.100.10", "203.0.113.20"
	findL3Preset("").apply(&iran)
	applyNewL3CarrierDefaults(&iran)
	for _, rotation := range []bool{true, false} {
		t.Run(map[bool]string{true: "default rotation", false: "explicit rotation off"}[rotation], func(t *testing.T) {
			iran.Spoof.SpoofShufflePort = rotation
			if !rotation {
				iran.Spoof.SpoofPortMin, iran.Spoof.SpoofPortMax = 0, 0
			}
			link, err := DecodeShareLink(pendingShareLinkFrom(iran, "203.0.113.10", linkExtras{}))
			if err != nil {
				t.Fatal(err)
			}
			form := MirrorForPeer(link)
			form.Spoof.SrcIPs = "198.51.100.20"
			kharej, err := form.ToNewDirectTunnel().spec()
			if err != nil {
				t.Fatal(err)
			}
			for _, cfg := range []config.L3Config{decode(t, iran.Render()).L3, decode(t, kharej.Render()).L3} {
				if cfg.MTU != 1280 || cfg.AutoMTUEnabled() || cfg.SpoofShufflePort != rotation || cfg.SpoofSockBuf != 8<<20 {
					t.Fatalf("paired UDP settings lost: %+v", cfg)
				}
				if rotation && (cfg.SpoofPortMin != 49152 || cfg.SpoofPortMax != 65535) {
					t.Fatal("source-port range lost")
				}
			}
			if kharej.Spoof.SpoofPeerSrcIP != "198.51.100.10" || kharej.Spoof.SpoofPeerIP != "203.0.113.10" {
				t.Fatal("producer's forged and real IPs were not mirrored correctly")
			}
		})
	}
}

func TestSpoofDefaultsRespectImportedMTUAndOtherCarriers(t *testing.T) {
	on := true
	n := NewDirectTunnel{Side: "kharej", Carrier: "spoof", Name: "imported",
		Token: "regression-test-token", TunnelPort: "9000", SpoofPeerIP: "203.0.113.10", MTU: 1360, AutoMTU: &on}
	spec, err := n.spec()
	if err != nil {
		t.Fatal(err)
	}
	if spec.MTU != 1360 || !*spec.AutoMTU {
		t.Fatal("imported MTU choice was overridden")
	}
	for _, carrier := range []string{"udp", "pck", "quic", "sni", "xdi"} {
		s := l3Spec{Carrier: carrier, MTU: defaultL3MTU}
		before := s
		applyNewL3CarrierDefaults(&s)
		if !reflect.DeepEqual(s, before) {
			t.Fatalf("defaults affected %s", carrier)
		}
	}
	for _, profile := range []string{"tcp", "icmp"} {
		s := l3Spec{Carrier: "spoof", MTU: defaultL3MTU, Spoof: newSpoofConfig(profile, "", "")}
		before := s
		applyNewL3CarrierDefaults(&s)
		if !reflect.DeepEqual(s, before) {
			t.Fatalf("UDP defaults affected %s", profile)
		}
	}
}

func TestSpoofOldLinksKeepTheProducersImplicitAutoMTU(t *testing.T) {
	link := ShareLink{Kind: "direct", From: "iran", Tr: "spoof", Profile: "udp",
		Name: "old-link", Tok: "regression-test-token", Port: "9000", Host: "203.0.113.10",
		LocalIP: "10.71.0.1/30", PeerIP: "10.71.0.2", MTU: 1400}
	form := MirrorForPeer(link)
	spec, err := form.ToNewDirectTunnel().spec()
	if err != nil {
		t.Fatal(err)
	}
	if spec.MTU != 1400 || spec.AutoMTU == nil || !*spec.AutoMTU {
		t.Fatal("old link's implicit Auto MTU was replaced with the new setup default")
	}
	if !spec.Spoof.SpoofShufflePort {
		t.Fatal("old UDP link did not receive independent local port rotation")
	}
}

func TestSpoofEditKeepsRotationAndHiddenSettingsWhenStealthIsOff(t *testing.T) {
	off := false
	l := config.L3Config{Carrier: "spoof", Mode: "dial", MTU: 1280, AutoMTU: &off,
		SpoofConfig: newSpoofConfig("udp", "", "")}
	l.SpoofSrcIP, l.SpoofPeerIP, l.SpoofPeerSrcIP = "198.51.100.10", "203.0.113.20", "198.51.100.20"
	l.SpoofTTLJitter, l.SpoofMTU = true, 1500
	settings := directSettingsFrom("test", l)
	settings.Spoof.SrcIPs = "198.51.100.11"
	got, err := applyDirectEdit(l, DirectEdit{Spoof: &settings.Spoof, Stealth: &off})
	if err != nil {
		t.Fatal(err)
	}
	want := l
	want.SpoofSrcIP = "198.51.100.11"
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IP edit reset other settings:\ngot %+v\nwant %+v", got, want)
	}
	applySpoofStealth(&got.SpoofConfig)
	clearSpoofStealth(&got.SpoofConfig)
	if !got.SpoofShufflePort || got.SpoofPortMin != 49152 || got.SpoofPortMax != 65535 {
		t.Fatal("Stealth toggle disabled independent rotation")
	}
}
