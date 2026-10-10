package externaltunnel

import (
	"strings"
	"testing"
)

func TestIPsecCleanupSelectsOnlyReservedReqid(t *testing.T) {
	body := `src 10.0.0.1/32 dst 10.0.0.2/32 proto udp sport 1701 dport 1701
 dir out priority 399999 index 10
 tmpl src 10.0.0.1 dst 10.0.0.2 proto esp reqid 10000 mode transport
src 10.0.0.1/32 dst 10.0.0.3/32
 dir in priority 399999 index 17
 tmpl src 10.0.0.1 dst 10.0.0.3 proto esp reqid 10100 mode transport
src 10.0.0.2/32 dst 10.0.0.1/32
 dir in priority 399999 index 18
 tmpl src 10.0.0.2 dst 10.0.0.1 proto esp reqid 10000 mode transport
`
	got := ownedIPsecPolicies(body, 10000)
	if len(got) != 2 || got[0].index != "10" || got[1].index != "18" {
		t.Fatalf("wrong owners: %+v", got)
	}
	if len(ownedIPsecPolicies(body, 50000)) != 0 {
		t.Fatal("unrelated policy selected")
	}
	if len(ownedIPsecPolicies(strings.ReplaceAll(body, "index 10", "index ;reboot"), 10000)) != 1 {
		t.Fatal("unsafe index selected")
	}
}
func TestL2TPIPsecRequiresExactVPNPair(t *testing.T) {
	s := New("pair", "e3-l2tp-ipsec", "iran")
	s.LocalIP, s.PeerIP = "192.0.2.10", "192.0.2.11"
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.IranIP = "10.203.0.3/30"
	if s.Validate() == nil {
		t.Fatal("invalid PPP server address accepted")
	}
	s = New("pair", "e3-l2tp-ipsec", "iran")
	s.LocalIP, s.PeerIP = "192.0.2.10", "192.0.2.11"
	s.Port = 20000
	if s.Validate() == nil {
		t.Fatal("misleading arbitrary L2TP port accepted")
	}
}
