package types

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	messageTypeDiscordMessage = "discord_message"
	maxDiscordContentLength   = 2000
)

type NotificationMessage struct {
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

func (m NotificationMessage) Validate() error {
	if m.Type != messageTypeDiscordMessage {
		return fmt.Errorf("unsupported message type %q", m.Type)
	}
	if strings.TrimSpace(m.ChannelID) == "" {
		return fmt.Errorf("channel_id is required")
	}
	if strings.TrimSpace(m.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if utf8.RuneCountInString(m.Content) > maxDiscordContentLength {
		return fmt.Errorf("content exceeds %d characters", maxDiscordContentLength)
	}

	return nil
}
