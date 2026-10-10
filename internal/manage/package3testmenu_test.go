package manage

import (
	"reflect"
	"strings"
	"testing"

	"github.com/topgsmir/bk/internal/externaltunnel"
	"github.com/topgsmir/bk/internal/tui"
)

func testGroupIDs(group []externaltunnel.Kind) []string {
	var ids []string
	for _, kind := range group {
		ids = append(ids, kind.ID)
	}
	return ids
}

func TestPackage3GroupsBoundActualRowsAndKeepSpoofSeparate(t *testing.T) {
	seen := map[string]bool{}
	for _, family := range package3Families(externaltunnel.Package3Kinds) {
		for _, spoof := range []bool{false, true} {
			groups := package3TestGroups(family, spoof)
			for _, group := range groups {
				rows := 0
				if len(group) == 0 {
					t.Fatal("empty group")
				}
				for _, kind := range group {
					if seen[kind.ID] {
						t.Fatalf("duplicate method: %s", kind.ID)
					}
					seen[kind.ID] = true
					isSpoof := externaltunnel.DaggerSpoof(kind.ID) || kind.ID == "s3-spoof"
					if isSpoof != spoof {
						t.Fatalf("mixed normal and spoof group: %s", kind.ID)
					}
					rows += package3TestCaseCount(kind.ID)
				}
				if rows > 10 {
					t.Fatalf("group starts %d rows", rows)
				}
			}
		}
	}
	if len(seen) != len(externaltunnel.Package3Kinds) {
		t.Fatalf("methods lost: %d of %d", len(seen), len(externaltunnel.Package3Kinds))
	}
}

func TestPackage3MenuSelectsOnlyChosenGroupOrMethod(t *testing.T) {
	for familyIndex, family := range package3Families(externaltunnel.Package3Kinds) {
		for _, spoof := range []bool{false, true} {
			groups := package3TestGroups(family, spoof)
			for groupIndex, group := range groups {
				for methodIndex := 0; methodIndex <= len(group); methodIndex++ {
					choices := []int{familyIndex}
					if len(package3TestGroups(family, true)) > 0 {
						scope := 0
						if spoof {
							scope = 1
						}
						choices = append(choices, scope)
					}
					choices = append(choices, groupIndex, methodIndex)
					step := 0
					got, ok := choosePackage3TestCatalog(externaltunnel.Package3Kinds, func(title string, opts []tui.Option) int {
						if step >= len(choices) {
							t.Fatalf("unexpected extra prompt: %s", title)
						}
						if len(opts) > 11 {
							t.Fatalf("oversized menu: %s (%d)", title, len(opts))
						}
						for _, opt := range opts {
							if strings.Contains(opt.Title, "All") {
								t.Fatalf("unbounded selection: %s", opt.Title)
							}
						}
						pick := choices[step]
						step++
						return pick
					})
					want := group
					if methodIndex > 0 {
						want = group[methodIndex-1 : methodIndex]
					}
					if !ok || step != len(choices) || !reflect.DeepEqual(testGroupIDs(got), testGroupIDs(want)) {
						t.Fatalf("choices %v selected %v; want %v", choices, testGroupIDs(got), testGroupIDs(want))
					}
				}
			}
		}
	}
}

func TestPackage3MenuCancellationNeverStartsDefaultMatrix(t *testing.T) {
	for cancelAt := 0; cancelAt < 4; cancelAt++ {
		step := 0
		got, ok := choosePackage3TestCatalog(externaltunnel.Package3Kinds, func(_ string, _ []tui.Option) int {
			pick := 0
			if step == cancelAt {
				pick = -1
			}
			step++
			return pick
		})
		if ok || got != nil || step != cancelAt+1 {
			t.Fatalf("cancel at %d selected %v after %d prompts", cancelAt, got, step)
		}
	}
}
