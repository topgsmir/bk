package webui

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

// Edit on a direct tunnel used to open the reverse form and post its fields at
// the top of the request, where the direct edit reads nothing — so it showed
// the wrong dialog and could not save. The direct form posts under "direct",
// and every key it can send has to be one the server reads: one unknown key
// fails the whole save.
func TestTheDirectEditFormPostsWhatTheServerReads(t *testing.T) {
	src, err := os.ReadFile("panel/js/views/edit.js")
	if err != nil {
		t.Fatal(err)
	}
	js := string(src)
	if !strings.Contains(js, "direct: readDirect(root, settings)") {
		t.Fatal("the direct form no longer posts under \"direct\"")
	}
	src, err = os.ReadFile("panel/js/ui/directedit.js")
	if err != nil {
		t.Fatal(err)
	}
	js = string(src)
	start := strings.Index(js, "function directMarkup(")
	end := strings.Index(js, "const DIRECT_NUMBERS")
	if start < 0 || end < start {
		t.Fatal("directedit.js no longer has the direct form")
	}
	markup := js[start:end]
	var keys []string
	for _, m := range regexp.MustCompile(`(?:field|sw)\('([A-Za-z.]+)'`).FindAllStringSubmatch(markup, -1) {
		keys = append(keys, m[1])
	}
	keys = append(keys, "preset") // the menu
	if len(keys) < 8 {
		t.Fatalf("found only %v in the direct form", keys)
	}

	direct := map[string]any{}
	spoof := map[string]any{}
	for _, k := range keys {
		var v any = 1
		switch k {
		case "ports", "preset", "spoof.srcIPs", "spoof.peerIP", "spoof.peerSrcIP", "spoof.interface":
			v = "x"
		case "acceptUdp", "autoMtu", "fec", "stealth", "spoof.shufflePort":
			v = true
		}
		if strings.HasPrefix(k, "spoof.") {
			spoof[strings.TrimPrefix(k, "spoof.")] = v
		} else {
			direct[k] = v
		}
	}
	direct["spoof"] = spoof
	raw, err := json.Marshal(map[string]any{"name": "t", "direct": direct})
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	var req tunnelEditRequest
	r := httptest.NewRequest("POST", "/api/tunnel/edit", strings.NewReader(body))
	if err := decodeJSON(httptest.NewRecorder(), r, &req); err != nil {
		t.Fatalf("the server refuses what the direct form sends: %v\n%s", err, body)
	}
	if req.Direct.MTU == nil || req.Direct.Ports == nil || req.Direct.FEC == nil {
		t.Errorf("the direct fields did not arrive: %+v", req.Direct)
	}
	if req.Direct.Spoof == nil || req.Direct.Spoof.SrcIPs != "x" || req.Direct.Spoof.PeerIP != "x" || req.Direct.Spoof.PeerSrcIP != "x" || !req.Direct.Spoof.ShufflePort {
		t.Fatal("the visible spoof address fields did not reach the nested carrier edit")
	}
}
