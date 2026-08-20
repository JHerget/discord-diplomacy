package game

import "github.com/bwmarrin/discordgo"

const createSubcommandName = "create"

func createSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        createSubcommandName,
		Description: "Create a new game",
	}
}

func (Module) handleCreate(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) error {
	return respondEphemeral(session, interaction, "Game creation is not available yet.")
}
