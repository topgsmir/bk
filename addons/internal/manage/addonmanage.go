package manage

import (
	"fmt"
	"github.com/topgsmir/bk/addons/internal/tui"
)

// ManageTunnels owns only optional-package configurations. Core management
// remains in the independently installed BackPack component.
func ManageTunnels() {
	for {
		tunnels, warnings := listManagedTunnels(nil, externalDir)
		for _, warning := range warnings {
			tui.Warn(warning)
		}
		if len(tunnels) == 0 {
			tui.Warn("No Package 2 Or 3 Tunnels Yet.")
			tui.PressEnter()
			return
		}
		options := make([]tui.Option, len(tunnels))
		for i, t := range tunnels {
			state := "stopped"
			if IsActive(t.tunnel.Service) {
				state = "running"
			}
			options[i] = tui.Option{Title: t.tunnel.Name, Desc: fmt.Sprintf("Package %d | %s %s — %s", t.packageNumber, t.tunnel.Role, t.tunnel.Transport, state)}
		}
		choice := tui.ChooseOpt("Package 2 And 3 Tunnels", options)
		if choice < 0 {
			return
		}
		manageOneExternal(*tunnels[choice].external)
	}
}
