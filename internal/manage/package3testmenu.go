package manage

import (
	"fmt"
	"strings"

	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
)

// Count TCP and UDP separately so even a group of dual-protocol methods
// produces at most ten measured rows, while the runner keeps four-case waves.
const package3TestGroupLimit = 10

func package3TestCaseCount(kind string) int {
	if externaltunnel.BothProtocols(kind) && !externaltunnel.UDPOnly(kind) {
		return 2
	}
	return 1
}

func package3TestGroups(catalog []externaltunnel.Kind, spoof bool) [][]externaltunnel.Kind {
	var groups [][]externaltunnel.Kind
	count := 0
	for _, k := range catalog {
		isSpoof := externaltunnel.DaggerSpoof(k.ID) || k.ID == "s3-spoof"
		if isSpoof != spoof {
			continue
		}
		cost := package3TestCaseCount(k.ID)
		if len(groups) == 0 || count+cost > package3TestGroupLimit {
			groups = append(groups, nil)
			count = 0
		}
		last := len(groups) - 1
		groups[last] = append(groups[last], k)
		count += cost
	}
	return groups
}

// There is deliberately no full-matrix shortcut: selecting a normal group
// cannot require spoof addresses or prepare another family's dependencies.
func choosePackage3TestCatalog(catalog []externaltunnel.Kind, choose func(string, []tui.Option) int) ([]externaltunnel.Kind, bool) {
	families := package3Families(catalog)
	names := []string{"Dagger", "Solarpass", "Backhaul", "Eylan VPN methods"}
	opts := make([]tui.Option, len(families))
	for i, family := range families {
		opts[i] = tui.Option{Title: names[i], Desc: fmt.Sprintf("%d normal groups; %d separate IP Spoof groups", len(package3TestGroups(family, false)), len(package3TestGroups(family, true)))}
	}
	index := choose("Package 3 Test Family", opts)
	if index < 0 || index >= len(families) {
		return nil, false
	}
	family := families[index]
	spoof := false
	if len(package3TestGroups(family, true)) > 0 {
		pick := choose(names[index]+" Test Category", []tui.Option{
			{Title: "Normal Tests", Desc: "no IP spoofing; no spoof addresses required"},
			{Title: "IP Spoof Tests (Separate)", Desc: "only choose this when your network supports spoofing"},
		})
		if pick < 0 || pick > 1 {
			return nil, false
		}
		spoof = pick == 1
	}
	groups := package3TestGroups(family, spoof)
	opts = make([]tui.Option, len(groups))
	for i, group := range groups {
		var titles []string
		cases := 0
		for _, k := range group {
			titles = append(titles, k.Title)
			cases += package3TestCaseCount(k.ID)
		}
		opts[i] = tui.Option{Title: fmt.Sprintf("Group %d (%d TCP/UDP Tests)", i+1, cases), Desc: strings.Join(titles, "; ")}
	}
	pick := choose("Select One Test Group (Maximum 10 Tests)", opts)
	if pick < 0 || pick >= len(groups) {
		return nil, false
	}
	group := groups[pick]
	opts = []tui.Option{{Title: "Test This Group", Desc: "only the listed methods; at most four test cases run together"}}
	for _, k := range group {
		opts = append(opts, tui.Option{Title: k.Title, Desc: k.Requirement})
	}
	pick = choose("Test This Group Or One Method", opts)
	if pick < 0 || pick > len(group) {
		return nil, false
	}
	if pick > 0 {
		return group[pick-1 : pick], true
	}
	return group, true
}
