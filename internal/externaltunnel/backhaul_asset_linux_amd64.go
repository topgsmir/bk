//go:build linux && amd64

package externaltunnel

import _ "embed"

//go:embed assets/backhaul-2.0.3-linux-amd64.tar.gz
var backhaulArchive []byte

func bundledBackhaul() ([]byte, error) { return backhaulArchive, nil }
