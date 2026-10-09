package externaltunnel

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"os"
	"strconv"
	"time"
)

// Prefer the dedicated key, but recognize usable existing OpenSSH identities.
// Encrypted keys require the interactive dedicated-key preparation, not a stored password.
func DefaultSSHKey() string {
	for _, path := range []string{"/root/.ssh/bk_external", "/root/.ssh/id_ed25519", "/root/.ssh/id_rsa", "/root/.ssh/id_ecdsa"} {
		if b, e := os.ReadFile(path); e == nil {
			if _, e = ssh.ParsePrivateKey(b); e == nil {
				return path
			}
		}
	}
	return "/root/.ssh/bk_external"
}
func openSSH(ctx context.Context, s Spec) (*ssh.Client, net.Conn, error) {
	b, e := os.ReadFile(s.SSHKey)
	if e != nil {
		return nil, nil, fmt.Errorf("SSH key %s: %w", s.SSHKey, e)
	}
	key, e := ssh.ParsePrivateKey(b)
	if e != nil {
		return nil, nil, fmt.Errorf("SSH key: %w; prepare a dedicated key through bk", e)
	}
	check, e := knownhosts.New(s.KnownHosts)
	if e != nil {
		return nil, nil, fmt.Errorf("SSH known hosts: %w", e)
	}
	address := net.JoinHostPort(s.PeerIP, strconv.Itoa(s.Port))
	cfg := &ssh.ClientConfig{User: s.SSHUser, Auth: []ssh.AuthMethod{ssh.PublicKeys(key)}, HostKeyCallback: check}
	raw, e := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, "tcp", address)
	if e != nil {
		return nil, nil, fmt.Errorf("SSH connection: %w", e)
	}
	_ = raw.SetDeadline(time.Now().Add(10 * time.Second))
	stopped := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			raw.Close()
		case <-stopped:
		}
	}()
	cc, ch, req, e := ssh.NewClientConn(raw, address, cfg)
	close(stopped)
	if e != nil {
		raw.Close()
		return nil, nil, fmt.Errorf("SSH authentication/host verification: %w", e)
	}
	_ = raw.SetDeadline(time.Time{})
	return ssh.NewClient(cc, ch, req), raw, nil
}
func CheckSSHAuthentication(ctx context.Context, s Spec) error {
	client, _, e := openSSH(ctx, s)
	if e != nil {
		return e
	}
	return client.Close()
}
