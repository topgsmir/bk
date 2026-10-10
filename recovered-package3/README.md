# Package 3 binary analysis

This folder contains analysis of the provided Backhaul and Solarpass binaries. Original files under `new2` have not been changed. These results are not the original complete, buildable source code.

## Backhaul

- Recovered 11,612 Go runtime function records, with corrected text-base addresses in `Backhaul/go-function-table.txt`.
- Exported 12 C-like diagnostic functions. The original application names and file paths were obfuscated; several Go library names remain available.
- The “no public IP address found” path enumerates locally assigned interface addresses and filters them. It is not evidence of a failed remote IP lookup. Both lab peers initially had only private addresses; adding lab-only documentation-range addresses made startup succeed.
- All nine supplied stream types transferred actual TCP payloads in the isolated lab. WSS, WSSMUX and AnyTLS require certificate files.

## Solarpass Spoof

Two independent faults were reproduced:

1. Iran binds a Tokio frontend UDP socket, then binds a second return UDP socket to the same wildcard address and port. The first socket does not have `SO_REUSEADDR`; the second fails with `EADDRINUSE`.
2. The receive thread compares byte-swapped `sockaddr_in.sin_addr` with a stored native-endian peer filter. Correct IPv4 configuration can consequently reject every received packet. See `Solarpass/spoof-receive-disassembly.txt`, offsets `0xda011`–`0xda01a`.

A narrowly scoped, version-pinned repair exists in a separate ELF copy: change only the Iran Spoof return-socket call to use frontend port + 1. The manifest records exact offsets and original/replacement bytes. The outside peer must send returns to that new port. The pinned core's filter can be corrected by reversing the four peer-IP octets **only in `spoof.peer_spoof_ip`**; actual endpoints and routing addresses remain normal. This is a version-specific compatibility measure, not a recommendation to reverse ordinary server addresses.

The initial repaired-core round trip passed against the actual two native cores. The separate native-core test passed 60/60 byte-for-byte UDP exchanges over 60 seconds, with payloads of 5, 68, 516, 1204 and 1404 bytes. The adapter applies only the hash-verified repair to a separate copy and preserves the original provided binary.

`Solarpass/function-boundaries.tsv` contains 8,945 function ranges recovered from unwind metadata. Automatic Rust decompilation produced unresolved switch-flow errors and timed out on selected functions; no successful full-source reconstruction is claimed. `spoof-receive-recovered.c` is a readable manual reconstruction of the confirmed receive-loop logic, not directly compilable original source.

No licensing, signature, authentication, or vendor-access checks were patched. The only executable change is the diagnosed return-port call and its short trampoline in verified inter-function padding. The original hashes are recorded alongside the patch manifest.

The bk adapter also uses a stable socket towards the native UDP core and an authenticated, encrypted session envelope between local adapters. This preserves the native carrier while identifying each application's replies, including when the connection-test phases open a new UDP socket. The adapter runs the core in a private mount namespace with read-only `/proc/sys`, and restricts its internal frontend and synchronization listeners using temporary, named firewall rules. This is new addon code, not a reconstruction of the original complete Solarpass source.
