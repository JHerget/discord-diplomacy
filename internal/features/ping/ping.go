package ping

import (
	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Command struct{}

func New() Command {
	return Command{}
}

func (c Command) Register(registry *interactions.Registry) error {
	return registry.AddCommand(&discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check whether the bot is responding",
	}, c)
}

func (Command) Handle(cctx *utils.CommandContext) error {
	return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
