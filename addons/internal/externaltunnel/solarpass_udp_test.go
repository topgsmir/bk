package externaltunnel

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestSolarpassUDPEnvelopeAuthentication(t *testing.T) {
	key := bytes.Repeat([]byte{7}, 32)
	var id solarUDPID
	id[0] = 9
	frame, err := solarUDPFrame(key, id, 0, 42, []byte("test payload"))
	if err != nil {
		t.Fatal(err)
	}
	gotID, seq, body, err := solarUDPDecode(key, frame, 0)
	if err != nil || gotID != id || seq != 42 || string(body) != "test payload" {
		t.Fatal(gotID, seq, string(body), err)
	}
	// Tampering with any routing field, sequence, nonce, payload or tag fails.
	for _, position := range []int{4, 5, 21, 29, 41, len(frame) - 1} {
		changed := append([]byte(nil), frame...)
		changed[position] ^= 1
		if _, _, _, err := solarUDPDecode(key, changed, 0); err == nil {
			t.Fatal("tampering accepted", position)
		}
	}
	if _, _, _, err := solarUDPDecode(key, frame, 1); err == nil {
		t.Fatal("reflected request accepted as reply")
	}
	if _, _, _, err := solarUDPDecode(bytes.Repeat([]byte{8}, 32), frame, 0); err == nil {
		t.Fatal("wrong peer key accepted")
	}
	repeated, err := solarUDPFrame(key, id, 0, 42, []byte("test payload"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(frame[29:41], repeated[29:41]) {
		t.Fatal("nonce reused across session reconstruction")
	}
	if _, err := solarUDPFrame(key, id, 0, ^uint64(0), nil); err == nil {
		t.Fatal("sequence overflow accepted")
	}
}
func TestSolarpassUDPReplayAndReordering(t *testing.T) {
	var window solarUDPWindow
	for _, seq := range []uint64{100, 102, 101, 163, 162} {
		if !window.accept(seq) {
			t.Fatal("valid reordered packet rejected", seq)
		}
	}
	for _, seq := range []uint64{162, 163, 100, 99} {
		if window.accept(seq) {
			t.Fatal("duplicate or old packet accepted", seq)
		}
	}
	if !window.accept(1000) || window.accept(900) {
		t.Fatal("replay window jump failed")
	}
}

func TestPackage3UDPBackendRestartPreservesSession(t *testing.T) {
	backend, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	target := backend.LocalAddr().(*net.UDPAddr)
	echo := func(conn *net.UDPConn) {
		go func() {
			buffer := make([]byte, 4096)
			for {
				n, a, e := conn.ReadFromUDP(buffer)
				if e != nil {
					return
				}
				conn.WriteToUDP(buffer[:n], a)
			}
		}()
	}
	echo(backend)
	defer func() { backend.Close() }()
	reserve, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	port := reserve.LocalAddr().(*net.UDPAddr).Port
	reserve.Close()
	s := New("udp-restart", "s3-spoof", "kharej")
	s.SourcePort = port - 3
	s.Target = target.String()
	key := bytes.Repeat([]byte{17}, 32)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- solarpassUDPKharej(ctx, s, key) }()
	defer func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(time.Second):
			t.Error("UDP bridge did not stop")
		}
	}()
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var id solarUDPID
	id[0] = 6
	var window solarUDPWindow
	exchange := func(seq uint64) error {
		body := []byte(fmt.Sprintf("packet-%d", seq))
		frame, e := solarUDPFrame(key, id, 0, seq, body)
		if e != nil {
			return e
		}
		conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
		if _, e = conn.Write(frame); e != nil {
			return e
		}
		buffer := make([]byte, 4096)
		n, e := conn.Read(buffer)
		if e != nil {
			return e
		}
		gotID, replySeq, reply, e := solarUDPDecode(key, buffer[:n], 1)
		if e != nil {
			return e
		}
		if gotID != id || !bytes.Equal(body, reply) || !window.accept(replySeq) {
			return fmt.Errorf("wrong reply or replayed/reset sequence %d", replySeq)
		}
		return nil
	}
	// A live counter makes a reset on reconnection observable.
	time.Sleep(20 * time.Millisecond)
	for i := uint64(0); i < 40; i++ {
		if e := exchange(i); e != nil {
			t.Fatal("before restart", i, e)
		}
	}
	backend.Close()
	for i := uint64(40); i < 43; i++ {
		_ = exchange(i)
	}
	backend, err = net.ListenUDP("udp4", target)
	if err != nil {
		t.Fatal(err)
	}
	echo(backend)
	if e := exchange(43); e != nil {
		t.Fatal("after backend restart", e)
	}
}
