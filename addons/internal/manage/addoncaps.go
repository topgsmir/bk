package manage

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

type Tunnel struct {
	Name, Role, Transport, Addr, Service string
	Ports                                []string
}

const (
	PresetTurbo      = "turbo"
	PresetAggressive = "aggressive"
)

var PublicIPv4 = func() string {
	addresses, _ := net.InterfaceAddrs()
	for _, address := range addresses {
		ip, _, _ := net.ParseCIDR(address.String())
		if ip.To4() != nil && !ip.IsLoopback() {
			return ip.String()
		}
	}
	return "-"
}

func serviceAction(action, service string) error {
	if !strings.HasPrefix(service, "bk-ext-") {
		return fmt.Errorf("not an addon service: %s", service)
	}
	out, err := exec.Command("systemctl", action, service).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}
	return nil
}

var StartService = func(service string) error { return serviceAction("start", service) }
var StopService = func(service string) error { return serviceAction("stop", service) }
var RestartService = func(service string) error { return serviceAction("restart", service) }
var IsActive = func(service string) bool {
	out, err := exec.Command("systemctl", "is-active", service).Output()
	return err == nil && strings.TrimSpace(string(out)) == "active"
}

func FollowLog(service string) {
	command := exec.Command("journalctl", "-fu", service)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	_ = command.Run()
}
