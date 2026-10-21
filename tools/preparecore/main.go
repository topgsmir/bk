package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := prepare(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func prepare() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	source := filepath.Join(root, "upstream", "BackPack")
	var hashes map[string]string
	if err = readJSON(filepath.Join(root, "integration", "upstream-files.json"), &hashes); err != nil {
		return err
	}
	var keep map[string][]string
	if err = readJSON(filepath.Join(root, "integration", "normalization.json"), &keep); err != nil {
		return err
	}
	var modes map[string]string
	if err = readJSON(filepath.Join(root, "integration", "upstream-modes.json"), &modes); err != nil {
		return err
	}
	build := filepath.Join(root, ".build")
	if err = os.MkdirAll(build, 0700); err != nil {
		return err
	}
	next, err := os.MkdirTemp(build, "core-next-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(next)
	count := 0
	err = filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, e := filepath.Rel(source, path)
		if e != nil {
			return e
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("upstream symlink refused: %s", relative)
		}
		target := filepath.Join(next, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		raw, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		digest := sha256.Sum256(raw)
		name := filepath.ToSlash(relative)
		if hashes[name] != hex.EncodeToString(digest[:]) {
			return fmt.Errorf("pristine upstream changed: %s", name)
		}
		count++
		text := string(raw)
		if !strings.ContainsRune(text, 0) {
			for i, literal := range keep[name] {
				text = strings.ReplaceAll(text, literal, fmt.Sprintf("__BK_PRESERVE_%d__", i))
			}
			text = strings.NewReplacer("github.com/backpack/backpack", "github.com/topgsmir/bk", "AminMGMT/BackPack", "topgsmir/bk", "BACKPACK", "BK", "BackPack", "bk", "Backpack", "bk", "backpack", "bk").Replace(text)
			for i, literal := range keep[name] {
				text = strings.ReplaceAll(text, fmt.Sprintf("__BK_PRESERVE_%d__", i), literal)
			}
			raw = []byte(text)
		}
		var mode fs.FileMode
		switch modes[name] {
		case "100644":
			mode = 0644
		case "100755":
			mode = 0755
		default:
			return fmt.Errorf("unsupported upstream mode for %s: %s", name, modes[name])
		}
		return os.WriteFile(target, raw, mode)
	})
	if err != nil {
		return err
	}
	if count != len(hashes) || count != len(modes) {
		return fmt.Errorf("upstream manifest has %d files, found %d", len(hashes), count)
	}
	if err = run(next, "gofmt", "-w", "."); err != nil {
		return err
	}
	patch := filepath.Join(root, "integration", "core.patch")
	if err = run(next, "git", "apply", "--check", patch); err != nil {
		return fmt.Errorf("core overlay conflicts; current generated core has been kept: %w", err)
	}
	if err = run(next, "git", "apply", patch); err != nil {
		return err
	}
	if err = run(next, "gofmt", "-w", "."); err != nil {
		return err
	}
	target := filepath.Join(build, "core")
	// Only this disposable directory inside the current workspace may be removed.
	absolute, e := filepath.Abs(target)
	if e != nil {
		return e
	}
	relative, e := filepath.Rel(build, absolute)
	if e != nil || relative != "core" {
		return fmt.Errorf("unsafe generated-core target: %s", target)
	}
	if err = os.RemoveAll(absolute); err != nil {
		return err
	}
	if err = os.Rename(next, absolute); err != nil {
		return err
	}
	fmt.Println("Prepared isolated core at", target)
	return nil
}
func readJSON(path string, value any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, value)
}
func run(dir, name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Dir = dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
