package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

type boundedReadFS struct{ fs.FS }
type boundedReadFile struct{ fs.File }

func (source boundedReadFS) Open(name string) (fs.File, error) {
	file, err := source.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return boundedReadFile{file}, nil
}
func (file boundedReadFile) Read(buffer []byte) (int, error) {
	if len(buffer) > 256*1024 {
		return 0, fmt.Errorf("whole-component allocation during startup")
	}
	return file.File.Read(buffer)
}
func TestStartingAComponentDoesNotReadItsWholeExecutableIntoAMemoryBuffer(t *testing.T) {
	data := bytes.Repeat([]byte("component"), 256*1024)
	digest := sha256.Sum256(data)
	files := fstest.MapFS{"payload/bk-core": {Data: data}, "payload/manifest.json": {Data: []byte(fmt.Sprintf(`{"bk-core_sha256":%q}`, hex.EncodeToString(digest[:])))}}
	if _, err := materialize(boundedReadFS{files}, filepath.Join(t.TempDir(), "cache"), "bk-core"); err != nil {
		t.Fatal(err)
	}
}

func TestDistributionUninstallUsesTheSharedDryRun(t *testing.T) {
	// This checks launcher delegation, independently of the test host's
	// crontab and systemd state. Removal behavior has disposable-container fixtures.
	commands := t.TempDir()
	for name, script := range map[string]string{
		"crontab":   "#!/bin/sh\necho 'no crontab for fixture' >&2\nexit 1\n",
		"systemctl": "#!/bin/sh\nexit 0\n",
	} {
		if err := os.WriteFile(filepath.Join(commands, name), []byte(script), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", commands+string(os.PathListSeparator)+os.Getenv("PATH"))
	var output, errors bytes.Buffer
	if err := run([]string{"distribution-uninstall", "--dry-run"}, strings.NewReader(""), &output, &errors); err != nil {
		t.Fatalf("%v\nstdout: %s\nstderr: %s", err, output.String(), errors.String())
	}
	for _, path := range []string{"/etc/bk", "/usr/local/bin/bk", "/root/.cache/bk/components"} {
		if !strings.Contains(output.String(), path) {
			t.Fatalf("shared dry run omits %s", path)
		}
	}
}

func TestExistingServiceCommandsRouteToTheirOwningComponent(t *testing.T) {
	for _, test := range []struct {
		args      []string
		component string
		forward   []string
	}{
		{[]string{"-c", "/etc/bk/a.toml"}, "bk-core", []string{"-c", "/etc/bk/a.toml"}},
		{[]string{"--webui"}, "bk-core", []string{"--webui"}},
		{[]string{"external", "run", "-c", "/etc/bk/external/a.json"}, "bk-addons", []string{"external", "run", "-c", "/etc/bk/external/a.json"}},
		{[]string{"core", "short-https", "serve"}, "bk-core", []string{"short-https", "serve"}},
		{[]string{"package3"}, "bk-addons", []string{"package3"}},
	} {
		component, forward := route(test.args)
		if component != test.component || !reflect.DeepEqual(forward, test.forward) {
			t.Fatalf("%v -> %s %v", test.args, component, forward)
		}
	}
}
func TestMenuDoesNotConsumeTheChildProgramsInput(t *testing.T) {
	in := strings.NewReader("2\nanswer for child\n")
	line, err := readLine(in)
	if err != nil || line != "2" {
		t.Fatal(line, err)
	}
	next, err := readLine(in)
	if err != nil || next != "answer for child" {
		t.Fatal(next, err)
	}
}
func TestComponentExtractionVerifiesAndRepairsItsPrivateExecutable(t *testing.T) {
	data := []byte("component bytes")
	digest := sha256.Sum256(data)
	hash := hex.EncodeToString(digest[:])
	files := fstest.MapFS{"payload/bk-core": {Data: data}, "payload/manifest.json": {Data: []byte(fmt.Sprintf(`{"bk-core_sha256":%q}`, hash))}}
	cache := filepath.Join(t.TempDir(), "cache")
	path, err := materialize(files, cache, "bk-core")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	path, err = materialize(files, cache, "bk-core")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(data) {
		t.Fatal("corruption survived")
	}
	files["payload/bk-core"] = &fstest.MapFile{Data: []byte("changed bundle")}
	if _, err = materialize(files, cache, "bk-core"); err == nil {
		t.Fatal("bad embedded hash accepted")
	}
}
func TestComponentCacheRejectsDirectorySymlinks(t *testing.T) {
	base := t.TempDir()
	cache := filepath.Join(base, "cache")
	if err := os.Symlink(filepath.Join(base, "elsewhere"), cache); err != nil {
		t.Skip(err)
	}
	data := []byte("test")
	hash := sha256.Sum256(data)
	files := fstest.MapFS{"payload/bk-core": {Data: data}, "payload/manifest.json": {Data: []byte(fmt.Sprintf(`{"bk-core_sha256":%q}`, hex.EncodeToString(hash[:])))}}
	if _, err := materialize(files, cache, "bk-core"); err == nil {
		t.Fatal("cache symlink accepted")
	}
}
