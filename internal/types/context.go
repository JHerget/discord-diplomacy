package types

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

func (cc *CommandContext) InteractionUserID() (string, bool) {
	if cc.Interaction.Member != nil && cc.Interaction.Member.User != nil {
		return cc.Interaction.Member.User.ID, true
	}
	if cc.Interaction.User != nil {
		return cc.Interaction.User.ID, true
	}

	return "", false
}

func (cc *CommandContext) InteractionUserName() string {
	if cc.Interaction.Member != nil {
		if cc.Interaction.Member.Nick != "" {
			return cc.Interaction.Member.Nick
		}
		if cc.Interaction.Member.User != nil {
			if cc.Interaction.Member.User.Username != "" {
				return cc.Interaction.Member.User.Username
			}
			if cc.Interaction.Member.User.GlobalName != "" {
				return cc.Interaction.Member.User.GlobalName
			}
		}
	}
	if cc.Interaction.User != nil {
		if cc.Interaction.User.Username != "" {
			return cc.Interaction.User.Username
		}
		if cc.Interaction.User.GlobalName != "" {
			return cc.Interaction.User.GlobalName
		}
	}

	return "Unknown user"
}

func (cc *CommandContext) basicResponse(content string, flag discordgo.MessageFlags) error {
	return cc.Session.InteractionRespond(cc.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flag,
		},
	})
}

func (cc *CommandContext) imageResponse(content string, image ImageResponse, flag discordgo.MessageFlags) error {
	return cc.Session.InteractionRespond(cc.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flag,
			Files: []*discordgo.File{
				{
					Name:        image.Filename,
					ContentType: image.ContentType,
					Reader:      bytes.NewReader(image.Data),
				},
			},
		},
	})
}

func (cc *CommandContext) BasicResponse(content string) error {
	return cc.basicResponse(content, discordgo.MessageFlags(0))
}

func (cc *CommandContext) BasicEphemeralResponse(content string) error {
	return cc.basicResponse(content, discordgo.MessageFlagsEphemeral)
}

func (cc *CommandContext) ImageResponse(content string, image ImageResponse) error {
	return cc.imageResponse(content, image, discordgo.MessageFlags(0))
}

func (cc *CommandContext) EphemeralImageResponse(content string, image ImageResponse) error {
	return cc.imageResponse(content, image, discordgo.MessageFlagsEphemeral)
}
