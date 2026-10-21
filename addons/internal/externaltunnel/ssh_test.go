package externaltunnel

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestSSHForwardingAuthenticatesPeerAndCarriesRealTraffic(t *testing.T) {
	_, private, e := ed25519.GenerateKey(rand.Reader)
	if e != nil {
		t.Fatal(e)
	}
	signer, _ := ssh.NewSignerFromKey(private)
	_, hostPrivate, _ := ed25519.GenerateKey(rand.Reader)
	hostSigner, _ := ssh.NewSignerFromKey(hostPrivate)
	server, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer server.Close()
	cfg := &ssh.ServerConfig{PublicKeyCallback: func(meta ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		if !bytes.Equal(key.Marshal(), signer.PublicKey().Marshal()) {
			return nil, fmt.Errorf("unauthorized key")
		}
		return nil, nil
	}}
	cfg.AddHostKey(hostSigner)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			raw, e := server.Accept()
			if e != nil {
				return
			}
			go func() {
				conn, chans, reqs, e := ssh.NewServerConn(raw, cfg)
				if e != nil {
					raw.Close()
					return
				}
				defer conn.Close()
				go func() { <-ctx.Done(); conn.Close() }()
				go func() {
					for req := range reqs {
						req.Reply(true, nil)
					}
				}()
				for ch := range chans {
					if ch.ChannelType() != "direct-tcpip" {
						ch.Reject(ssh.UnknownChannelType, "unsupported")
						continue
					}
					channel, reqs, e := ch.Accept()
					if e != nil {
						continue
					}
					go ssh.DiscardRequests(reqs)
					go func() { defer channel.Close(); io.Copy(channel, channel) }()
				}
			}()
		}
	}()
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id")
	der, _ := x509.MarshalPKCS8PrivateKey(private)
	os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600)
	known := filepath.Join(dir, "known_hosts")
	os.WriteFile(known, []byte(knownhosts.Line([]string{server.Addr().String()}, hostSigner.PublicKey())+"\n"), 0600)
	s := fixture("ssh")
	s.PeerIP = "127.0.0.1"
	s.Port = server.Addr().(*net.TCPAddr).Port
	s.SSHKey, s.KnownHosts = keyPath, known
	s.Listen = "127.0.0.1:0"
	probe, _ := net.Listen("tcp", "127.0.0.1:0")
	s.Listen = probe.Addr().String()
	probe.Close()
	done := make(chan error, 1)
	go func() { done <- runSSH(ctx, s) }()
	for n := 0; n < 8; n++ {
		var c net.Conn
		for i := 0; i < 100; i++ {
			c, e = net.DialTimeout("tcp", s.Listen, time.Second)
			if e == nil {
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		if e != nil {
			cancel()
			t.Fatal(e)
		}
		payload := bytes.Repeat([]byte(strconv.Itoa(n)), 32<<10)
		go c.Write(payload)
		b := make([]byte, len(payload))
		c.SetReadDeadline(time.Now().Add(3 * time.Second))
		if _, e = io.ReadFull(c, b); e != nil || !bytes.Equal(payload, b) {
			cancel()
			t.Fatal(e)
		}
		c.Close()
	}
	// A changed host key must fail before opening any forwarded listener.
	bad := s
	bad.KnownHosts = filepath.Join(dir, "wrong_known_hosts")
	os.WriteFile(bad.KnownHosts, []byte(knownhosts.Line([]string{server.Addr().String()}, signer.PublicKey())+"\n"), 0600)
	if e := runSSH(ctx, bad); e == nil || !bytes.Contains([]byte(e.Error()), []byte("key mismatch")) {
		t.Fatalf("untrusted host accepted: %v", e)
	}
	cancel()
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("SSH pool did not stop")
	}
}
