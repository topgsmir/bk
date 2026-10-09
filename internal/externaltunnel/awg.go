package externaltunnel

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed scripts/install-awg.sh
var awgInstaller string

func awgTool() string {
	if _, e := exec.LookPath("awg"); e == nil {
		return "awg"
	}
	return filepath.Join(CoreDir, "awg")
}
func installAWG(ctx context.Context) error {
	if e := runCommand(ctx, command("apt-get", "install", "-y", "build-essential", "curl", "ca-certificates")); e != nil {
		return e
	}
	cmd := exec.CommandContext(ctx, "bash", "-c", awgInstaller)
	cmd.Env = append(os.Environ(), "BK_EXTRA_CORES="+CoreDir)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
func startAWGUserspace(ctx context.Context, iface string) (func(), <-chan struct{}, error) {
	cmd := exec.Command(filepath.Join(CoreDir, "amneziawg-go"), "-f", iface)
	cmd.Env = append(os.Environ(), "WG_I_PREFER_BUGGY_USERSPACE_TO_POLISHED_KMOD=1")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if e := cmd.Start(); e != nil {
		return nil, nil, e
	}
	done := make(chan struct{})
	var result error
	go func() { result = cmd.Wait(); close(done) }()
	stop := func() {
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-done:
		case <-time.After(4 * time.Second):
			_ = cmd.Process.Kill()
			<-done
		}
	}
	for i := 0; i < 50; i++ {
		if _, e := os.Stat("/var/run/amneziawg/" + iface + ".sock"); e == nil {
			return stop, done, nil
		}
		if _, e := os.Stat("/var/run/wireguard/" + iface + ".sock"); e == nil {
			return stop, done, nil
		}
		select {
		case <-done:
			return nil, nil, fmt.Errorf("AWG userspace startup: %v", result)
		case <-ctx.Done():
			stop()
			return nil, nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	stop()
	return nil, nil, fmt.Errorf("AWG userspace interface did not become ready")
}
