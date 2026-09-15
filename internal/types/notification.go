package types

import (
	"fmt"
	"strings"
)

const maxDiscordContentLength = 2000

type NotificationMessage struct {
	ChannelID string `json:"channelId"`
	GameID    string `json:"gameId"`
	TurnID    string `json:"turnID"`
}

func (m NotificationMessage) Validate() error {
	if strings.TrimSpace(m.ChannelID) == "" {
		return fmt.Errorf("channelId is required")
	}
	if strings.TrimSpace(m.GameID) == "" {
		return fmt.Errorf("gameId is required")
	}
	if strings.TrimSpace(m.TurnID) == "" {
		return fmt.Errorf("turnId is required")
	}

	return nil
}
