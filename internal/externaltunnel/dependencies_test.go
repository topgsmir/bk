package externaltunnel

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestMissingSelectedCoresArePreparedAndSharedRGTInstalledOnce(t *testing.T) {
	ready := map[string]bool{}
	calls := map[string]int{}
	check := func(kind string) error {
		if kind == "rgt-udp" {
			kind = "rgt-tcp"
		}
		if ready[kind] {
			return nil
		}
		return fmt.Errorf("missing core")
	}
	install := func(ctx context.Context, kind string, out io.Writer) error {
		calls[kind]++
		ready[kind] = true
		return nil
	}
	issues := prepareDependencies(context.Background(), []string{"paqet", "rgt-tcp", "rgt-udp", "paqet"}, io.Discard, check, install)
	if len(issues) != 0 || calls["paqet"] != 1 || calls["rgt-tcp"] != 1 || calls["rgt-udp"] != 0 {
		t.Fatal(issues, calls)
	}
}
func TestFailedCoreInstallationRemainsAResultForEverySelectedMethod(t *testing.T) {
	calls := 0
	check := func(string) error { return fmt.Errorf("missing") }
	install := func(context.Context, string, io.Writer) error {
		calls++
		return fmt.Errorf("download checksum mismatch")
	}
	issues := prepareDependencies(context.Background(), []string{"rgt-tcp", "rgt-udp"}, io.Discard, check, install)
	if calls != 1 || len(issues) != 2 || !strings.Contains(issues["rgt-udp"], "checksum mismatch") {
		t.Fatal(calls, issues)
	}
}
