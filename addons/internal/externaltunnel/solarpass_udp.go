package externaltunnel

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

// Solarpass Spoof does not identify application UDP sessions. One stable socket
// towards the native core plus an authenticated session envelope prevents replies
// being assigned to the wrong application when test phases or clients change.
// Native Solarpass still carries every byte between the two server namespaces.
const solarUDPHeader = 41
const solarUDPMaxSessions = 64

type solarUDPID [16]byte

func solarUDPFrame(key []byte, id solarUDPID, side byte, seq uint64, body []byte) ([]byte, error) {
	if seq == ^uint64(0) || len(body) > 65450 {
		return nil, fmt.Errorf("Solarpass UDP frame limit reached")
	}
	header := make([]byte, solarUDPHeader)
	copy(header, "BK3U")
	header[4] = side
	copy(header[5:21], id[:])
	binary.BigEndian.PutUint64(header[21:29], seq)
	if _, err := rand.Read(header[29:41]); err != nil {
		return nil, err
	}
	aead, err := solarUDPAEAD(key, id, side)
	if err != nil {
		return nil, err
	}
	return aead.Seal(header, header[29:41], body, header), nil
}
func solarUDPAEAD(key []byte, id solarUDPID, side byte) (cipher.AEAD, error) {
	h := hmac.New(sha256.New, key)
	h.Write([]byte("bk/package3/solarpass/udp/session/v1"))
	h.Write(id[:])
	h.Write([]byte{side})
	block, err := aes.NewCipher(h.Sum(nil))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
func solarUDPDecode(key, frame []byte, expected byte) (id solarUDPID, seq uint64, body []byte, err error) {
	if len(frame) < solarUDPHeader+16 || string(frame[:4]) != "BK3U" || frame[4] != expected {
		err = fmt.Errorf("invalid Solarpass UDP envelope")
		return
	}
	copy(id[:], frame[5:21])
	seq = binary.BigEndian.Uint64(frame[21:29])
	aead, e := solarUDPAEAD(key, id, expected)
	if e != nil {
		err = e
		return
	}
	body, err = aead.Open(nil, frame[29:41], frame[solarUDPHeader:], frame[:solarUDPHeader])
	return
}

type solarUDPWindow struct {
	initialized bool
	max, bits   uint64
}

func (w *solarUDPWindow) accept(seq uint64) bool {
	if !w.initialized {
		w.initialized = true
		w.max = seq
		w.bits = 1
		return true
	}
	if seq > w.max {
		delta := seq - w.max
		if delta >= 64 {
			w.bits = 1
		} else {
			w.bits = w.bits<<delta | 1
		}
		w.max = seq
		return true
	}
	delta := w.max - seq
	if delta >= 64 || w.bits&(uint64(1)<<delta) != 0 {
		return false
	}
	w.bits |= uint64(1) << delta
	return true
}
func runSolarpassUDPBridge(ctx context.Context, s Spec) error {
	return runPackage3UDPBridge(ctx, s, "solarpass/udp/envelope")
}
func runPackage3UDPBridge(ctx context.Context, s Spec, label string) error {
	key, err := base64.StdEncoding.DecodeString(s.Key(label))
	if err != nil {
		return err
	}
	if s.Side == "iran" {
		return solarpassUDPIran(ctx, s, key)
	}
	return solarpassUDPKharej(ctx, s, key)
}
func solarpassUDPIran(ctx context.Context, s Spec, key []byte) error {
	address, err := net.ResolveUDPAddr("udp4", s.Listen)
	if err != nil {
		return err
	}
	front, err := net.ListenUDP("udp4", address)
	if err != nil {
		return err
	}
	defer front.Close()
	upstreamAddress := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: s.SourcePort}
	upstream, err := net.DialUDP("udp4", nil, upstreamAddress)
	if err != nil {
		return err
	}
	defer upstream.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() { <-ctx.Done(); front.Close(); upstream.Close() }()
	type session struct {
		id       solarUDPID
		client   *net.UDPAddr
		seen     time.Time
		sent     uint64
		received solarUDPWindow
	}
	byClient := map[string]*session{}
	byID := map[solarUDPID]*session{}
	var mu sync.Mutex
	var reader sync.WaitGroup
	reader.Add(1)
	defer func() { cancel(); upstream.Close(); reader.Wait() }()
	go func() {
		defer reader.Done()
		buffer := make([]byte, 65535)
		for {
			n, err := upstream.Read(buffer)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				time.Sleep(20 * time.Millisecond)
				continue
			}
			id, seq, body, err := solarUDPDecode(key, buffer[:n], 1)
			if err != nil {
				continue
			}
			mu.Lock()
			peer := byID[id]
			if peer == nil || !peer.received.accept(seq) {
				mu.Unlock()
				continue
			}
			client := peer.client
			peer.seen = time.Now()
			mu.Unlock()
			front.WriteToUDP(body, client)
		}
	}()
	for {
		buffer := make([]byte, 65535)
		n, client, err := front.ReadFromUDP(buffer)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		mu.Lock()
		peer := byClient[client.String()]
		if peer != nil && time.Since(peer.seen) > 50*time.Second {
			delete(byClient, client.String())
			delete(byID, peer.id)
			peer = nil
		}
		if peer == nil {
			for old, p := range byClient {
				if time.Since(p.seen) > 60*time.Second {
					delete(byClient, old)
					delete(byID, p.id)
				}
			}
			if len(byClient) >= solarUDPMaxSessions {
				mu.Unlock()
				continue
			}
			peer = &session{client: client}
			if _, err = rand.Read(peer.id[:]); err != nil {
				mu.Unlock()
				return err
			}
			byClient[client.String()] = peer
			byID[peer.id] = peer
		}
		peer.seen = time.Now()
		seq := peer.sent
		peer.sent++
		id := peer.id
		mu.Unlock()
		frame, err := solarUDPFrame(key, id, 0, seq, buffer[:n])
		if err != nil {
			return err
		}
		upstream.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if _, err = upstream.Write(frame); err != nil && ctx.Err() != nil {
			return nil
		}
	}
}
func solarpassUDPKharej(ctx context.Context, s Spec, key []byte) error {
	front, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: s.SourcePort + 3})
	if err != nil {
		return err
	}
	defer front.Close()
	target, err := net.ResolveUDPAddr("udp4", s.Target)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() { <-ctx.Done(); front.Close() }()
	type session struct {
		conn     *net.UDPConn
		core     *net.UDPAddr
		seen     time.Time
		sent     uint64
		received solarUDPWindow
	}
	sessions := map[solarUDPID]*session{}
	var mu sync.Mutex
	var readers sync.WaitGroup
	defer func() {
		cancel()
		mu.Lock()
		for _, p := range sessions {
			if p.conn != nil {
				p.conn.Close()
			}
		}
		mu.Unlock()
		readers.Wait()
	}()
	startReader := func(id solarUDPID, p *session, conn *net.UDPConn) {
		readers.Add(1)
		go func() {
			defer readers.Done()
			defer conn.Close()
			defer func() {
				mu.Lock()
				if p.conn == conn {
					p.conn = nil
				}
				mu.Unlock()
			}()
			buffer := make([]byte, 65535)
			for {
				conn.SetReadDeadline(time.Now().Add(65 * time.Second))
				n, err := conn.Read(buffer)
				if err != nil {
					return
				}
				mu.Lock()
				// Keep counters and replay state after a target socket error. A reopened
				// target must not reset the reply counter of an active client session.
				seq := p.sent
				p.sent++
				address := p.core
				p.seen = time.Now()
				mu.Unlock()
				frame, err := solarUDPFrame(key, id, 1, seq, buffer[:n])
				if err != nil {
					return
				}
				front.WriteToUDP(frame, address)
			}
		}()
	}
	for {
		buffer := make([]byte, 65535)
		n, core, err := front.ReadFromUDP(buffer)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		id, seq, body, err := solarUDPDecode(key, buffer[:n], 0)
		if err != nil {
			continue
		}
		mu.Lock()
		p := sessions[id]
		if p == nil {
			for old, peer := range sessions {
				if time.Since(peer.seen) > 60*time.Second {
					if peer.conn != nil {
						peer.conn.Close()
					}
					delete(sessions, old)
				}
			}
			if len(sessions) >= solarUDPMaxSessions {
				mu.Unlock()
				continue
			}
			p = &session{core: core, seen: time.Now()}
			sessions[id] = p
		}
		if !p.received.accept(seq) {
			mu.Unlock()
			continue
		}
		p.seen = time.Now()
		p.core = core
		if p.conn == nil {
			conn, e := net.DialUDP("udp4", nil, target)
			if e != nil {
				mu.Unlock()
				return e
			}
			p.conn = conn
			startReader(id, p, conn)
		}
		conn := p.conn
		mu.Unlock()
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if _, err = conn.Write(body); err != nil {
			mu.Lock()
			if p.conn == conn {
				p.conn = nil
			}
			mu.Unlock()
			conn.Close()
		}
	}
}
