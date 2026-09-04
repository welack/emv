package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Account struct {
	IMAPHost      string `json:"imap_host"`
	IMAPPort      int    `json:"imap_port"`
	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	User          string `json:"user"`
	TLSSkipVerify bool   `json:"tls_skip_verify"`
	IMAPTLSMode   string `json:"imap_tls_mode"`
}

func load_config(app_name string) (*Account, error) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config dir: %w", err)
	}

	config_path := filepath.Join(config_dir, app_name, "config.json")

	file_bytes, err := os.ReadFile(config_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var acc Account
	if err := json.Unmarshal(file_bytes, &acc); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &acc, nil
}

func save_config(app_name string, acc *Account) error {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user config dir: %w", err)
	}

	path := filepath.Join(config_dir, app_name, "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("could not create config dir: %w", err)
	}
	b, err := json.MarshalIndent(acc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return os.WriteFile(path, b, 0o600)
}
