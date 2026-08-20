package game

import "github.com/bwmarrin/discordgo"

const statusSubcommandName = "status"

func statusSubcommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        statusSubcommandName,
		Description: "Show the active game",
	}
}

func (m Module) handleStatus(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) error {
	if m.cfg == nil || m.cfg.ActiveGame == nil {
		return respondEphemeral(session, interaction, "There is no active game.")
	}

	return respondEphemeral(session, interaction, "Active game: "+*m.cfg.ActiveGame)
}
