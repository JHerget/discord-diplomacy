package ping

import (
	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/interactions"

	"github.com/bwmarrin/discordgo"
)

type Module struct {}

func New(cfg *config.Shared) Module {
	return Module{}
}

func (Module) Register(registry *interactions.Registry) error {
	return registry.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check whether the bot is responding",
	}, handle)
}

func handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
