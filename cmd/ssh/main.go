package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	internalSSH "github.com/dev-harshhh19/harshad-tui/internal/ssh"
)

func main() {
	cfg := internalSSH.DefaultConfig()

	if portEnv := os.Getenv("SSH_PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil && p > 0 {
			cfg.Port = p
		}
	} else if portEnv := os.Getenv("PORT"); portEnv != "" {
		// If running in an environment where PORT is specified for the SSH binary
		if p, err := strconv.Atoi(portEnv); err == nil && p > 0 && p != 3000 {
			cfg.Port = p
		}
	}

	if hostEnv := os.Getenv("SSH_HOST"); hostEnv != "" {
		cfg.Host = hostEnv
	}

	if keyEnv := os.Getenv("SSH_KEY_PATH"); keyEnv != "" {
		cfg.KeyPath = keyEnv
	}

	log.SetLevel(log.InfoLevel)
	log.Info("Harshad Nikam — SSH Terminal Portfolio Service", "version", "1.0.0")
	fmt.Printf("Configured to serve SSH sessions on %s:%d\n", cfg.Host, cfg.Port)

	ctx := context.Background()
	if err := internalSSH.StartServer(ctx, cfg); err != nil {
		log.Fatal("SSH Server failure", "error", err)
	}
}
