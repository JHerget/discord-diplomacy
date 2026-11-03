package orders

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	submitModalCustomID        = "orders:submit"
	greatPowerInputCustomID    = "orders:great-power"
	submitInputCustomID        = "orders:value"
	unavailableGreatPowerValue = "unavailable"
	maxResponseLength          = 2000
)

var SubmitSubcommand = types.Subcommand{
	Name:        "submit",
	Description: "Submit orders for the current turn.",
	Handler: func(cctx *types.CommandContext) error {
		return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: submitModalCustomID,
				Title:    "Submit Orders",
				Components: []discordgo.MessageComponent{
					discordgo.Label{
						Label: "Great Power",
						Component: discordgo.SelectMenu{
							CustomID:    greatPowerInputCustomID,
							Placeholder: "Select a Great Power",
							MinValues:   new(1),
							MaxValues:   1,
							Required:    new(true),
							Options:     greatPowerSelectOptions(cctx),
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  submitInputCustomID,
								Label:     "What are your orders?",
								Style:     discordgo.TextInputParagraph,
								Required:  new(true),
								MaxLength: 4000,
							},
						},
					},
				},
			},
		})
	},
}

func (c Command) Submit(cctx *types.CommandContext) error {
	if cctx.ActiveGame == nil {
		return cctx.BasicEphemeralResponse("There is no active game.")
	}

	userID, ok := cctx.InteractionUserID()
	if !ok {
		return cctx.BasicEphemeralResponse("Could not identify the invoking user.")
	}

	API := diplomacy.NewAPI()
	g, err := API.GetGame(*cctx.ActiveGame)
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}
	if !g.InProgress {
		return cctx.BasicEphemeralResponse("The current game is not in progress.")
	}

	if _, err := g.FindPlayerByUserID(userID); err != nil {
		return cctx.BasicEphemeralResponse("You are not in the active game.")
	}

	turn := g.CurrentTurn()
	if turn == nil {
		return cctx.BasicEphemeralResponse("The current game has no turns in progress.")
	}

	values := interactions.ModalValues(cctx.Interaction.ModalSubmitData())

	orders := strings.FieldsFunc(values[submitInputCustomID], func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	for i := range orders {
		orders[i] = strings.TrimSpace(orders[i])
	}

	value := strings.Join(orders, ",")
	if value == "" {
		return cctx.BasicEphemeralResponse("Enter orders before submitting.")
	}

	commands, err := API.CreateOrder(g.ID, turn.ID, models.CreateOrderRequest{
		PlayerName: values[greatPowerInputCustomID],
		Value:      value,
	})
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	return cctx.BasicEphemeralResponse(submissionResponse(turn.Name(), commands))
}

func greatPowerSelectOptions(cctx *types.CommandContext) []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0)

	if cctx.ActiveGame != nil {
		API := diplomacy.NewAPI()

		allPlayers, err := API.GetPlayers(*cctx.ActiveGame)
		userID, ok := cctx.InteractionUserID()
		if err == nil && ok {
			for _, player := range allPlayers {
				if player.UserID == nil || *player.UserID != userID {
					continue
				}
				if player.IsPlaying == false {
					continue
				}

				options = append(options, discordgo.SelectMenuOption{
					Label:   player.Name,
					Value:   player.Name,
					Default: false,
				})
			}
		}
	}

	switch len(options) {
	case 0:
		return []discordgo.SelectMenuOption{
			{
				Label:   "No Great Powers available",
				Value:   unavailableGreatPowerValue,
				Default: true,
			},
		}
	case 1:
		options[0].Default = true
		return options
	default:
		return options
	}
}

func submissionResponse(turnName string, commands models.CommandValidationResults) string {
	response := fmt.Sprintf("Orders submitted for %s.", turnName)

	invalid := make(models.CommandValidationResults, 0)
	for _, command := range commands {
		if !command.Valid {
			invalid = append(invalid, command)
		}
	}
	if len(invalid) == 0 {
		return response
	}

	response += "\n\nThere were validation errors with the following commands:"
	included := 0
	for i, command := range invalid {
		entry := fmt.Sprintf("\n- %s", escapeInlineCode(command.Value))
		if command.Error != "" {
			entry += fmt.Sprintf(" (%s)", command.Error)
		}
		remaining := len(invalid) - i - 1
		candidate := response + entry
		if remaining > 0 {
			candidate += omittedCommandsMessage(remaining)
		}
		if len([]rune(candidate)) > maxResponseLength {
			break
		}

		response += entry
		included++
	}

	if omitted := len(invalid) - included; omitted > 0 {
		response += omittedCommandsMessage(omitted)
	}

	return response
}

func omittedCommandsMessage(count int) string {
	return fmt.Sprintf("\n- ... and %d additional invalid command(s) were omitted.", count)
}

func escapeInlineCode(value string) string {
	return strings.ReplaceAll(utils.ToTitle(value), "`", "\\`")
}
