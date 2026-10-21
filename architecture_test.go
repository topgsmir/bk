package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestUpstreamSnapshotIsExactlyThePinnedPristineSource(t *testing.T) {
	raw, err := os.ReadFile("integration/upstream-files.json")
	if err != nil {
		t.Fatal(err)
	}
	hashes := map[string]string{}
	if err = json.Unmarshal(raw, &hashes); err != nil {
		t.Fatal(err)
	}
	count := 0
	err = filepath.WalkDir("upstream/BackPack", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			t.Errorf("upstream symlink: %s", path)
			return nil
		}
		relative, _ := filepath.Rel("upstream/BackPack", path)
		data, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		digest := sha256.Sum256(data)
		if hashes[filepath.ToSlash(relative)] != hex.EncodeToString(digest[:]) {
			t.Errorf("upstream file was edited: %s", relative)
		}
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != len(hashes) {
		t.Fatalf("snapshot file count %d != manifest %d", count, len(hashes))
	}
}
func TestAddonsHaveNoBuildDependencyOnTheBackPackCore(t *testing.T) {
	raw, err := os.ReadFile("addons/go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "replace ") {
		t.Fatal("addon module must not replace or import the core module")
	}
	err = filepath.WalkDir("addons", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, e := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if e != nil {
			return e
		}
		for _, spec := range file.Imports {
			value, _ := strconv.Unquote(spec.Path.Value)
			if strings.HasPrefix(value, "github.com/backpack/backpack") || strings.HasPrefix(value, "github.com/topgsmir/bk/") && !strings.HasPrefix(value, "github.com/topgsmir/bk/addons/") {
				t.Errorf("addon imports core implementation: %s: %s", path, value)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
