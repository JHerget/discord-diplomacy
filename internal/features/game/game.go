package game

import (
	"errors"
	"fmt"

	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	cctx        *utils.CommandContext
	subcommands types.SubcommandMap
}

func New(cctx *utils.CommandContext) Command {
	return Command{
		cctx: cctx,
		subcommands: types.NewSubcommandMap([]types.Subcommand{
			CreateSubcommand,
			JoinSubcommand,
			LeaveSubcommand,
			StatusSubcommand,
		}),
	}
}

func (c Command) Register(registry *interactions.Registry) error {
	return registry.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "game",
		Description: "Manage a Diplomacy game",
		Options:     c.subcommands.Options(),
	}, c.handle)
}

func (c Command) handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	options := interaction.ApplicationCommandData().Options
	if len(options) != 1 {
		return errors.New("game command requires exactly one subcommand")
	}

	option := options[0]
	if option.Type != discordgo.ApplicationCommandOptionSubCommand {
		return fmt.Errorf("game option %q is not a subcommand", option.Name)
	}

	subcommand, ok := c.subcommands[option.Name]
	if !ok {
		return fmt.Errorf("unknown game subcommand %q", option.Name)
	}

	return subcommand.Handler(c.cctx)
}
