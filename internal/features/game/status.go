package game

import (
	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
	"fmt"
	"strings"
)

var StatusSubcommand = types.Subcommand{
	Name:        "status",
	Description: "Get the status of the current game",
	Handler: func(cctx *types.CommandContext) error {
		API := diplomacy.NewAPI()

		if cctx.ActiveGame == nil {
			return cctx.BasicEphemeralResponse("No game is being tracked.")
		}

		g, err := API.GetGame(*cctx.ActiveGame)
		if err != nil {
			return cctx.BasicEphemeralResponse("No game is being tracked.")
		}
		if !g.InProgress {
			return cctx.BasicEphemeralResponse("The current game is not in progress.")
		}

		board, err := API.GetBoard(g.ID)
		if err != nil {
			return cctx.BasicEphemeralResponse(fmt.Sprintf("Error getting board: %v", err))
		}

		turn := g.CurrentTurn()
		turnName := "No turns yet"
		if turn != nil {
			turnName = turn.Name()
		}

		message := fmt.Sprintf(`
		# %s
		## %s
		`, g.Map.Name, turnName)
		filename := fmt.Sprintf("diplomacy-%s.png", strings.ReplaceAll(strings.ToLower(turnName), " ", "-"))

		return cctx.EphemeralImageResponse(message, types.ImageResponse{
			Filename:    filename,
			ContentType: "image/png",
			Data:        board,
		})
	},
}
