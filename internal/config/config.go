package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds the full application configuration.
type Config struct {
	Telegram TelegramConfig `toml:"telegram"`
	Sync     SyncConfig     `toml:"sync"`
}

// TelegramConfig holds Telegram bot credentials and chat filters.
type TelegramConfig struct {
	Token   string  `toml:"token"`
	ChatIDs []int64 `toml:"chat_ids"`
}

// SyncConfig holds output settings.
type SyncConfig struct {
	InboxDir string `toml:"inbox_dir"`
}

// Load reads the config from
// ~/.config/telegram-obsidian-sync/config.toml.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(
		home,
		".config",
		"telegram-obsidian-sync",
		"config.toml",
	)
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Sync.InboxDir) > 1 && cfg.Sync.InboxDir[0] == '~' {
		cfg.Sync.InboxDir = filepath.Join(
			home, cfg.Sync.InboxDir[1:],
		)
	}
	return &cfg, nil
}
