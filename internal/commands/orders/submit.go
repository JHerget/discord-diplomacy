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
	submitModalCustomID = "orders:submit"
	submitInputCustomID = "orders:value"
	maxResponseLength   = 2000
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
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  submitInputCustomID,
								Label:     "What are your orders?",
								Style:     discordgo.TextInputParagraph,
								Required:  utils.BoolPtr(true),
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

	player, err := g.FindPlayerByUserID(userID)
	if err != nil {
		return cctx.BasicEphemeralResponse("You are not in the active game.")
	}

	turn := g.CurrentTurn()
	if turn == nil {
		return cctx.BasicEphemeralResponse("The current game has no turns in progress.")
	}

	values := interactions.TextInputValues(cctx.Interaction.ModalSubmitData())

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

	order, err := API.CreateOrder(g.ID, turn.ID, models.CreateOrderRequest{
		PlayerName: player.Name,
		Value:      value,
	})
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	commands, err := API.GetOrderCommands(g.ID, turn.ID, order.ID)
	return cctx.BasicEphemeralResponse(submissionResponse(turn.Name(), commands, err))
}

func submissionResponse(turnName string, commands []models.OrderCommand, commandErr error) string {
	response := fmt.Sprintf("Orders submitted for %s.", turnName)
	if commandErr != nil {
		return response + "\n\nYour order formatting could not be checked."
	}

	invalid := make([]models.OrderCommand, 0)
	for _, command := range commands {
		if !command.Valid {
			invalid = append(invalid, command)
		}
	}
	if len(invalid) == 0 {
		return response
	}

	response += "\n\nThe formatting of these commands was not recognized:"
	included := 0
	for i, command := range invalid {
		entry := fmt.Sprintf("\n- %s", escapeInlineCode(command.Value))
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
	return strings.ReplaceAll(value, "`", "\\`")
}
