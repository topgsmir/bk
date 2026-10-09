//go:build linux

package externaltunnel

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestKernelSetupFailureRollsBackOnlyItsOwnTunnel(t *testing.T) {
	dir := t.TempDir()
	log := filepath.Join(dir, "calls")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '" + log + "'\ncase \"$*\" in *'addr add'*) exit 1;; esac\n"
	os.WriteFile(filepath.Join(dir, "ip"), []byte(script), 0700)
	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))
	s := fixture("gre")
	if Run(context.Background(), s) == nil {
		t.Fatal("setup failure was ignored")
	}
	b, _ := os.ReadFile(log)
	if !strings.Contains(string(b), "tunnel del "+s.Interface()) || strings.Contains(string(b), "flush") {
		t.Fatal(string(b))
	}
}
func TestCoreDownloadChecksumFailureLeavesExistingCoreUntouched(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("tampered")) }))
	defer srv.Close()
	p := filepath.Join(t.TempDir(), "core")
	os.WriteFile(p, []byte("previous"), 0700)
	e := installAsset(context.Background(), Asset{URL: srv.URL, SHA256: strings.Repeat("0", 64), Member: "rgt"}, p, io.Discard)
	if e == nil {
		t.Fatal("accepted invalid checksum")
	}
	b, _ := os.ReadFile(p)
	if string(b) != "previous" {
		t.Fatal(string(b))
	}
}
func TestVerifiedCoreIsInstalledAtomically(t *testing.T) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	f, _ := z.Create("rgt")
	f.Write([]byte{127, 'E', 'L', 'F', 0, 1})
	z.Close()
	hash := sha256.Sum256(b.Bytes())
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(b.Bytes()) }))
	defer srv.Close()
	p := filepath.Join(t.TempDir(), "core")
	if e := installAsset(context.Background(), Asset{URL: srv.URL, SHA256: hex.EncodeToString(hash[:]), Member: "rgt"}, p, io.Discard); e != nil {
		t.Fatal(e)
	}
	st, e := os.Stat(p)
	if e != nil || st.Mode().Perm() != 0700 {
		t.Fatal(st, e)
	}
}
func TestTCPForwarderCarriesByteIdenticalTrafficAndStops(t *testing.T) {
	echo, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer echo.Close()
	go func() {
		for {
			c, e := echo.Accept()
			if e != nil {
				return
			}
			go func() { defer c.Close(); io.Copy(c, c) }()
		}
	}()
	probe, _ := net.Listen("tcp", "127.0.0.1:0")
	listen := probe.Addr().String()
	probe.Close()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- serveTCP(ctx, listen, func(ctx context.Context) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", echo.Addr().String())
		})
	}()
	var c net.Conn
	for i := 0; i < 100; i++ {
		c, e = net.Dial("tcp", listen)
		if e == nil {
			break
		}
		time.Sleep(time.Millisecond * 5)
	}
	if e != nil {
		cancel()
		t.Fatal(e)
	}
	defer c.Close()
	payload := bytes.Repeat([]byte("packet"), 10000)
	go c.Write(payload)
	received := make([]byte, len(payload))
	c.SetReadDeadline(time.Now().Add(time.Second))
	if _, e = io.ReadFull(c, received); e != nil || !bytes.Equal(payload, received) {
		t.Fatal(e)
	}
	cancel()
	select {
	case e = <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("forwarder did not stop")
	}
}
func TestUnitsRunOwnAdditionalConfig(t *testing.T) {
	s := fixture("gre")
	body := Unit(s, "/etc/bk/external/sample.json", "/usr/local/bin/bk")
	if !strings.Contains(body, "external run -c") || !strings.Contains(body, "KillMode=control-group") || strings.Contains(body, "/etc/rc.local") {
		t.Fatal(body)
	}
}

func TestCoreSecretsAreRedactedAcrossOutputChunks(t *testing.T) {
	r, w := io.Pipe()
	secret := strings.Repeat("a1", 32)
	go func() {
		defer w.Close()
		w.Write([]byte("core -k " + secret[:13]))
		w.Write([]byte(secret[13:] + " ready\n" + secret))
	}()
	var out bytes.Buffer
	redactLines(&out, r, []string{secret})
	r.Close()
	if strings.Contains(out.String(), secret) || strings.Count(out.String(), "[redacted]") != 2 {
		t.Fatal(out.String())
	}
}
func TestPaqetRejectsIncompatibleParallelFixedPorts(t *testing.T) {
	s := fixture("paqet")
	if s.Connections != 1 || s.Validate() != nil {
		t.Fatal("default Paqet configuration cannot run")
	}
	s.Connections = 4
	if s.Validate() == nil {
		t.Fatal("accepted the core's invalid connection/port combination")
	}
}
