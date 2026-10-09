package externaltunnel

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Asset struct{ URL, SHA256, Member string }

func CoreAsset(kind, arch string) (Asset, error) {
	switch kind {
	case "rgt-tcp", "rgt-udp":
		if arch != "amd64" {
			return Asset{}, fmt.Errorf("upstream RGT publishes only amd64; this machine is %s", arch)
		}
		return Asset{"https://raw.githubusercontent.com/black-sec/RGT/d71d5e18a7e783f80c2961db86fc4fe84bdb8abf/core/RGT-x86-64-linux.zip", "3033358eb769ae72948a302fed7ca458369f55befc570182bbb6e1e17fb7aa34", "rgt"}, nil
	case "paqet":
		hashes := map[string]string{"amd64": "4f2f69e0493746726598485ef18df43372470a1514f8021594137e522f0f4fc7", "arm64": "be5c7b35a0a93832063e1dc7fddeb76d4597fb9de76abdfdf7b8574b775e3885", "arm": "ad7c4868a1dd1db9ae77f32ecc5859df0c7081a58745c40b754d25526586bf70"}
		h := hashes[arch]
		if h == "" {
			return Asset{}, fmt.Errorf("no pinned Paqet core for %s", arch)
		}
		if arch == "arm" {
			arch = "arm32"
		}
		return Asset{"https://github.com/hanselime/paqet/releases/download/v1.0.0-alpha.21/paqet-linux-" + arch + "-v1.0.0-alpha.21.tar.gz", h, "paqet"}, nil
	case "alghadir":
		return udp2rawAsset(arch)
	}
	return Asset{}, nil
}
func InstallDependencies(ctx context.Context, kind string, out io.Writer) error {
	if runtime.GOOS != "linux" || os.Geteuid() != 0 {
		return fmt.Errorf("install dependencies as root on Linux")
	}
	packages := []string{"iproute2", "iputils-ping"}
	switch kind {
	case "awg":
		// AWG uses a pinned userspace build, independent of kernel packages.
	case "paqet":
		packages = append(packages, "libpcap-dev", "iptables")
	case "ssh":
		packages = append(packages, "openssh-client")
	case "alghadir":
		packages = append(packages, "obfs4proxy", "iptables")
	}
	// Use the administrator's configured package sources. Never add an Ubuntu
	// PPA to Debian, or silently replace an existing kernel/module installation.
	if e := runCommand(ctx, command("apt-get", "update")); e != nil {
		return e
	}
	if e := runCommand(ctx, command("apt-get", append([]string{"install", "-y"}, packages...)...)); e != nil {
		return fmt.Errorf("dependencies: %w", e)
	}
	if kind == "awg" {
		return installAWG(ctx)
	}
	a, e := CoreAsset(kind, runtime.GOARCH)
	if e != nil {
		return e
	}
	if a.URL == "" {
		return nil
	}
	name := "rgt"
	if kind == "paqet" {
		name = "paqet"
	}
	if kind == "alghadir" {
		name = "udp2raw"
	}
	return installAsset(ctx, a, filepath.Join(CoreDir, name), out)
}
func installAsset(ctx context.Context, a Asset, path string, out io.Writer) error {
	fmt.Fprintf(out, "Downloading pinned core: %s\n", a.URL)
	req, e := http.NewRequestWithContext(ctx, "GET", a.URL, nil)
	if e != nil {
		return e
	}
	client := http.Client{Timeout: 3 * time.Minute}
	r, e := client.Do(req)
	if e != nil {
		return e
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("core download: HTTP %d", r.StatusCode)
	}
	b, e := io.ReadAll(io.LimitReader(r.Body, 64<<20))
	if e != nil {
		return e
	}
	h := sha256.Sum256(b)
	if hex.EncodeToString(h[:]) != a.SHA256 {
		return fmt.Errorf("core checksum mismatch; nothing installed")
	}
	core, e := extractCore(b, a.Member)
	if e != nil {
		return e
	}
	if len(core) < 4 || !bytes.Equal(core[:4], []byte{127, 'E', 'L', 'F'}) {
		return fmt.Errorf("archive member is not a Linux executable")
	}
	if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		return e
	}
	tmp, e := os.CreateTemp(filepath.Dir(path), ".core-")
	if e != nil {
		return e
	}
	defer os.Remove(tmp.Name())
	if _, e = tmp.Write(core); e != nil {
		tmp.Close()
		return e
	}
	if e = tmp.Chmod(0700); e != nil {
		tmp.Close()
		return e
	}
	if e = tmp.Close(); e != nil {
		return e
	}
	return os.Rename(tmp.Name(), path)
}
func extractCore(b []byte, member string) ([]byte, error) {
	matches := func(name string) bool {
		return filepath.Base(name) == member || (member == "paqet" && strings.HasPrefix(filepath.Base(name), "paqet_linux_"))
	}
	if bytes.HasPrefix(b, []byte("PK")) {
		z, e := zip.NewReader(bytes.NewReader(b), int64(len(b)))
		if e != nil {
			return nil, e
		}
		for _, f := range z.File {
			if !f.FileInfo().IsDir() && matches(f.Name) {
				r, e := f.Open()
				if e != nil {
					return nil, e
				}
				b, e := io.ReadAll(io.LimitReader(r, 32<<20))
				r.Close()
				return b, e
			}
		}
	} else {
		g, e := gzip.NewReader(bytes.NewReader(b))
		if e != nil {
			return nil, e
		}
		defer g.Close()
		t := tar.NewReader(g)
		for {
			h, e := t.Next()
			if e == io.EOF {
				break
			}
			if e != nil {
				return nil, e
			}
			if h.Typeflag == tar.TypeReg && matches(h.Name) {
				return io.ReadAll(io.LimitReader(t, 32<<20))
			}
		}
	}
	return nil, fmt.Errorf("archive does not contain %s", member)
}
