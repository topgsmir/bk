package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/topgsmir/bk/internal/externaltunnel"
)

func runExternal(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: bk external run -c /etc/bk/external/name.json")
		os.Exit(2)
	}
	fs := flag.NewFlagSet("external "+args[0], flag.ExitOnError)
	cfg := fs.String("c", "", "additional-tunnel JSON config")
	bridge := fs.String("bridge", "", "private Alghadir packet bridge")
	kind := fs.String("kind", "", "additional tunnel kind for dependency installation")
	_ = fs.Parse(args[1:])
	if args[0] == "install" {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if _, ok := externaltunnel.Find(*kind); !ok {
			fmt.Fprintln(os.Stderr, "unknown tunnel kind")
			os.Exit(2)
		}
		if err := externaltunnel.InstallDependencies(ctx, *kind, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	s, e := externaltunnel.Load(*cfg)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	switch args[0] {
	case "run":
		e = externaltunnel.Run(ctx, s)
	case "check":
		e = externaltunnel.Check(s)
	case "alghadir-child":
		e = externaltunnel.AlghadirChild(ctx, s, *bridge)
	default:
		e = fmt.Errorf("unknown external command %q", args[0])
	}
	if e != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
