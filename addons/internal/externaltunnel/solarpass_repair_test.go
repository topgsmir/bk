package externaltunnel

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestSolarpassSpoofPeerFilter(t *testing.T) {
	got, err := solarpassPeerFilter("10.202.0.1")
	if err != nil || got != "1.0.202.10" {
		t.Fatal(got, err)
	}
	back, err := solarpassPeerFilter(got)
	if err != nil || back != "10.202.0.1" {
		t.Fatal(back, err)
	}
	for _, value := range []string{"", "example.test", "::1", "999.0.0.1"} {
		if _, err := solarpassPeerFilter(value); err == nil {
			t.Fatal("accepted invalid IPv4 filter", value)
		}
	}
}
func TestSolarpassSpoofRepairRefusesUnpinnedELF(t *testing.T) {
	if _, err := repairSolarpass([]byte("unexpected executable")); err == nil {
		t.Fatal("repair accepted unknown bytes")
	}
}
func TestSolarpassSpoofRepairExactProvidedCore(t *testing.T) {
	path := os.Getenv("BK_SOLARPASS_ORIGINAL")
	if path == "" {
		t.Skip("set BK_SOLARPASS_ORIGINAL for the supplied-core repair verification")
	}
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	preserved := append([]byte(nil), original...)
	repaired, err := repairSolarpass(original)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, preserved) {
		t.Fatal("original ELF changed")
	}
	hash := sha256.Sum256(repaired)
	if hex.EncodeToString(hash[:]) != solarpassRepairedSHA256 {
		t.Fatal("wrong repaired core")
	}
	// No byte outside the two verified instruction ranges may differ.
	for i := range original {
		if original[i] != repaired[i] && !(i >= 0x19c833 && i < 0x19c83b) && !(i >= 0x19e8fa && i < 0x19e8ff) {
			t.Fatalf("unexpected patch at %#x", i)
		}
	}
	altered := append([]byte(nil), original...)
	altered[0] ^= 1
	if _, err := repairSolarpass(altered); err == nil {
		t.Fatal("accepted tampered input")
	}
	if _, err := repairSolarpass(repaired); err == nil {
		t.Fatal("accepted a previously patched input")
	}
}
