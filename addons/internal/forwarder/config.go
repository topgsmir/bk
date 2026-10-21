package forwarder

// Config is owned by the addon TCP/UDP forwarder, independently of BackPack.
type Config struct {
	Ports          []string
	PeerIP         string
	AcceptUDP      bool
	MaxConnections int
	BandwidthMbps  int
}
