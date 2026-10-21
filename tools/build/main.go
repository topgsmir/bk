package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	out := flag.String("o", "bk", "output bundle")
	prepareOnly := flag.Bool("prepare-only", false, "prepare core source only")
	goos := flag.String("goos", "linux", "target operating system")
	goarch := flag.String("goarch", runtime.GOARCH, "target architecture")
	goarm := flag.String("goarm", "", "ARM variant")
	flag.Parse()
	if err := build(*out, *prepareOnly, *goos, *goarch, *goarm); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func build(out string, prepareOnly bool, goos, goarch, goarm string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	// Build tools must run on the host, even for a cross-compiled component.
	host := append(os.Environ(), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	command := exec.Command("go", "run", "./tools/preparecore")
	command.Env = host
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err = command.Run(); err != nil {
		return err
	}
	if prepareOnly {
		return nil
	}
	if err = os.MkdirAll("payload", 0755); err != nil {
		return err
	}
	var lock map[string]string
	raw, err := os.ReadFile("upstream.lock.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(raw, &lock); err != nil {
		return err
	}
	for path, field := range map[string]string{"VERSION": "distribution_version", "upstream/BackPack/VERSION": "version", "addons/VERSION": "addons_version"} {
		bytes, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		if strings.TrimSpace(string(bytes)) != lock[field] {
			return fmt.Errorf("%s and upstream.lock.json %s disagree", path, field)
		}
	}
	manifest := map[string]any{"distribution_version": lock["distribution_version"], "upstream_commit": lock["commit"], "core_version": lock["version"], "addons_version": lock["addons_version"]}
	targetEnv := append(os.Environ(), "CGO_ENABLED=0", "GOOS="+goos, "GOARCH="+goarch, "GOARM="+goarm)
	var sbom []byte
	for name, dir := range map[string]string{"bk-core": ".build/core", "bk-addons": "addons"} {
		target := filepath.Join(root, "payload", name)
		args := []string{"build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w", "-o", target, "."}
		if name == "bk-core" && goarm != "" {
			args[3] = "-ldflags=-s -w -X github.com/topgsmir/bk/internal/app.GOARM=" + goarm
		}
		command = exec.Command("go", args...)
		command.Dir = dir
		command.Env = targetEnv
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err = command.Run(); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		bytes, e := os.ReadFile(target)
		if e != nil {
			return e
		}
		digest := sha256.Sum256(bytes)
		manifest[name+"_sha256"] = hex.EncodeToString(digest[:])
		metadata := exec.Command("go", "version", "-m", target)
		metadata.Env = host
		info, e := metadata.CombinedOutput()
		if e != nil {
			return fmt.Errorf("component SBOM: %w: %s", e, info)
		}
		sbom = append(sbom, info...)
	}
	if err = os.MkdirAll(".build/sbom", 0755); err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(".build/sbom", goos+"-"+goarch+goarm+".txt"), sbom, 0644); err != nil {
		return err
	}
	raw, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile("payload/manifest.json", append(raw, '\n'), 0644); err != nil {
		return err
	}
	command = exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w", "-o", out, ".")
	command.Env = targetEnv
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
