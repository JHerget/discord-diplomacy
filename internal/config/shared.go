package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Shared struct {
	ActiveGame *string
}

func LoadShared(guildID string) (*Shared, error) {
	cfg := &Shared{}

	data, err := os.ReadFile("/var/lib/discord-diplomacy.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, err
	}

	if len(data) > 0 {
		var activeGames map[string]string
		if err := json.Unmarshal(data, &activeGames); err != nil {
			return cfg, err
		}

		gameId, ok := activeGames[guildID]
		if ok {
			cfg.ActiveGame = &gameId
		}
	}

	return cfg, nil
}
