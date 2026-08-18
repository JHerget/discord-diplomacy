package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config contains the values needed to connect to Discord and register guild commands.
type Config struct {
	BotToken string
	GuildID  string
}

// Load reads and validates configuration from the process environment.
func Load() (Config, error) {
	cfg := Config{
		BotToken: strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN")),
		GuildID:  strings.TrimSpace(os.Getenv("DISCORD_GUILD_ID")),
	}

	var validationErrors []error
	if cfg.BotToken == "" {
		validationErrors = append(validationErrors, fmt.Errorf("DISCORD_BOT_TOKEN is required"))
	}
	if cfg.GuildID == "" {
		validationErrors = append(validationErrors, fmt.Errorf("DISCORD_GUILD_ID is required"))
	}

	if err := errors.Join(validationErrors...); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
