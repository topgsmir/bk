//go:build !linux || !amd64

package externaltunnel

import "fmt"

func bundledSolarpass() ([]byte, error) {
	return nil, fmt.Errorf("the supplied Solarpass core is Linux amd64 only")
}
