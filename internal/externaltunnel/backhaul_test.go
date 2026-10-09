package externaltunnel

import (
	"context"
	"crypto/tls"
	"net"
	"testing"
	"time"
)

func TestBackhaulPairedTLSAndWrongSecret(t *testing.T) {
	s := New("tls-case", "b3-tcp", "iran")
	for _, wrong := range []bool{false, true} {
		server, err := backhaulTLS(s, t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		other := s.Mirror()
		if wrong {
			other.Secret = NewSecret()
		}
		client, err := backhaulTLS(other, t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		listener, err := tls.Listen("tcp", "127.0.0.1:0", server)
		if err != nil {
			t.Fatal(err)
		}
		done := make(chan error, 1)
		go func() {
			conn, err := listener.Accept()
			if err == nil {
				defer conn.Close()
				conn.SetDeadline(time.Now().Add(time.Second))
				err = conn.(*tls.Conn).Handshake()
			}
			done <- err
		}()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		dialer := tls.Dialer{Config: client, NetDialer: &net.Dialer{}}
		conn, clientErr := dialer.DialContext(ctx, "tcp", listener.Addr().String())
		if conn != nil {
			conn.Close()
		}
		cancel()
		listener.Close()
		serverErr := <-done
		if !wrong && (clientErr != nil || serverErr != nil) {
			t.Fatal("paired TLS", clientErr, serverErr)
		}
		if wrong && clientErr == nil {
			t.Fatal("wrong paired secret accepted")
		}
	}
}
