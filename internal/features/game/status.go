package game

import (
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
)

var StatusSubcommand = types.Subcommand{
	Name:        "status",
	Description: "Get the status of the current game",
	Handler: func(cctx *utils.CommandContext) error {
		return cctx.BasicEphemeralResponse("Get the status of the current game")
	},
}
