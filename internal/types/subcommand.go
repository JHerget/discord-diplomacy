package types

import (
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Subcommand struct {
	Name        string
	Description string
	Handler     func(cctx *utils.CommandContext) error
}

func (s Subcommand) Create() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        s.Name,
		Description: s.Description,
	}
}

type SubcommandMap map[string]Subcommand

func NewSubcommandMap(subcommands []Subcommand) SubcommandMap {
	sm := SubcommandMap{}

	for _, subcommand := range subcommands {
		sm[subcommand.Name] = subcommand
	}

	return sm
}

func (sm SubcommandMap) Options() []*discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{}

	for _, subcommand := range sm {
		options = append(options, subcommand.Create())
	}

	return options
}
