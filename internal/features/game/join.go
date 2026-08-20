package game

import "github.com/bwmarrin/discordgo"

const joinSubcommandName = "join"

func joinSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        joinSubcommandName,
		Description: "Join the active game",
	}
}

func (Module) handleJoin(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) error {
	return respondEphemeral(session, interaction, "Joining a game is not available yet.")
}
