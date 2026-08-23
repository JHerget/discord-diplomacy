package utils

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/bwmarrin/discordgo"
)

type CommandContext struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	ActiveGame  *string
}

func NewCommandContext(guildID string) (*CommandContext, error) {
	cctx := &CommandContext{}

	data, err := os.ReadFile("/var/lib/discord-diplomacy.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return cctx, err
	}

	if len(data) > 0 {
		var activeGames map[string]string
		if err := json.Unmarshal(data, &activeGames); err != nil {
			return cctx, err
		}

		gameId, ok := activeGames[guildID]
		if ok {
			cctx.ActiveGame = &gameId
		}
	}

	return cctx, nil
}

func (cc *CommandContext) BasicResponse(content string) error {
	return cc.Session.InteractionRespond(cc.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func (cc *CommandContext) BasicEphemeralResponse(content string) error {
	return cc.Session.InteractionRespond(cc.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
