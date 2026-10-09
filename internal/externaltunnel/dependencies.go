package externaltunnel

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func CheckDependencies(kind string) error {
	if _, ok := Find(kind); !ok {
		return fmt.Errorf("unknown additional tunnel %q", kind)
	}
	if runtime.GOOS != "linux" {
		return fmt.Errorf("additional tunnels require Linux")
	}
	if _, e := CoreAsset(kind, runtime.GOARCH); e != nil {
		return e
	}
	programs := []string{}
	if Layer3(kind) {
		programs = append(programs, "ip")
	}
	switch kind {
	case "l2tp-ip", "l2tp-udp":
		programs = append(programs, "modprobe")
	case "awg":
		programs = append(programs, awgTool(), filepath.Join(CoreDir, "amneziawg-go"))
	case "paqet":
		programs = append(programs, "ip", "iptables", filepath.Join(CoreDir, "paqet"))
	case "rgt-tcp", "rgt-udp":
		programs = append(programs, filepath.Join(CoreDir, "rgt"))
	case "alghadir":
		programs = append(programs, "ip", "iptables", "obfs4proxy", filepath.Join(CoreDir, "udp2raw"))
	}
	for _, program := range programs {
		if _, e := exec.LookPath(program); e != nil {
			return fmt.Errorf("missing %s", program)
		}
	}
	if kind == "l2tp-ip" || kind == "l2tp-udp" {
		modules := []string{"l2tp_netlink", "l2tp_eth"}
		if kind == "l2tp-ip" {
			modules = append(modules, "l2tp_ip")
		}
		for _, module := range modules {
			if _, e := os.Stat("/sys/module/" + module); e == nil {
				continue
			}
			if e := exec.Command("modprobe", "-n", module).Run(); e != nil {
				return fmt.Errorf("kernel module %s is unavailable", module)
			}
		}
	}
	return nil
}

// Every selected kind receives a preparation result. Shared cores are installed
// once, and every checksum/install error is retained for both test boards.
func PrepareDependencies(ctx context.Context, kinds []string, out io.Writer) map[string]string {
	return prepareDependencies(ctx, kinds, out, CheckDependencies, InstallDependencies)
}
func prepareDependencies(ctx context.Context, kinds []string, out io.Writer, check func(string) error, install func(context.Context, string, io.Writer) error) map[string]string {
	result := map[string]string{}
	attempted := map[string]error{}
	for _, kind := range kinds {
		if e := check(kind); e == nil {
			continue
		}
		key := kind
		if kind == "rgt-udp" {
			key = "rgt-tcp"
		}
		if kind == "ssh-reverse" {
			key = "ssh"
		}
		previous, seen := attempted[key]
		if !seen {
			fmt.Fprintf(out, "Preparing %s dependencies before the connection test...\n", kind)
			step, cancel := context.WithTimeout(ctx, 10*time.Minute)
			previous = install(step, kind, out)
			cancel()
			attempted[key] = previous
		}
		if previous != nil {
			result[kind] = "dependency preparation: " + previous.Error()
			continue
		}
		if e := check(kind); e != nil {
			result[kind] = "dependency preparation did not make the method ready: " + e.Error()
		}
	}
	return result
}
