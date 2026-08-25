package game

import (
	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
)

var StatusSubcommand = types.Subcommand{
	Name:        "status",
	Description: "Get the status of the current game",
	Handler: func(cctx *utils.CommandContext) error {
		API := diplomacy.NewAPI()

		if cctx.ActiveGame != nil {
			g, err := API.GetGame(*cctx.ActiveGame)
			if err == nil {
				if g.InProgress {
					return cctx.BasicEphemeralResponse("The current game is in progress.")
				} else {
					return cctx.BasicEphemeralResponse("The current game is not in progress.")
				}
			}

		}

		return cctx.BasicEphemeralResponse("There is no game in progress.")
	},
}
