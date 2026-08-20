package feedback

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/interactions"

	"github.com/bwmarrin/discordgo"
)

const (
	modalCustomID = "feedback:submit"
	inputCustomID = "feedback:message"
)

type Module struct {
	cfg *config.Shared
}

func New(cfg *config.Shared) Module {
	return Module{
		cfg: cfg,
	}
}

func (m Module) Register(registry *interactions.Registry) error {
	if err := registry.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Open the feedback form",
	}, openModal); err != nil {
		return fmt.Errorf("register feedback command: %w", err)
	}

	if err := registry.RegisterModal(modalCustomID, submitModal); err != nil {
		return fmt.Errorf("register feedback modal: %w", err)
	}
	return nil
}

func openModal(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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
