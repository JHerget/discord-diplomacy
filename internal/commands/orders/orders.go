package orders

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
			SubmitSubcommand,
		}),
	}
}

func (c Command) Register(registry *interactions.Registry) error {
	if err := registry.AddCommand(&discordgo.ApplicationCommand{
		Name:        "orders",
		Description: "Manage your orders",
		Options:     c.subcommands.Options(),
	}, c); err != nil {
		return err
	}

	return registry.AddModal(submitModalCustomID, c)
}

func (c Command) Handle(cctx *types.CommandContext) error {
	options := cctx.Interaction.ApplicationCommandData().Options
	if len(options) != 1 {
		return errors.New("orders command requires exactly one subcommand")
	}

	option := options[0]
	if option.Type != discordgo.ApplicationCommandOptionSubCommand {
		return fmt.Errorf("orders option %q is not a subcommand", option.Name)
	}

	subcommand, ok := c.subcommands[option.Name]
	if !ok {
		return fmt.Errorf("unknown orders subcommand %q", option.Name)
	}

	return subcommand.Handler(cctx)
}
