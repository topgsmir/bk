//go:build linux && amd64

package externaltunnel

import _ "embed"

//go:embed assets/solarpass-2.3.0-linux-amd64.tar.gz
var solarpassArchive []byte

func bundledSolarpass() ([]byte, error) { return solarpassArchive, nil }
