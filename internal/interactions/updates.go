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
	if !g.InProgress {
		return nil, fmt.Errorf("The game with id '%s' is not in progress.", gameID)
	}

	board, err := API.GetBoard(g.ID)
	if err != nil {
		return nil, fmt.Errorf("Error getting board: %s", err)
	}

	turn, ok := g.FindTurn(turnID)
	if !ok {
		return nil, fmt.Errorf("The turn id '%s' is invalid.", turnID)
	}

	message := fmt.Sprintf(`
	# %s
	## %s
	`, g.Map.Name, turn.Name())
	filename := fmt.Sprintf("diplomacy-%s.png", strings.ReplaceAll(strings.ToLower(turn.Name()), " ", "-"))

	return &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{
			{
				Name:        filename,
				ContentType: "image/png",
				Reader:      bytes.NewReader(board),
			},
		},
	}, nil
}
