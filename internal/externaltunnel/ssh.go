package externaltunnel

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func runSSH(ctx context.Context, s Spec) error {
	b, e := os.ReadFile(s.SSHKey)
	if e != nil {
		return fmt.Errorf("SSH key: %w", e)
	}
	key, e := ssh.ParsePrivateKey(b)
	if e != nil {
		return e
	}
	hostCheck, e := knownhosts.New(s.KnownHosts)
	if e != nil {
		return fmt.Errorf("SSH known hosts: %w; verify the peer with ssh before setup", e)
	}
	cfg := &ssh.ClientConfig{User: s.SSHUser, Auth: []ssh.AuthMethod{ssh.PublicKeys(key)}, HostKeyCallback: hostCheck, Timeout: 10 * time.Second}
	pool := make([]*ssh.Client, s.Connections)
	var mu sync.Mutex
	var turn atomic.Uint64
	defer func() {
		mu.Lock()
		defer mu.Unlock()
		for _, c := range pool {
			if c != nil {
				c.Close()
			}
		}
	}()
	dialSSH := func(i int) (*ssh.Client, error) {
		mu.Lock()
		defer mu.Unlock()
		if pool[i] != nil {
			return pool[i], nil
		}
		d := net.Dialer{Timeout: 10 * time.Second}
		raw, e := d.DialContext(ctx, "tcp", net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)))
		if e != nil {
			return nil, e
		}
		_ = raw.SetDeadline(time.Now().Add(10 * time.Second))
		cc, ch, req, e := ssh.NewClientConn(raw, net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port)), cfg)
		if e != nil {
			raw.Close()
			return nil, e
		}
		_ = raw.SetDeadline(time.Time{})
		client := ssh.NewClient(cc, ch, req)
		pool[i] = client
		go func() {
			_ = client.Wait()
			mu.Lock()
			if pool[i] == client {
				pool[i] = nil
			}
			mu.Unlock()
		}()
		return client, nil
	}
	// Establish at least one authenticated connection before reporting readiness.
	if _, e := dialSSH(0); e != nil {
		return e
	}
	stop := make(chan struct{})
	maintenanceDone := make(chan struct{})
	defer func() { close(stop); <-maintenanceDone }()
	go func() {
		defer close(maintenanceDone)
		tick := time.NewTicker(20 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			case <-tick.C:
				mu.Lock()
				snapshot := append([]*ssh.Client(nil), pool...)
				mu.Unlock()
				var checks sync.WaitGroup
				for i, c := range snapshot {
					if c == nil {
						continue
					}
					checks.Add(1)
					go func(i int, c *ssh.Client) {
						defer checks.Done()
						done := make(chan error, 1)
						go func() { _, _, e := c.SendRequest("keepalive@openssh.com", true, nil); done <- e }()
						failed := false
						select {
						case e := <-done:
							failed = e != nil
						case <-ctx.Done():
							failed = true
						case <-stop:
							failed = true
						case <-time.After(5 * time.Second):
							failed = true
						}
						if failed {
							c.Close()
							mu.Lock()
							if pool[i] == c {
								pool[i] = nil
							}
							mu.Unlock()
						}
					}(i, c)
				}
				checks.Wait()
			}
		}
	}()
	return serveTCP(ctx, s.Listen, func(ctx context.Context) (net.Conn, error) {
		var last error
		for n := 0; n < len(pool); n++ {
			i := int(turn.Add(1) % uint64(len(pool)))
			c, e := dialSSH(i)
			if e != nil {
				last = e
				continue
			}
			conn, e := c.DialContext(ctx, "tcp", s.Target)
			if e == nil {
				return conn, nil
			}
			last = e
			// Do not tear down healthy streams merely because the target refused a new one.
		}
		return nil, last
	})
}
