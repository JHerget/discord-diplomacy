package game

import "github.com/bwmarrin/discordgo"

type CreateSubcommand struct {
	Name        string
	Description string
}

func NewCreateSubcommand() CreateSubcommand {
	return CreateSubcommand{
		Name: "create",
		Description: "Create a new game",
	}
}

func (c CreateSubcommand) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c.Description,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

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
