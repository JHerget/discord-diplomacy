package game

import (
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
)

var JoinSubcommand = types.Subcommand{
	Name:        "join",
	Description: "Join the current game",
	Handler: func(cctx *utils.CommandContext) error {
		return cctx.BasicEphemeralResponse("Join the current game")
	},
}
