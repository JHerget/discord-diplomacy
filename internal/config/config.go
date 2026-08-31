package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BotToken    string
	GuildID     string
	AWSRegion   string
	SQSQueueURL string
	ActiveGame  *string
}

func Load() (Config, error) {
	cfg := Config{
		BotToken:    strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN")),
		GuildID:     strings.TrimSpace(os.Getenv("DISCORD_GUILD_ID")),
		AWSRegion:   strings.TrimSpace(os.Getenv("AWS_REGION")),
		SQSQueueURL: strings.TrimSpace(os.Getenv("SQS_QUEUE_URL")),
	}

	var validationErrors []error
	if cfg.BotToken == "" {
		validationErrors = append(validationErrors, fmt.Errorf("DISCORD_BOT_TOKEN is required"))
	}
	if cfg.GuildID == "" {
		validationErrors = append(validationErrors, fmt.Errorf("DISCORD_GUILD_ID is required"))
	}
	if cfg.AWSRegion == "" {
		validationErrors = append(validationErrors, fmt.Errorf("AWS_REGION is required"))
	}
	if cfg.SQSQueueURL == "" {
		validationErrors = append(validationErrors, fmt.Errorf("SQS_QUEUE_URL is required"))
	}
	if err := errors.Join(validationErrors...); err != nil {
		return Config{}, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/.diplomacy/discord-diplomacy.json", home))
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

func (c Config) Save() error {
	gameID := ""
	if c.ActiveGame != nil {
		gameID = *c.ActiveGame
	}

	content, err := json.MarshalIndent(map[string]string{
		c.GuildID: gameID,
	}, "", "    ")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := fmt.Sprintf("%s/.diplomacy/discord-diplomacy.json", home)
	if err := os.WriteFile(dir, content, 0644); err != nil {
		return err
	}

	return nil
}
