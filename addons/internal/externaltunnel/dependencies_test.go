package externaltunnel

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAnyConnectPreparationRepairsAnExecutableWithMissingSharedLibraries(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("native dependency checks require Linux")
	}
	previous := CoreDir
	CoreDir = t.TempDir()
	t.Cleanup(func() { CoreDir = previous })
	t.Setenv("PATH", CoreDir)
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(CoreDir, name), []byte("#!/bin/sh\n"+body), 0700); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"ip", "openconnect", "ocpasswd", "ocserv-worker"} {
		write(name, "exit 0\n")
	}
	write("ocserv", "echo 'error while loading shared libraries: libpcl.so.1: cannot open shared object file' >&2\nexit 127\n")
	if err := CheckDependencies("e3-anyconnect-tls"); err == nil || !strings.Contains(err.Error(), "libpcl.so.1") {
		t.Fatalf("broken native server passed preparation: %v", err)
	}
	installs := 0
	issues := prepareDependencies(context.Background(), []string{"e3-anyconnect-tls"}, io.Discard, CheckDependencies, func(context.Context, string, io.Writer) error {
		installs++
		write("ocserv", "[ \"$1\" = --version ] || exit 1\nexit 0\n")
		return nil
	})
	if installs != 1 || len(issues) != 0 {
		t.Fatalf("missing runtime dependencies were not repaired: installs=%d issues=%v", installs, issues)
	}
}

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
