package manage

import (
	"fmt"

	"github.com/topgsmir/bk/addons/internal/tui"
)

// plainState returns "running"/"stopped" without colors (for gray descriptions).
func plainState(service string) string {
	if IsActive(service) {
		return "running"
	}
	return "stopped"
}

func report(err error, action string) {
	if err != nil {
		tui.Error(fmt.Sprintf("Failed: %v", err))
	} else {
		tui.Success("Tunnel " + action + ".")
	}
	tui.PressEnter()
}
