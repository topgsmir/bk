//go:build linux

package externaltunnel

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func TestDaggerNativeRejectsWrongIdentity(t *testing.T) {
	if os.Getenv("BK_LIVE_DAGGER") != "1" {
		t.Skip("mandatory in Dagger native CI")
	}
	if err := checkDaggerCore(); err != nil {
		t.Fatal(err)
	}
	port := func() int {
		c, e := net.Listen("tcp4", "127.0.0.1:0")
		if e != nil {
			t.Fatal(e)
		}
		defer c.Close()
		return c.Addr().(*net.TCPAddr).Port
	}
	for _, kind := range []string{"d3-tcp", "d3-https"} {
		t.Run(kind, func(t *testing.T) {
			backend, err := net.Listen("tcp4", "127.0.0.2:0")
			if err != nil {
				t.Fatal(err)
			}
			defer backend.Close()
			go func() {
				for {
					c, e := backend.Accept()
					if e != nil {
						return
					}
					go func() { defer c.Close(); io.Copy(c, c) }()
				}
			}()
			s := New("auth-negative", kind, "iran")
			s.LocalIP, s.PeerIP = "127.0.0.1", "127.0.0.2"
			s.Port = port()
			s.Listen = fmt.Sprintf("127.0.0.1:%d", port())
			s.Target = backend.Addr().String()
			start := func(spec Spec) func() {
				ctx, cancel := context.WithCancel(context.Background())
				done := make(chan error, 1)
				dir := t.TempDir()
				go func() { done <- runDagger(ctx, spec, dir) }()
				return func() {
					cancel()
					select {
					case <-done:
					case <-time.After(7 * time.Second):
						t.Error("native Dagger did not stop")
					}
				}
			}
			stopServer := start(s)
			defer stopServer()
			echo := func(deadline time.Duration) error {
				end := time.Now().Add(deadline)
				var last error
				for time.Now().Before(end) {
					c, e := net.DialTimeout("tcp4", s.Listen, 200*time.Millisecond)
					if e != nil {
						last = e
						time.Sleep(50 * time.Millisecond)
						continue
					}
					c.SetDeadline(time.Now().Add(400 * time.Millisecond))
					_, e = c.Write([]byte("authenticated"))
					if e == nil {
						body := make([]byte, 13)
						_, e = io.ReadFull(c, body)
						if e == nil && string(body) != "authenticated" {
							e = fmt.Errorf("payload mismatch")
						}
					}
					c.Close()
					if e == nil {
						return nil
					}
					last = e
				}
				return last
			}
			stopGood := start(s.Mirror())
			if err := echo(8 * time.Second); err != nil {
				stopGood()
				t.Fatal("valid native pair failed", err)
			}
			stopGood()
			wrong := s.Mirror()
			wrong.Secret = NewSecret()
			stopBad := start(wrong)
			if err := echo(3 * time.Second); err == nil {
				stopBad()
				t.Fatal("wrong paired identity carried real payload")
			}
			stopBad()
			stopAgain := start(s.Mirror())
			defer stopAgain()
			if err := echo(8 * time.Second); err != nil {
				t.Fatal("paired identity did not recover after wrong-key rejection", err)
			}
		})
	}
}
