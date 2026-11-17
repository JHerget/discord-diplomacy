package orders

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
)

var ViewSubcommand = types.Subcommand{
	Name:        "view",
	Description: "View your orders for the current turn.",
	Handler: func(cctx *types.CommandContext) error {
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
		turn := g.CurrentTurn()
		if turn == nil {
			return cctx.BasicEphemeralResponse("The current game has no turns in progress.")
		}

		var sections []string
		for _, player := range g.Players {
			if player.UserID == nil || *player.UserID != userID {
				continue
			}

			var values []string
			for _, order := range turn.Orders {
				if order.PlayerName == player.Name {
					values = append(values, strings.ReplaceAll(order.Value, ",", "\n"))
				}
			}

			if len(values) > 0 {
				sections = append(sections, fmt.Sprintf("**%s**\n%s", player.Name, strings.Join(values, "\n")))
			}
		}
		if len(sections) == 0 {
			return cctx.BasicEphemeralResponse("No orders submitted yet.")
		}

		response := strings.Join(sections, "\n\n")
		runes := []rune(response)
		if len(runes) > maxResponseLength {
			response = string(runes[:maxResponseLength-len([]rune("\n… additional orders omitted."))]) + "\n… additional orders omitted."
		}

		return cctx.BasicEphemeralResponse(response)
	},
}
