package manage

import (
	"os"
	"strings"
	"testing"

	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
)

func TestKharejSetupStartsWithLinkInBothManagers(t *testing.T) {
	for _, catalog := range [][]externaltunnel.Kind{externaltunnel.Kinds, externaltunnel.Package3Kinds} {
		func() {
			restore := tui.SetInput(strings.NewReader("\n"))
			defer restore()
			out, err := os.CreateTemp(t.TempDir(), "menu-")
			if err != nil {
				t.Fatal(err)
			}
			oldStdout := os.Stdout
			os.Stdout = out
			defer func() { os.Stdout = oldStdout; out.Close() }()
			setupExternalCatalog("kharej", catalog)
			body, err := os.ReadFile(out.Name())
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(body), "bk://e. Link:") || strings.Contains(string(body), "Tunnel Type") || strings.Contains(string(body), "Tunnel Name:") {
				t.Fatalf("kharej opened the manual wizard instead of asking for the Iran link: %s", body)
			}
		}()
	}
}

func TestPackage3SetupLinkPreservesGamingICMPPair(t *testing.T) {
	iran := externaltunnel.New("gaming", "d3-gaming-icmp", "iran")
	iran.LocalIP, iran.PeerIP = "192.0.2.1", "192.0.2.2"
	iran.Listen, iran.Target = "0.0.0.0:2053", "127.0.0.1:2053"
	iran.SourcePort = 17011
	iran.WAN, iran.RouterMAC = "eth-iran", "02:00:00:00:00:01"
	link, err := externalSetupLink(iran)
	if err != nil {
		t.Fatal(err)
	}
	kharej, err := parseExternalSetupLink(link, externaltunnel.Package3Kinds)
	if err != nil {
		t.Fatal(err)
	}
	if kharej.Side != "kharej" || kharej.LocalIP != iran.PeerIP || kharej.PeerIP != iran.LocalIP || kharej.Secret != iran.Secret {
		t.Fatal("import did not retain paired identity and exchange endpoint roles")
	}
	if kharej.Listen != iran.Listen || kharej.Target != iran.Target || kharej.Port != iran.Port || kharej.SourcePort != iran.SourcePort || kharej.Kind != iran.Kind {
		t.Fatal("import changed the forwarding or raw packet settings")
	}
	if kharej.WAN != "" || kharej.RouterMAC != "" {
		t.Fatal("Iran-specific network interface leaked into kharej settings")
	}
	if link, err := externalSetupLink(kharej); err == nil || link != "" {
		t.Fatal("kharej must not issue a second setup link")
	}
	wrongDirection, err := kharej.Link()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parseExternalSetupLink(wrongDirection, externaltunnel.Package3Kinds); err == nil {
		t.Fatal("accepted a kharej-origin link as an Iran setup link")
	}
	if _, err := parseExternalSetupLink(link, externaltunnel.Kinds); err == nil {
		t.Fatal("accepted a Package 3 link in another package")
	}
}

func TestAdditionalSetupLinksPreserveEveryPair(t *testing.T) {
	for _, kind := range externaltunnel.Kinds {
		t.Run(kind.ID, func(t *testing.T) {
			iran := externaltunnel.New("additional", kind.ID, "iran")
			iran.LocalIP, iran.PeerIP = "192.0.2.1", "192.0.2.2"
			iran.ID, iran.MTU = 12001, 1300
			iran.IranIP, iran.KharejIP = "10.204.0.1/30", "10.204.0.2/30"
			iran.Listen, iran.Target = "0.0.0.0:2053", "127.0.0.1:2053"
			iran.WAN, iran.RouterMAC = "eth-iran", "02:00:00:00:00:01"
			iran.SSHKey, iran.KnownHosts = "/iran/private-key", "/iran/known-hosts"
			link, err := externalSetupLink(iran)
			if err != nil {
				t.Fatal(err)
			}
			peer, err := parseExternalSetupLink(link, externaltunnel.Kinds)
			if err != nil {
				t.Fatal(err)
			}
			if peer.Side != "kharej" || peer.LocalIP != iran.PeerIP || peer.PeerIP != iran.LocalIP || peer.Secret != iran.Secret || peer.Kind != iran.Kind || peer.Name != iran.Name {
				t.Fatal("import did not preserve the paired identity and reverse endpoint roles")
			}
			if peer.ID != iran.ID || peer.MTU != iran.MTU || peer.IranIP != iran.IranIP || peer.KharejIP != iran.KharejIP || peer.Port != iran.Port || peer.SourcePort != iran.SourcePort || peer.Listen != iran.Listen || peer.Target != iran.Target || peer.Protocol != iran.Protocol || peer.Connections != iran.Connections {
				t.Fatal("import changed shared tunnel or forwarding settings")
			}
			if peer.WAN != "" || peer.RouterMAC != "" || peer.SSHKey == iran.SSHKey || peer.KnownHosts == iran.KnownHosts {
				t.Fatal("import copied Iran-specific interface or SSH files")
			}
			if link, err := externalSetupLink(peer); err == nil || link != "" {
				t.Fatal("kharej issued a return setup link")
			}
			wrongDirection, err := peer.Link()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := parseExternalSetupLink(wrongDirection, externaltunnel.Kinds); err == nil {
				t.Fatal("accepted a kharej-origin link")
			}
			if _, err := parseExternalSetupLink(link, externaltunnel.Package3Kinds); err == nil {
				t.Fatal("accepted an additional link in Package 3")
			}
		})
	}
}
