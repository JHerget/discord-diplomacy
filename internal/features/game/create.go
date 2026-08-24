package game

import (
	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
)

var CreateSubcommand = types.Subcommand{
	Name:        "create",
	Description: "Create a new game",
	Handler: func(cctx *utils.CommandContext) error {
		API := diplomacy.NewAPI()

		if cctx.ActiveGame != nil {
			g, err := API.GetGame(*cctx.ActiveGame)
			if err == nil && g.InProgress {
				return cctx.BasicEphemeralResponse("There is already a game in progress for this server.")
			}

		}

		g, err := API.CreateGame(models.CreateGameRequest{
			ExternalID:    &cctx.GuildID,
			MapID:         "",
			DaysPerTurn:   14,
			TurnStartHour: 12,
			Timezone:      -7,
			StartDate:     0,
		})
		if err != nil {
			return cctx.BasicEphemeralResponse(err.Error())
		}

		return cctx.BasicEphemeralResponse(g.ID)
	},
}
