package manage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/topgsmir/bk/addons/internal/externaltunnel"
)

// Keep the optional configs separate from core.List: its other callers expect
// TOML configs and must not send an optional tunnel to the original editor.
type managedTunnel struct {
	packageNumber int
	tunnel        Tunnel
	external      *externaltunnel.Spec
}

func listManagedTunnels(original []Tunnel, directory string) ([]managedTunnel, []string) {
	tunnels := make([]managedTunnel, 0, len(original))
	for _, t := range original {
		tunnels = append(tunnels, managedTunnel{packageNumber: 1, tunnel: t})
	}
	var warnings []string
	files, err := os.ReadDir(directory)
	if err != nil && !os.IsNotExist(err) {
		warnings = append(warnings, "Cannot Read Package 2 And 3 Tunnel Configurations.")
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		s, err := externaltunnel.Load(filepath.Join(directory, file.Name()))
		// Actions address configs by their validated name. Do not offer a file
		// whose name would make those actions target a different config.
		if err != nil || file.Name() != s.Name+".json" || file.Type()&os.ModeSymlink != 0 {
			warnings = append(warnings, fmt.Sprintf("Cannot Manage Tunnel Config: %s", file.Name()))
			continue
		}
		pkg := 2
		for _, kind := range externaltunnel.Package3Kinds {
			if kind.ID == s.Kind {
				pkg = 3
				break
			}
		}
		tunnels = append(tunnels, managedTunnel{
			packageNumber: pkg,
			tunnel: Tunnel{Name: s.Name, Role: s.Side, Transport: s.Kind,
				Service: externaltunnel.ServiceName(s.Name)},
			external: &s,
		})
	}
	sort.Slice(tunnels, func(i, j int) bool {
		if tunnels[i].packageNumber != tunnels[j].packageNumber {
			return tunnels[i].packageNumber < tunnels[j].packageNumber
		}
		return tunnels[i].tunnel.Name < tunnels[j].tunnel.Name
	})
	return tunnels, warnings
}
