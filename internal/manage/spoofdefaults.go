package manage

import (
	"github.com/topgsmir/bk/config"
	"github.com/topgsmir/bk/internal/utils/network"
)

const spoofUDPBuffer = 8 << 20

func isSpoofUDP(sc config.SpoofConfig) bool {
	up, down := network.ResolveSpoofDirections(sc.SpoofProfile, sc.SpoofUplink, sc.SpoofDownlink)
	return up == "udp" && down == "udp"
}

// New UDP setups rotate source ports independently of padding/Stealth.
// Existing configs and explicit edits are never migrated by this helper.
func newSpoofConfig(profile, uplink, downlink string) config.SpoofConfig {
	sc := config.SpoofConfig{SpoofProfile: profile, SpoofUplink: uplink, SpoofDownlink: downlink}
	if isSpoofUDP(sc) {
		sc.SpoofShufflePort = true
		sc.SpoofPortMin, sc.SpoofPortMax = 49152, 65535
		sc.SpoofSockBuf = spoofUDPBuffer
	}
	return sc
}

// Use a conservative inner MTU for raw UDP; explicit advanced settings and
// imported MTU values are applied afterwards.
func applyNewL3CarrierDefaults(s *l3Spec) {
	if s.Carrier != "spoof" || !isSpoofUDP(s.Spoof) {
		return
	}
	s.MTU = 1280
	auto := false
	s.AutoMTU = &auto
	if s.Spoof.SpoofSockBuf == 0 {
		s.Spoof.SpoofSockBuf = spoofUDPBuffer
	}
}
