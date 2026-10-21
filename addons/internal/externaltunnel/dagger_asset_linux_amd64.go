//go:build linux && amd64

package externaltunnel

import _ "embed"

//go:embed assets/dagger-rs-linux-x86_64.tar.gz
var daggerArchive []byte

func bundledDagger() ([]byte, error) { return daggerArchive, nil }
