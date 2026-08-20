package ping

import (
	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/interactions"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Module struct {
	cfg *config.Shared
}

func New(cfg *config.Shared) Module {
	return Module{
		cfg: cfg,
	}
}

func (m Module) Register(registry *interactions.Registry) error {
	return registry.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check whether the bot is responding",
	}, m.handle)
}

func (m Module) handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	value := ""
	if m.cfg.ActiveGame != nil {
		value = *m.cfg.ActiveGame
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Pong! (Active game: %s)", value),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
