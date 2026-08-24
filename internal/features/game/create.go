package game

import (
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
)

var CreateSubcommand = types.Subcommand{
	Name:        "create",
	Description: "Create a new game",
	Handler: func(cctx *utils.CommandContext) error {
		// API := diplomacy.NewAPI()

		// g, err := API.GetGame(*cctx.ActiveGame)
		// if err != nil {
		// 	return cctx.BasicEphemeralResponse(err.Error())
		// }

		return cctx.BasicEphemeralResponse(*cctx.ActiveGame)
	},
}
