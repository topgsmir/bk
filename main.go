// bk is a launcher. Core and addons are separate executables and Go modules.
package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

//go:embed payload
var payload embed.FS

//go:embed uninstall.sh
var uninstallScript string

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			os.Exit(exit.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run(args []string, in io.Reader, out, errOut io.Writer) error {
	if len(args) > 0 {
		if args[0] == "distribution-uninstall" || args[0] == "uninstall" {
			command := exec.Command("bash", append([]string{"-s", "--"}, args[1:]...)...)
			command.Stdin = strings.NewReader(uninstallScript)
			command.Stdout, command.Stderr = out, errOut
			return command.Run()
		}
		if args[0] == "components" || args[0] == "-v" || args[0] == "--version" {
			raw, err := fs.ReadFile(payload, "payload/manifest.json")
			if err != nil {
				return fmt.Errorf("component bundle is missing; build with make build: %w", err)
			}
			if args[0] == "components" {
				_, err = out.Write(raw)
				return err
			}
			var manifest struct {
				Version string `json:"distribution_version"`
			}
			if err = json.Unmarshal(raw, &manifest); err != nil {
				return err
			}
			_, err = fmt.Fprintln(out, manifest.Version+"\nhttps://github.com/topgsmir/bk")
			return err
		}
		component, forward := route(args)
		return executeComponent(component, forward, in, out, errOut, true)
	}
	launcherPath, _ := os.Executable()
	for {
		fmt.Fprintln(out, "\nbk — independent tunnel components\n1) BackPack core\n2) Package 2\n3) Package 3\n4) Manage Package 2 and 3\n0) Exit")
		fmt.Fprint(out, "Select an option: ")
		choice, err := readLine(in)
		if errors.Is(err, io.EOF) && choice == "" {
			return nil
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		var component string
		var forwarded []string
		switch strings.TrimSpace(choice) {
		case "0":
			return nil
		case "1":
			component = "bk-core"
		case "2":
			component = "bk-addons"
			forwarded = []string{"package2"}
		case "3":
			component = "bk-addons"
			forwarded = []string{"package3"}
		case "4":
			component = "bk-addons"
			forwarded = []string{"manage"}
		default:
			fmt.Fprintln(out, "Choose 0, 1, 2, 3 or 4.")
			continue
		}
		if err = execute(component, forwarded, in, out, errOut); err != nil {
			fmt.Fprintln(errOut, err)
		}
		// A confirmed distribution uninstall removes this installed launcher.
		// Do not redraw a menu that could recreate its just-removed components.
		if launcherPath != "" {
			if _, err = os.Stat(launcherPath); errors.Is(err, os.ErrNotExist) {
				return nil
			}
		}
	}
}
func route(args []string) (string, []string) {
	if len(args) == 0 {
		return "bk-core", nil
	}
	switch args[0] {
	case "core":
		return "bk-core", args[1:]
	case "addons":
		return "bk-addons", args[1:]
	case "external", "package2", "package3":
		return "bk-addons", args
	case "manage-addons":
		return "bk-addons", append([]string{"manage"}, args[1:]...)
	}
	return "bk-core", args
}

// Read one byte at a time: buffered read-ahead would consume the answers meant
// for the selected component when stdin is a pipe or scripted terminal input.
func readLine(in io.Reader) (string, error) {
	var line strings.Builder
	var b [1]byte
	for {
		n, err := in.Read(b[:])
		if n > 0 {
			if b[0] == '\n' {
				return strings.TrimSuffix(line.String(), "\r"), nil
			}
			line.WriteByte(b[0])
			if line.Len() > 4096 {
				return "", errors.New("menu input is too long")
			}
		}
		if err != nil {
			return line.String(), err
		}
	}
}
func execute(name string, args []string, in io.Reader, out, errOut io.Writer) error {
	return executeComponent(name, args, in, out, errOut, false)
}
func executeComponent(name string, args []string, in io.Reader, out, errOut io.Writer, replace bool) error {
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	path, err := materialize(payload, filepath.Join(cache, "bk", "components"), name)
	if err != nil {
		return err
	}
	// Service commands keep their PID and signal semantics and retain no
	// additional launcher runtime after selecting their owning component.
	if replace && in == os.Stdin && out == os.Stdout && errOut == os.Stderr {
		return syscall.Exec(path, append([]string{name}, args...), os.Environ())
	}
	command := exec.Command(path, args...)
	command.Stdin = in
	command.Stdout = out
	command.Stderr = errOut
	if err = command.Start(); err != nil {
		return err
	}
	signals := make(chan os.Signal, 4)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	done := make(chan struct{})
	defer close(done)
	defer signal.Stop(signals)
	go func() {
		for {
			select {
			case sig := <-signals:
				_ = command.Process.Signal(sig)
			case <-done:
				return
			}
		}
	}()
	return command.Wait()
}
func materialize(source fs.FS, cache, name string) (string, error) {
	if name != "bk-core" && name != "bk-addons" {
		return "", errors.New("unknown component")
	}
	component, err := source.Open("payload/" + name)
	if err != nil {
		return "", fmt.Errorf("%s is absent; build with make build: %w", name, err)
	}
	hash, err := digestReader(component)
	component.Close()
	if err != nil {
		return "", err
	}
	manifestBytes, err := fs.ReadFile(source, "payload/manifest.json")
	if err != nil {
		return "", err
	}
	manifest := map[string]any{}
	if err = json.Unmarshal(manifestBytes, &manifest); err != nil {
		return "", err
	}
	if manifest[name+"_sha256"] != hash {
		return "", errors.New("embedded component checksum mismatch")
	}
	directory := filepath.Join(cache, hash)
	if err = os.MkdirAll(directory, 0700); err != nil {
		return "", err
	}
	for current := directory; ; current = filepath.Dir(current) {
		info, e := os.Lstat(current)
		if e != nil {
			return "", e
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("component cache must use real directories")
		}
		if current == cache {
			break
		}
	}
	target := filepath.Join(directory, name)
	if info, e := os.Lstat(target); e == nil && info.Mode().IsRegular() {
		if cached, e := os.Open(target); e == nil {
			actual, digestErr := digestReader(cached)
			cached.Close()
			if digestErr == nil && actual == hash {
				if e = os.Chmod(target, 0700); e != nil {
					return "", e
				}
				return target, nil
			}
		}
	}
	temporary, err := os.CreateTemp(directory, "component-")
	if err != nil {
		return "", err
	}
	defer os.Remove(temporary.Name())
	component, err = source.Open("payload/" + name)
	if err != nil {
		temporary.Close()
		return "", err
	}
	written := sha256.New()
	_, err = io.Copy(io.MultiWriter(temporary, written), component)
	component.Close()
	if err != nil {
		temporary.Close()
		return "", err
	}
	if hex.EncodeToString(written.Sum(nil)) != hash {
		temporary.Close()
		return "", errors.New("copied component checksum mismatch")
	}
	if err = temporary.Chmod(0700); err != nil {
		temporary.Close()
		return "", err
	}
	if err = temporary.Sync(); err != nil {
		temporary.Close()
		return "", err
	}
	if err = temporary.Close(); err != nil {
		return "", err
	}
	if err = os.Rename(temporary.Name(), target); err != nil {
		return "", err
	}
	return target, nil
}

// Components are streamed: starting many services must not allocate one whole
// embedded executable per launcher process.
func digestReader(source io.Reader) (string, error) {
	digest := sha256.New()
	if _, err := io.Copy(digest, source); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}
