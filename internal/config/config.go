package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BotToken   string
	GuildID    string
	ActiveGame *string
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

	data, err := os.ReadFile("/var/lib/discord-diplomacy.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, err
	}

	if len(data) > 0 {
		var activeGames map[string]string
		if err := json.Unmarshal(data, &activeGames); err != nil {
			return cfg, err
		}

		gameId, ok := activeGames[cfg.GuildID]
		if ok {
			cfg.ActiveGame = &gameId
		}
	}

	return cfg, nil
}
