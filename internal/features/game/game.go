package game

import (
	"errors"
	"fmt"

	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/interactions"

	"github.com/bwmarrin/discordgo"
)

type subcommandHandler func(
	*discordgo.Session,
	*discordgo.InteractionCreate,
	*discordgo.ApplicationCommandInteractionDataOption,
) error

type Module struct {
	cfg      *config.Shared
	handlers map[string]subcommandHandler
}

func New(cfg *config.Shared) Module {
	m := Module{
		cfg:      cfg,
		handlers: make(map[string]subcommandHandler),
	}
	m.handlers[createSubcommandName] = m.handleCreate
	m.handlers[joinSubcommandName] = m.handleJoin
	m.handlers[statusSubcommandName] = m.handleStatus
	return m
}

func (m Module) Register(registry *interactions.Registry) error {
	return registry.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "game",
		Description: "Manage a Diplomacy game",
		Options: []*discordgo.ApplicationCommandOption{
			createSubcommand(),
			joinSubcommand(),
			statusSubcommand(),
		},
	}, m.handle)
}

func (m Module) handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	options := interaction.ApplicationCommandData().Options
	if len(options) != 1 {
		return errors.New("game command requires exactly one subcommand")
	}

	subcommand := options[0]
	if subcommand.Type != discordgo.ApplicationCommandOptionSubCommand {
		return fmt.Errorf("game option %q is not a subcommand", subcommand.Name)
	}

	handler, ok := m.handlers[subcommand.Name]
	if !ok {
		return fmt.Errorf("unknown game subcommand %q", subcommand.Name)
	}

	return handler(session, interaction, subcommand)
}

func respondEphemeral(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	content string,
) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
