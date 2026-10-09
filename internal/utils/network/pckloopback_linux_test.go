//go:build linux

package network

import (
	"bytes"
	"errors"
	"net"
	"os"
	"testing"
	"time"
)

// Exercise real AF_PACKET sockets, not a mocked frame: /proc/net/route does
// not contain the local routing table. A default Ethernet interface must not
// capture the listener away from loopback, or supply its reply source address.
func TestPckLoopback(t *testing.T) {
	for _, target := range []string{"127.0.0.1", "127.0.0.2"} {
		t.Run(target, func(t *testing.T) {
			ifaceName, gateway, err := routeToward(net.ParseIP(target))
			if err != nil {
				t.Fatal(err)
			}
			iface, err := net.InterfaceByName(ifaceName)
			if err != nil || iface.Flags&net.FlagLoopback == 0 || gateway != nil {
				t.Fatalf("local route: interface=%v gateway=%v error=%v", iface, gateway, err)
			}
			reservation, err := net.Listen("tcp4", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			port := uint16(reservation.Addr().(*net.TCPAddr).Port)
			reservation.Close()
			serverPC, err := newPckConn(true, port, PcapCarrier{Port: port, Token: t.Name()})
			if errors.Is(err, os.ErrPermission) && os.Getenv("BK_REQUIRE_PCK") != "1" {
				t.Skip("real PCK sockets require root or CAP_NET_RAW")
			}
			if err != nil {
				t.Fatal(err)
			}
			defer serverPC.Close()
			clientPC, err := newPckConn(false, port, PcapCarrier{Port: port, PeerIP: target, Token: t.Name()})
			if err != nil {
				t.Fatal(err)
			}
			defer clientPC.Close()
			server, client := serverPC.(*pckConn), clientPC.(*pckConn)
			server.SetDeadline(time.Now().Add(5 * time.Second))
			client.SetDeadline(time.Now().Add(5 * time.Second))
			destination := &net.UDPAddr{IP: net.ParseIP(target), Port: int(port)}
			request := bytes.Repeat([]byte("bk-pck-request"), 31)
			reply := bytes.Repeat([]byte("bk-pck-reply"), 37)
			if _, err := client.WriteTo(request, destination); err != nil {
				t.Fatal(err)
			}
			bufs := [][]byte{make([]byte, 2048)}
			sizes := make([]int, 1)
			froms := make([]net.Addr, 1)
			n, err := server.ReadBatch(bufs, sizes, froms)
			if err != nil || n != 1 || !bytes.Equal(bufs[0][:sizes[0]], request) {
				t.Fatalf("server batch received %d packets, error=%v", n, err)
			}
			// Loopback has no Ethernet injection. Callers must fall back to
			// WriteTo when the batch path declines it.
			if _, err := server.WriteBatch([][]byte{reply}, froms[0]); !errors.Is(err, errPckNoBatch) {
				t.Fatalf("loopback batch must decline L2 injection: %v", err)
			}
			if _, err := server.WriteTo(reply, froms[0]); err != nil {
				t.Fatal(err)
			}
			buf := make([]byte, 2048)
			got, from, err := client.ReadFrom(buf)
			if err != nil || !bytes.Equal(buf[:got], reply) {
				t.Fatalf("reply: bytes=%d error=%v", got, err)
			}
			if !from.(*net.UDPAddr).IP.Equal(destination.IP) {
				t.Fatalf("reply source %v differs from contacted address %v", from, destination)
			}
		})
	}
}
