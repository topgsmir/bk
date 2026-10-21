package manage

import (
	"github.com/topgsmir/bk/addons/internal/externaltunnel"
	"github.com/topgsmir/bk/addons/internal/tui"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSeparatedManagerListsBothAddonPackagesAndStartsOnlyTheSelectedService(t *testing.T) {
	directory := t.TempDir()
	for _, item := range []struct{ name, kind string }{{"second", "gre"}, {"third", "d3-tcp"}} {
		spec := externaltunnel.New(item.name, item.kind, "iran")
		spec.LocalIP = "192.0.2.1"
		spec.PeerIP = "192.0.2.2"
		if err := externaltunnel.Save(filepath.Join(directory, item.name+".json"), spec); err != nil {
			t.Fatal(err)
		}
	}
	oldDir, oldActive, oldStart := externalDir, IsActive, StartService
	t.Cleanup(func() { externalDir, IsActive, StartService = oldDir, oldActive, oldStart })
	externalDir = directory
	IsActive = func(string) bool { return false }
	var started []string
	StartService = func(service string) error { started = append(started, service); return nil }
	restore := tui.SetInput(strings.NewReader("2\n1\n\n0\n0\n"))
	defer restore()
	out := captureAddonStdout(t, ManageTunnels)
	if len(started) != 1 || started[0] != externaltunnel.ServiceName("third") {
		t.Fatalf("wrong addon service: %v", started)
	}
	for _, want := range []string{"Package 2 |", "Package 3 |", "stopped"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q: %s", want, out)
		}
	}
	if strings.Contains(out, "Package 1 |") {
		t.Fatal("addon manager claimed ownership of core tunnels")
	}
}

func captureAddonStdout(t *testing.T, action func()) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = original }()
	result := make(chan string, 1)
	go func() { raw, _ := io.ReadAll(reader); reader.Close(); result <- string(raw) }()
	action()
	writer.Close()
	return <-result
}

func TestAddonServiceActionsRejectCoreServicesBeforeRunningSystemctl(t *testing.T) {
	if err := serviceAction("start", "bk-core.service"); err == nil {
		t.Fatal("addon service control accepted a core service")
	}
}
