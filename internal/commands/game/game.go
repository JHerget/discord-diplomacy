package game

import (
	"errors"
	"fmt"

	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/types"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	subcommands types.SubcommandMap
}

func New() Command {
	return Command{
		subcommands: types.NewSubcommandMap([]types.Subcommand{
			CreateSubcommand,
			JoinSubcommand,
			LeaveSubcommand,
			StatusSubcommand,
		}),
	}
}

func (c Command) Register(registry *interactions.Registry) error {
	if err := registry.AddCommand(&discordgo.ApplicationCommand{
		Name:        "game",
		Description: "Manage a Diplomacy game",
		Options:     c.subcommands.Options(),
	}, c); err != nil {
		return err
	}

	return registry.AddModal(createModalCustomID, c)
}

func (c Command) Handle(cctx *types.CommandContext) error {
	options := cctx.Interaction.ApplicationCommandData().Options
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

	return subcommand.Handler(cctx)
}
