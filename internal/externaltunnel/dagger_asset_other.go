//go:build !linux || !amd64

package externaltunnel

import "fmt"

func bundledDagger() ([]byte, error) {
	return nil, fmt.Errorf("the supplied Dagger release requires Linux amd64")
}
