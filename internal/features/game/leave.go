package game

import (
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
)

var LeaveSubcommand = types.Subcommand{
	Name:        "leave",
	Description: "Leave the current game",
	Handler: func(cctx *utils.CommandContext) error {
		return cctx.BasicEphemeralResponse("Leave the current game")
	},
}
