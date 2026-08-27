package game

import (
	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"
	"fmt"
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
					turn := g.CurrentTurn()
					turnName := "No turns yet"
					if turn != nil {
						turnName = turn.Name()
					}

					message := fmt.Sprintf(`# %s\n## %s`, g.Map.Name, turnName)

					// TODO: return the board image too.
					return cctx.BasicEphemeralResponse(message)
				} else {
					return cctx.BasicEphemeralResponse("The current game is not in progress.")
				}
			}

		}

		return cctx.BasicEphemeralResponse("No game is being tracked.")
	},
}
