package manage

import (
	"context"
	"fmt"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func externalKindIDs(cases []externaltunnel.Spec) []string {
	ids := []string{}
	seen := map[string]bool{}
	for _, s := range cases {
		if !seen[s.Kind] {
			seen[s.Kind] = true
			ids = append(ids, s.Kind)
		}
	}
	return ids
}
func selectedExternalKinds(kinds []string) []string {
	if kinds != nil {
		return kinds
	}
	for _, k := range externaltunnel.Kinds {
		kinds = append(kinds, k.ID)
	}
	return kinds
}
func ensureExternalSSH(ctx context.Context, s externaltunnel.Spec, out io.Writer) error {
	if !externaltunnel.SSHInitiator(s) {
		return nil
	}
	if e := s.Validate(); e != nil {
		return e
	}
	if e := externaltunnel.CheckSSHAuthentication(ctx, s); e == nil {
		return nil
	}
	// Passwords remain in OpenSSH's interactive prompt and are never stored in
	// configuration, logs, setup links or the test coordinator.
	for _, program := range []string{"ssh-keygen", "ssh-copy-id"} {
		if _, e := exec.LookPath(program); e != nil {
			if e = externaltunnel.InstallDependencies(ctx, s.Kind, out); e != nil {
				return e
			}
			break
		}
	}
	if e := os.MkdirAll(filepath.Dir(s.SSHKey), 0700); e != nil {
		return e
	}
	if e := os.MkdirAll(filepath.Dir(s.KnownHosts), 0700); e != nil {
		return e
	}
	if _, e := os.Lstat(s.SSHKey); os.IsNotExist(e) {
		cmd := exec.CommandContext(ctx, "ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", s.SSHKey)
		cmd.Stdout, cmd.Stderr = out, out
		if e = cmd.Run(); e != nil {
			return fmt.Errorf("SSH key creation: %w", e)
		}
	}
	if _, e := os.Stat(s.SSHKey + ".pub"); os.IsNotExist(e) {
		b, e := os.ReadFile(s.SSHKey)
		if e != nil {
			return e
		}
		key, e := ssh.ParsePrivateKey(b)
		if e != nil {
			return fmt.Errorf("SSH key: %w; use a dedicated unencrypted key", e)
		}
		if e = os.WriteFile(s.SSHKey+".pub", ssh.MarshalAuthorizedKey(key.PublicKey()), 0644); e != nil {
			return e
		}
	}
	fmt.Fprintf(out, "SSH authentication for %s@%s:%d: verify the displayed host fingerprint against that server. Enter its login password only in OpenSSH's prompt if requested.\n", s.SSHUser, s.PeerIP, s.Port)
	cmd := exec.CommandContext(ctx, "ssh-copy-id", "-p", strconv.Itoa(s.Port), "-i", s.SSHKey+".pub", "-o", "StrictHostKeyChecking=ask", "-o", "UserKnownHostsFile="+s.KnownHosts, "-o", "ConnectTimeout=10", "-o", "NumberOfPasswordPrompts=1", s.SSHUser+"@"+s.PeerIP)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, out, out
	if e := cmd.Run(); e != nil {
		return fmt.Errorf("SSH authentication preparation failed: %w; SSH login alone is not proof that port forwarding is allowed", e)
	}
	return externaltunnel.CheckSSHAuthentication(ctx, s)
}

func promptExternalSSH(kind, side, local, peer string) externaltunnel.Spec {
	s := externaltunnel.New("ssh-auth", kind, side)
	s.LocalIP, s.PeerIP = local, peer
	s.Port = externalNumber("Peer SSH Port", 22)
	s.SSHUser = tui.PromptDefault("Peer SSH User", s.SSHUser)
	s.SSHKey = tui.PromptDefault("SSH Private Key (Existing Or New Dedicated Key)", externaltunnel.DefaultSSHKey())
	s.KnownHosts = tui.PromptDefault("Verified Known Hosts File", s.KnownHosts)
	return s
}
func prepareExternalPeer(ctx context.Context, cases []externaltunnel.Spec, out io.Writer) map[string]string {
	preparation := externaltunnel.PrepareDependencies(ctx, externalKindIDs(cases), out)
	for i, spec := range cases {
		local := spec.Mirror()
		if externaltunnel.SSHInitiator(local) {
			local.SSHKey = tui.PromptDefault("SSH Reverse Private Key On Kharej", externaltunnel.DefaultSSHKey())
			local.KnownHosts = tui.PromptDefault("Verified Iran Known Hosts File", local.KnownHosts)
			cases[i].SSHKey, cases[i].KnownHosts = local.SSHKey, local.KnownHosts
			if e := ensureExternalSSH(ctx, local, out); e != nil {
				preparation[local.Kind] = e.Error()
			}
		}
	}
	return preparation
}
