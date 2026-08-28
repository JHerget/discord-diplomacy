package utils

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
)

type CommandContext struct {
	Session           *discordgo.Session
	Interaction       *discordgo.InteractionCreate
	ActiveGame        *string
	GuildID           string
	SetActiveGameFunc func(string)
}

func (cc *CommandContext) SetActiveGame(gameID string) {
	cc.ActiveGame = &gameID

	if cc.SetActiveGameFunc != nil {
		cc.SetActiveGameFunc(gameID)
	}
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

func (cc *CommandContext) EphemeralImageResponse(content string, filename string, contentType string, data []byte) error {
	return cc.Session.InteractionRespond(cc.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
			Files: []*discordgo.File{
				{
					Name:        filename,
					ContentType: contentType,
					Reader:      bytes.NewReader(data),
				},
			},
		},
	})
}
