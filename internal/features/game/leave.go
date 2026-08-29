package game

import (
	"fmt"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
)

var LeaveSubcommand = types.Subcommand{
	Name:        "leave",
	Description: "Leave the current game",
	Handler: func(cctx *types.CommandContext) error {
		if cctx.ActiveGame == nil {
			return cctx.BasicEphemeralResponse(ErrNoActiveGame)
		}

		userID, ok := cctx.InteractionUserID()
		if !ok {
			return cctx.BasicEphemeralResponse(ErrUnknownUser)
		}

		API := diplomacy.NewAPI()
		g, err := API.GetGame(*cctx.ActiveGame)
		if err != nil {
			return cctx.BasicEphemeralResponse(err.Error())
		}

		for _, player := range g.Players {
			if player.UserID == nil || *player.UserID != userID {
				continue
			}

			if err := API.DeletePlayer(g.ID, player.ID); err != nil {
				return cctx.BasicEphemeralResponse(err.Error())
			}

			return cctx.BasicResponse(fmt.Sprintf("%s (%s) left the game.", cctx.InteractionUserName(), player.Name))
		}

		return cctx.BasicEphemeralResponse("You are not in the active game.")
	},
}
