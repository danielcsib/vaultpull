package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/sync"
	"github.com/example/vaultpull/internal/vault"
)

const version = "0.1.0"

func main() {
	var (
		configPath  string
		showVersion bool
		dryRun      bool
	)

	flag.StringVar(&configPath, "config", "vaultpull.yaml", "path to config file")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&dryRun, "dry-run", false, "print what would be written without writing files")
	flag.Parse()

	if showVersion {
		fmt.Printf("vaultpull v%s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating vault client: %v\n", err)
		os.Exit(1)
	}

	syncer := sync.New(client, dryRun)
	report, err := syncer.Run(cfg.Mappings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync failed: %v\n", err)
		os.Exit(1)
	}

	sync.PrintReport(report)

	if sync.HasErrors(report) {
		os.Exit(2)
	}
}
