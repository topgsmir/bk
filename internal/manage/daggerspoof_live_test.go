//go:build linux

package manage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"net"
	"sync"
	"syscall"
	"testing"
)

// Observe the real Ethernet frames, independently of a successful echo. Never
// print payloads, key material or complete packets.
func observeDaggerSpoofHeaders(t *testing.T, kinds []string, host, peer string, spoof [4]string) func() {
	t.Helper()
	protocols := map[string]byte{"tcp": 6, "udp": 17, "icmp": 1, "bip": 1, "gre": 47, "ipip": 4, "raw": 253}
	wanted := map[string]bool{}
	for _, kind := range kinds {
		if !externaltunnel.DaggerSpoof(kind) {
			continue
		}
		_, profile, _, _ := externaltunnel.DaggerOptions(kind)
		for _, direction := range []string{"iran", "kharej"} {
			wanted[fmt.Sprintf("%s/%d", direction, protocols[profile])] = true
		}
	}
	if len(wanted) == 0 {
		return func() {}
	}
	iface, err := net.InterfaceByName("bk-test-host")
	if err != nil {
		t.Fatal(err)
	}
	protocol := binary.NativeEndian.Uint16([]byte{0x08, 0x00})
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(protocol))
	if err != nil {
		t.Fatal(err)
	}
	if err = syscall.Bind(fd, &syscall.SockaddrLinklayer{Protocol: protocol, Ifindex: iface.Index}); err != nil {
		syscall.Close(fd)
		t.Fatal(err)
	}
	if err = syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{Usec: 300000}); err != nil {
		syscall.Close(fd)
		t.Fatal(err)
	}
	endpoints := [2][2]string{{host, peer}, {peer, host}}
	if spoof[0] != "" {
		endpoints[0][0] = spoof[0]
	}
	if spoof[1] != "" {
		endpoints[1][0] = spoof[1]
	}
	if spoof[2] != "" {
		endpoints[0][1] = spoof[2]
	}
	if spoof[3] != "" {
		endpoints[1][1] = spoof[3]
	}
	stop, done := make(chan struct{}), make(chan struct{})
	seen := map[string]bool{}
	var captureErr error
	go func() {
		defer close(done)
		frame := make([]byte, 65549)
		for {
			select {
			case <-stop:
				return
			default:
			}
			n, _, e := syscall.Recvfrom(fd, frame, 0)
			if e == syscall.EAGAIN || e == syscall.EWOULDBLOCK || e == syscall.EINTR {
				continue
			}
			if e != nil {
				captureErr = e
				return
			}
			if n < 34 || frame[12] != 8 || frame[13] != 0 {
				continue
			}
			ip := frame[14:n]
			ihl := int(ip[0]&15) * 4
			total := int(binary.BigEndian.Uint16(ip[2:4]))
			if ip[0]>>4 != 4 || ihl < 20 || total < ihl || total > len(ip) {
				continue
			}
			body := ip[ihl:total]
			offset := 0
			switch ip[9] {
			case 6:
				if len(body) < 20 {
					continue
				}
				offset = int(body[12]>>4) * 4
			case 17, 1:
				offset = 8
			case 47:
				offset = 4
			case 4, 253:
			default:
				continue
			}
			if offset > len(body) {
				continue
			}
			payload := body[offset:]
			if !bytes.HasPrefix(payload, []byte("DRTU\x01")) && !(len(payload) > 6 && payload[4] == 0xf1 && payload[5] == 0 || len(payload) > 6 && payload[4] == 0xf2 && payload[5] == 0) {
				continue
			}
			direction := ""
			for i, addresses := range endpoints {
				if net.IP(ip[12:16]).String() == addresses[0] && net.IP(ip[16:20]).String() == addresses[1] {
					direction = []string{"iran", "kharej"}[i]
					break
				}
			}
			if direction == "" {
				continue
			}
			sum := uint32(0)
			for i := 0; i < ihl; i += 2 {
				sum += uint32(binary.BigEndian.Uint16(ip[i : i+2]))
			}
			for sum>>16 != 0 {
				sum = (sum & 65535) + (sum >> 16)
			}
			if uint16(sum) != 65535 {
				captureErr = fmt.Errorf("invalid outer IPv4 checksum")
				return
			}
			seen[fmt.Sprintf("%s/%d", direction, ip[9])] = true
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() {
			close(stop)
			<-done
			syscall.Close(fd)
			if captureErr != nil {
				t.Errorf("spoof header observation: %v", captureErr)
			}
			for key := range wanted {
				if !seen[key] {
					t.Errorf("no real spoofed source/destination header observed for %s", key)
				}
			}
			t.Logf("Verified actual outer spoof source/destination headers and IPv4 checksums in both directions for %d carrier/protocol witnesses", len(seen))
		})
	}
}
