package externaltunnel

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
)

// The supplied Solarpass 2.3.0 ELF has two Spoof defects. This narrowly scoped
// repair changes only Iran's return-socket bind to frontend port + 1. The peer
// must return to that port. Originals are never changed in place.
const solarpassOriginalSHA256 = "26360c3f2f5debe6b34b0033e3312b0c4ab85de85832bd1f32b5f7071ff6416d"
const solarpassRepairedSHA256 = "5d6f510a7eae4a9be90e0aa5fcfc0066032025fbd2024da201befa69f82109d9"

func repairSolarpass(core []byte) ([]byte, error) {
	hash := sha256.Sum256(core)
	if hex.EncodeToString(hash[:]) != solarpassOriginalSHA256 {
		return nil, fmt.Errorf("Solarpass repair requires the exact pinned 2.3.0 ELF; checksum differs")
	}
	const cave, site = 0x19c833, 0x19e8fa
	beforeCave := []byte{0x66, 0x2e, 0x0f, 0x1f, 0x84, 0, 0, 0}
	beforeCall := []byte{0xe8, 0xf1, 0xd6, 0x3b, 0}
	if len(core) < site+5 || !bytes.Equal(core[cave:cave+8], beforeCave) || !bytes.Equal(core[site:site+5], beforeCall) {
		return nil, fmt.Errorf("Solarpass repair instruction precondition failed")
	}
	repaired := append([]byte(nil), core...)
	// lea esi,[rsi+1]; jmp original UDP-return-socket helper.
	copy(repaired[cave:cave+4], []byte{0x8d, 0x76, 0x01, 0xe9})
	binary.LittleEndian.PutUint32(repaired[cave+4:cave+8], uint32(0x55bff0-(cave+8)))
	displacement := int32(cave - (site + 5))
	repaired[site] = 0xe8
	binary.LittleEndian.PutUint32(repaired[site+1:site+5], uint32(displacement))
	hash = sha256.Sum256(repaired)
	if hex.EncodeToString(hash[:]) != solarpassRepairedSHA256 {
		return nil, fmt.Errorf("Solarpass repair result checksum differs")
	}
	return repaired, nil
}

// A second native defect compares the host-order sockaddr IP with a filter
// stored in native byte order. Reverse only this filter field on the pinned
// little-endian amd64 core. Actual endpoint, source and route IPs remain normal.
func solarpassPeerFilter(ip string) (string, error) {
	address := net.ParseIP(ip).To4()
	if address == nil {
		return "", fmt.Errorf("Solarpass Spoof filter needs an IPv4 literal")
	}
	return net.IPv4(address[3], address[2], address[1], address[0]).String(), nil
}
