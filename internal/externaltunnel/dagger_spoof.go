package externaltunnel

import (
	"fmt"
	"net"
	"strings"
)

func DaggerSpoof(kind string) bool { return Dagger(kind) && strings.HasSuffix(kind, "-spoof") }

func (s Spec) validateDaggerSpoof() error {
	values := []string{s.IranSpoofSource, s.KharejSpoofSource, s.IranSpoofDestination, s.KharejSpoofDestination}
	count := 0
	for _, value := range values {
		if value == "" {
			continue
		}
		count++
		ip := net.ParseIP(value)
		if ip == nil || ip.To4() == nil || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLoopback() || value == "255.255.255.255" {
			return fmt.Errorf("Dagger spoof addresses must be usable literal unicast IPv4 addresses")
		}
	}
	if DaggerSpoof(s.Kind) {
		if count == 0 {
			return fmt.Errorf("IP Spoof needs at least one explicit outer source or destination address")
		}
	} else if count != 0 {
		return fmt.Errorf("spoof addresses require a Dagger IP Spoof method; DCPI and spoofing cannot be combined")
	}
	return nil
}
