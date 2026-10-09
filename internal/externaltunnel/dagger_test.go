package externaltunnel

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestDaggerCatalogAndPairedSettings(t *testing.T) {
	if len(daggerKinds()) != 44 {
		t.Fatal(len(daggerKinds()))
	}
	seen := map[string]bool{}
	for _, k := range daggerKinds() {
		if seen[k.ID] {
			t.Fatal(k.ID)
		}
		seen[k.ID] = true
		s := New("dagger-test", k.ID, "iran")
		s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
		s.LocalIPv6, s.PeerIPv6 = "fd42:3::1", "fd42:3::2"
		if e := s.Validate(); e != nil {
			t.Fatal(k.ID, e)
		}
		if !BothProtocols(k.ID) {
			t.Fatal(k.ID)
		}
		if got := s.Mirror().Mirror(); got.LocalIPv6 != s.LocalIPv6 || got.PeerIPv6 != s.PeerIPv6 {
			t.Fatal(k.ID, "IPv6 pairing")
		}
		if daggerPublic(s, "iran") == daggerPublic(s, "kharej") {
			t.Fatal("identities not separated")
		}
	}
}
func TestEveryDaggerConfigurationAcceptedBySuppliedCore(t *testing.T) {
	core := os.Getenv("BK_DAGGER_BINARY")
	if core == "" {
		t.Skip("requires provided real Dagger core; mandatory in package3 CI")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, kind := range daggerKinds() {
		for _, side := range []string{"iran", "kharej"} {
			for _, proto := range []string{"tcp", "udp"} {
				s := New("verify", kind.ID, side)
				s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
				s.LocalIPv6, s.PeerIPv6 = "fd42:3::1", "fd42:3::2"
				s.Protocol = proto
				dir := t.TempDir()
				cfg, e := daggerConfig(s, dir)
				if e != nil {
					t.Fatal(kind.ID, side, e)
				}
				b, _ := json.Marshal(cfg)
				path := filepath.Join(dir, "core.json")
				if e = os.WriteFile(path, b, 0600); e != nil {
					t.Fatal(e)
				}
				out, e := exec.CommandContext(ctx, core, "--config", path, "--check").CombinedOutput()
				if e != nil {
					t.Fatalf("%s %s %s: %v: %s", kind.ID, side, proto, e, out)
				}
			}
		}
	}
}

func TestPinnedDaggerNoticesIncludeCompleteRustCopyright(t *testing.T) {
	archive, err := bundledDagger()
	if err != nil {
		t.Skip("supplied amd64 asset only")
	}
	previous := CoreDir
	CoreDir = t.TempDir()
	defer func() { CoreDir = previous }()
	if err := installDaggerNotices(archive); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(CoreDir, "dagger-notices", "licenses", "third-party", "rust-stdlib-1.99.0", "COPYRIGHT-library.html"))
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != 1499465 {
		t.Fatalf("copyright notice truncated: %d", len(body))
	}
}
