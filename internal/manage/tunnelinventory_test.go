package manage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/topgsmir/bk/internal/app"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
)

func inventoryFixture(t *testing.T, directory, name, kind string) externaltunnel.Spec {
	t.Helper()
	s := externaltunnel.New(name, kind, "iran")
	s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, name+".json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestInventoryIncludesEveryPackageAndKeepsDuplicateNamesSeparate(t *testing.T) {
	dir := t.TempDir()
	inventoryFixture(t, dir, "shared", "gre")
	inventoryFixture(t, dir, "game", "d3-tcp")
	original := Tunnel{Name: "shared", Role: "server", Transport: "tcp", Service: app.ServiceName("shared")}
	got, warnings := listManagedTunnels([]Tunnel{original}, dir)
	if len(got) != 3 || len(warnings) != 0 {
		t.Fatalf("got %d tunnels, warnings %v", len(got), warnings)
	}
	for i, want := range []struct {
		pkg           int
		name, service string
	}{
		{1, "shared", app.ServiceName("shared")},
		{2, "shared", externaltunnel.ServiceName("shared")},
		{3, "game", externaltunnel.ServiceName("game")},
	} {
		if got[i].packageNumber != want.pkg || got[i].tunnel.Name != want.name || got[i].tunnel.Service != want.service {
			t.Errorf("entry %d has wrong package/name/service", i)
		}
	}
	if got[0].external != nil || got[1].external == nil || got[2].external == nil || got[1].external.Name != "shared" || got[2].external.Name != "game" {
		t.Fatal("original and optional action targets were mixed")
	}
}

func TestInventoryReportsBrokenFilesWithoutHidingWorkingTunnels(t *testing.T) {
	dir := t.TempDir()
	inventoryFixture(t, dir, "valid", "gre")
	if err := os.WriteFile(filepath.Join(dir, "broken.json"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "valid.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "wrong-name.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	got, warnings := listManagedTunnels(nil, dir)
	if len(got) != 1 || got[0].tunnel.Name != "valid" || len(warnings) != 2 {
		t.Fatalf("got %d tunnels and %v warnings", len(got), warnings)
	}
	for _, name := range []string{"broken.json", "wrong-name.json"} {
		if !strings.Contains(strings.Join(warnings, "\n"), name) {
			t.Errorf("missing warning for %s", name)
		}
	}
}

func TestInventoryWorksWhenOptionalPackagesAreNotInstalled(t *testing.T) {
	original := []Tunnel{{Name: "main"}}
	got, warnings := listManagedTunnels(original, filepath.Join(t.TempDir(), "absent"))
	if len(got) != 1 || got[0].packageNumber != 1 || len(warnings) != 0 {
		t.Fatalf("got %d tunnels, warnings %v", len(got), warnings)
	}
}

func TestUnifiedManagerStartsTheSelectedPackageService(t *testing.T) {
	dir := t.TempDir()
	inventoryFixture(t, dir, "shared", "gre")
	inventoryFixture(t, dir, "game", "d3-tcp")
	oldDir, oldList, oldActive, oldStart := externalDir, List, IsActive, StartService
	t.Cleanup(func() { externalDir, List, IsActive, StartService = oldDir, oldList, oldActive, oldStart })
	externalDir = dir
	List = func() []Tunnel {
		return []Tunnel{{Name: "shared", Role: "server", Transport: "tcp", Service: app.ServiceName("shared")}}
	}
	IsActive = func(string) bool { return false }
	for _, tc := range []struct{ choice, action, service string }{
		{"1", "2", app.ServiceName("shared")},
		{"2", "1", externaltunnel.ServiceName("shared")},
		{"3", "1", externaltunnel.ServiceName("game")},
	} {
		t.Run(tc.choice, func(t *testing.T) {
			var started []string
			StartService = func(service string) error { started = append(started, service); return nil }
			restore := tui.SetInput(strings.NewReader(tc.choice + "\n" + tc.action + "\n\n0\n0\n"))
			defer restore()
			out := captureStdout(t, ManageTunnels)
			if len(started) != 1 || started[0] != tc.service {
				t.Fatalf("started %v, want only %s", started, tc.service)
			}
			for _, want := range []string{"Package 1 |", "Package 2 |", "Package 3 |", "stopped"} {
				if !strings.Contains(out, want) {
					t.Errorf("unified list is missing %q", want)
				}
			}
		})
	}
}

func TestPackageManagersOpenTheUnifiedList(t *testing.T) {
	dir := t.TempDir()
	inventoryFixture(t, dir, "second", "gre")
	inventoryFixture(t, dir, "third", "d3-tcp")
	oldDir, oldList, oldActive := externalDir, List, IsActive
	t.Cleanup(func() { externalDir, List, IsActive = oldDir, oldList, oldActive })
	externalDir = dir
	List = func() []Tunnel { return []Tunnel{{Name: "first", Service: app.ServiceName("first")}} }
	IsActive = func(string) bool { return false }
	for _, screen := range []func(){AdditionalTunnels, Package3Tunnels} {
		restore := tui.SetInput(strings.NewReader("5\n0\n0\n"))
		out := captureStdout(t, screen)
		restore()
		for _, want := range []string{"Package 1 |", "Package 2 |", "Package 3 |"} {
			if !strings.Contains(out, want) {
				t.Errorf("package Manage entry is missing %q", want)
			}
		}
	}
}

func TestUnifiedManagerFindsOptionalTunnelsWithoutPackageOne(t *testing.T) {
	dir := t.TempDir()
	inventoryFixture(t, dir, "only-optional", "d3-tcp")
	oldDir, oldList, oldActive := externalDir, List, IsActive
	t.Cleanup(func() { externalDir, List, IsActive = oldDir, oldList, oldActive })
	externalDir = dir
	List = func() []Tunnel { return nil }
	IsActive = func(string) bool { return true }
	restore := tui.SetInput(strings.NewReader("0\n"))
	defer restore()
	out := captureStdout(t, ManageTunnels)
	if !strings.Contains(out, "only-optional") || !strings.Contains(out, "Package 3 |") || !strings.Contains(out, "running") || strings.Contains(out, "No Tunnels Yet") {
		t.Fatalf("optional-only installation is not discoverable: %s", out)
	}
}
