package manage

import (
	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
	"strings"
)

func setDaggerSpoof(s *externaltunnel.Spec, a [4]string) {
	s.IranSpoofSource, s.KharejSpoofSource = a[0], a[1]
	s.IranSpoofDestination, s.KharejSpoofDestination = a[2], a[3]
}
func promptDaggerSpoof() ([4]string, bool) {
	tui.Info("IP Spoof Changes Outer IPv4 Headers. Enter Addresses That Your Paired Network Can Actually Deliver. Blank Keeps The Real Address; bk Does Not Change Your Routes.")
	labels := []string{"Iran Outer Source IPv4 (Blank = Real Iran IP)", "Kharej Outer Source IPv4 (Blank = Real Kharej IP)", "Iran Outer Destination IPv4 (Blank = Real Kharej IP)", "Kharej Outer Destination IPv4 (Blank = Real Iran IP)"}
	var a [4]string
	for i, label := range labels {
		a[i] = strings.TrimSpace(tui.PromptDefault(label, ""))
	}
	s := externaltunnel.New("spoof-check", "d3-quantum-tcp-spoof", "iran")
	s.LocalIP, s.PeerIP = "192.0.2.1", "192.0.2.2"
	setDaggerSpoof(&s, a)
	if err := s.Validate(); err != nil {
		tui.Error(err.Error())
		tui.PressEnter()
		return a, false
	}
	return a, true
}
