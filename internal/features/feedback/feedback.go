package feedback

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/types"

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

func (c Command) Register(registry *interactions.Registry) error {
	if err := registry.AddCommand(&discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Open the feedback form",
	}, c); err != nil {
		return err
	}

	if err := registry.AddModal(modalCustomID, c); err != nil {
		return err
	}

	return nil
}

func (c Command) Handle(cctx *types.CommandContext) error {
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
							Required:    boolPtr(true),
							MaxLength:   1000,
						},
					},
				},
			},
		},
	})
}

func boolPtr(value bool) *bool {
	return &value
}

func (c Command) Submit(cctx *types.CommandContext) error {
	inputs := interactions.TextInputValues(cctx.Interaction.ModalSubmitData())
	value := inputs[inputCustomID]

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
