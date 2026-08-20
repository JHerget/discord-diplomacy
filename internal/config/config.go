package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BotToken string
	GuildID  string
}

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
