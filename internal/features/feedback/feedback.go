package feedback

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	modalCustomID = "feedback:submit"
	inputCustomID = "feedback:message"
)

type Command struct{}

func New() Command {
	return Command{}
}

func (Command) Registration() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Open the feedback form",
	}
}

func (Command) Handle(cctx *utils.CommandContext) error {
	return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: modalCustomID,
			Title:    "Feedback",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    inputCustomID,
							Label:       "What would you like to share?",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Enter your feedback",
							Required:    true,
							MaxLength:   1000,
						},
					},
				},
			},
		},
	})
}

func submitModal(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	value, err := interactions.TextInputValue(interaction.ModalSubmitData(), inputCustomID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("feedback message is empty")
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for your feedback!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
