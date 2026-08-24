package ping

import (
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Command struct{}

func New() Command {
	return Command{}
}

func (Command) Registration() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check whether the bot is responding",
	}
}

func (Command) Handle(cctx *utils.CommandContext) error {
	return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
