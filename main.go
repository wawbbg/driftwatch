package main

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	fmt.Printf("driftwatch starting\n")
	fmt.Printf("config file: %s\n", cfg.ConfigFile)
	fmt.Printf("services: %d\n", len(cfg.Services))
	return nil
}
