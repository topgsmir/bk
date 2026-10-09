//go:build linux

package externaltunnel

import (
	"encoding/json"
	"net"
	"strings"
	"testing"
)

func TestAlghadirAuthenticationCreatesFreshPairedSessionKeys(t *testing.T) {
	s := fixture("alghadir")
	var previous string
	for i := 0; i < 2; i++ {
		a, b := net.Pipe()
		defer a.Close()
		defer b.Close()
		type result struct {
			k string
			e error
		}
		r := make(chan result, 2)
		go func() { k, e := sessionKey(a, s); r <- result{k, e} }()
		go func() { k, e := sessionKey(b, s.Mirror()); r <- result{k, e} }()
		one, two := <-r, <-r
		if one.e != nil || two.e != nil || one.k != two.k || one.k == previous {
			t.Fatal(one, two)
		}
		previous = one.k
	}
}
func TestAlghadirRejectsWrongSecret(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	s := fixture("alghadir")
	peer := s.Mirror()
	peer.Secret = NewSecret()
	err := make(chan error, 2)
	go func() { _, e := sessionKey(a, s); err <- e }()
	go func() { _, e := sessionKey(b, peer); err <- e }()
	if <-err == nil || <-err == nil {
		t.Fatal("wrong secret authenticated")
	}
}
func TestObfsStateAndIPsecAreDedicated(t *testing.T) {
	s := fixture("alghadir")
	state, cert := obfsState(s)
	var st map[string]any
	if e := json.Unmarshal(state, &st); e != nil || cert == "" || len(st["drbg-seed"].(string)) != 48 {
		t.Fatal(e, string(state))
	}
	joined := ""
	for _, c := range ipsecGRECommands(s) {
		joined += strings.Join(c.Args, " ") + "\n"
	}
	if !strings.Contains(joined, "proto gre") || !strings.Contains(joined, "rfc4106(gcm(aes))") || strings.Contains(joined, "flush") {
		t.Fatal(joined)
	}
}
