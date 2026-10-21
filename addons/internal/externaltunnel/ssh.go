package externaltunnel

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
)

func runSSH(ctx context.Context, s Spec) error {
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
		client, _, e := openSSH(ctx, s)
		if e != nil {
			return nil, e
		}
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
		fmt.Fprintln(os.Stderr, "SSH forwarding failed (check the target and AllowTcpForwarding/PermitOpen):", last)
		return nil, last
	})
}
