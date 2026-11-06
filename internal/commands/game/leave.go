package game

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
)

var LeaveSubcommand = types.Subcommand{
	Name:        "leave",
	Description: "Leave the current game.",
	Handler: func(cctx *types.CommandContext) error {
		if cctx.ActiveGame == nil {
			return cctx.BasicEphemeralResponse("There is no active game.")
		}

		userID, ok := cctx.InteractionUserID()
		if !ok {
			return cctx.BasicEphemeralResponse("Could not identify the invoking user.")
		}

		API := diplomacy.NewAPI()
		game, err := API.GetGame(*cctx.ActiveGame)
		if err != nil {
			return cctx.BasicEphemeralResponse(err.Error())
		}

		names := make([]string, 0)
		for _, player := range game.Players {
			if player.UserID == nil || *player.UserID != userID {
				continue
			}

			if err := API.DeletePlayer(game.ID, player.ID); err != nil {
				return cctx.BasicEphemeralResponse(err.Error())
			}

			names = append(names, player.Name)
		}

		if len(names) > 0 {
			if game.ExternalID != nil {
				playerNames := strings.Join(names, "/")
				cctx.BasicChannelResponse(*game.ExternalID, fmt.Sprintf("%s (%s) left the game.", cctx.InteractionUserName(), playerNames))
			}
			return cctx.BasicEphemeralResponse("You left the game.")
		}

		return cctx.BasicEphemeralResponse("You are not in the active game.")
	},
}
