package game

import (
	"fmt"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"discord-diplomacy/internal/types"
)

var JoinSubcommand = types.Subcommand{
	Name:        "join",
	Description: "Join the current game.",
	Handler: func(cctx *types.CommandContext) error {
		if cctx.ActiveGame == nil {
			return cctx.BasicEphemeralResponse("There is no active game.")
		}

		userID, ok := cctx.InteractionUserID()
		if !ok {
			return cctx.BasicEphemeralResponse("Could not identify the invoking user.")
		}

		API := diplomacy.NewAPI()
		player, err := API.CreatePlayer(*cctx.ActiveGame, models.CreatePlayerRequest{
			UserID: &userID,
		})
		if err != nil {
			return cctx.BasicEphemeralResponse(err.Error())
		}

		return cctx.BasicResponse(fmt.Sprintf("%s joined the game as %s", cctx.InteractionUserName(), player.Name))
	},
}
