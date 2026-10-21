//go:build !linux || !amd64

package externaltunnel

import "fmt"

func bundledBackhaul() ([]byte, error) {
	return nil, fmt.Errorf("the supplied Backhaul core is Linux amd64 only")
}
