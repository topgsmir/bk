package manage

// linkHost is this server's address as the other end reaches it, for a link
// made from Manage. The wizards ask for it and put it in their link; this link
// was built with none, so a kharej given it through Setup from a link — which
// had no way to ask — was refused with "the server address is required",
// while the same tunnel typed in by hand came up. Reported on v1.8.4 as "the
// setup link did not work for a reverse tunnel, doing it manually did". Empty
// when the address is not known; the other side then asks for it.
var linkHost = func() string {
	if ip := PublicIPv4(); ip != "-" {
		return ip
	}
	return ""
}
