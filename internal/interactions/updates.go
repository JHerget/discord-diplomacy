package interactions

import (
	"bytes"
	"discord-diplomacy/internal/apis/diplomacy"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func TurnSummary(gameID string, turnID string) (*discordgo.MessageSend, error) {
	API := diplomacy.NewAPI()

	g, err := API.GetGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("The game id '%s' is invalid.", gameID)
	}

	board, err := API.GetBoard(g.ID)
	if err != nil {
		return nil, fmt.Errorf("Error getting board: %s", err)
	}

	turn, ok := g.FindTurn(turnID)
	if !ok {
		return nil, fmt.Errorf("The turn id '%s' is invalid.", turnID)
	}

	var message string
	previousTurn, previousTurnFound := g.FindTurnByNumber(turn.TurnNumber - 1)

	if g.InProgress {
		switch previousTurnFound {
		case true:
			message = fmt.Sprintf("# %s\n## Previous Turn: %s\n%s\n## New Turn: %s",
				g.Map.Name,
				previousTurn.Name(),
				previousTurn.PrettyPrint(),
				turn.Name(),
			)
		case false:
			message = fmt.Sprintf("# %s\n## New Turn: %s", g.Map.Name, turn.Name())
		}
	} else {
		message = fmt.Sprintf("# %s\n## %s\n%s", g.Map.Name, turn.Name(), g.GameEndSummary())
	}

	return &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{
			{
				Name:        fmt.Sprintf("diplomacy-%s.png", strings.ReplaceAll(strings.ToLower(turn.Name()), " ", "-")),
				ContentType: "image/png",
				Reader:      bytes.NewReader(board),
			},
		},
	}, nil
}
