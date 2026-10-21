package externaltunnel

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// Kharej authenticates to Iran and requests a loopback remote listener. Iran's
// own forwarder exposes the configured local port, so sshd GatewayPorts is not needed.
func runSSHReverse(ctx context.Context, s Spec) error {
	established := false
	for ctx.Err() == nil {
		client, raw, e := openSSH(ctx, s)
		if e == nil {
			_ = raw.SetDeadline(time.Now().Add(10 * time.Second))
			var listener net.Listener
			listener, e = client.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(s.SourcePort)))
			_ = raw.SetDeadline(time.Time{})
			if e == nil {
				established = true
				e = serveSSHReverseSession(ctx, client, listener, s.Target)
			} else {
				e = fmt.Errorf("SSH remote forward refused on Iran port %d: %w; check AllowTcpForwarding/PermitListen and port conflicts", s.SourcePort, e)
			}
			client.Close()
		}
		if ctx.Err() != nil {
			return nil
		}
		if !established {
			return e
		}
		fmt.Fprintln(os.Stderr, "SSH reverse reconnecting:", e)
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
		}
	}
	return nil
}
func serveSSHReverseSession(ctx context.Context, client *ssh.Client, listener net.Listener, target string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var connections sync.WaitGroup
	defer func() { cancel(); listener.Close(); client.Close(); connections.Wait() }()
	stopped := make(chan struct{})
	defer close(stopped)
	go func() {
		select {
		case <-ctx.Done():
			client.Close()
			listener.Close()
		case <-stopped:
		}
	}()
	go func() { _ = client.Wait(); cancel(); listener.Close() }()
	go func() {
		tick := time.NewTicker(20 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				done := make(chan error, 1)
				go func() { _, _, e := client.SendRequest("keepalive@openssh.com", true, nil); done <- e }()
				select {
				case e := <-done:
					if e != nil {
						client.Close()
						return
					}
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					client.Close()
					return
				}
			}
		}
	}()
	for {
		incoming, e := listener.Accept()
		if e != nil {
			return e
		}
		connections.Add(1)
		go func(a net.Conn) {
			defer connections.Done()
			b, e := (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", target)
			if e != nil {
				a.Close()
				fmt.Fprintln(os.Stderr, "SSH reverse target:", e)
				return
			}
			relay(ctx, a, b)
		}(incoming)
	}
}
