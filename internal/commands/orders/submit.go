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
	value := strings.TrimSpace(values[submitInputCustomID])
	if value == "" {
		return cctx.BasicEphemeralResponse("Enter orders before submitting.")
	}

	if _, err := API.CreateOrder(g.ID, turn.ID, models.CreateOrderRequest{
		PlayerName: player.Name,
		PhaseID:    turn.PhaseID,
		Value:      value,
	}); err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	return cctx.BasicEphemeralResponse(fmt.Sprintf("Orders submitted for %s.", turn.Name()))
}
