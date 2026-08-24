package feedback

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

const ()

type Command struct {
	modalCustomID string
	inputCustomID string
}

func New() Command {
	return Command{
		modalCustomID: "feedback:submit",
		inputCustomID: "feedback:message",
	}
}

func (c Command) Register(registry *interactions.Registry) error {
	if err := registry.AddCommand(&discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Open the feedback form",
	}, c); err != nil {
		return err
	}

	if err := registry.AddModal(c.modalCustomID, c); err != nil {
		return err
	}

	return nil
}

func (c Command) Handle(cctx *utils.CommandContext) error {
	return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: c.modalCustomID,
			Title:    "Feedback",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    c.inputCustomID,
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

func (c Command) Submit(cctx *utils.CommandContext) error {
	inputs := interactions.TextInputValues(cctx.Interaction.ModalSubmitData())
	value := inputs[c.inputCustomID]

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("feedback message is empty")
	}

	return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for your feedback!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
