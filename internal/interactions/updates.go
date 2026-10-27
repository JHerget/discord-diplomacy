package interactions

import (
	"bytes"
	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/types"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func TurnSummary(notification types.NotificationMessage) (*discordgo.MessageSend, error) {
	API := diplomacy.NewAPI()

	g, err := API.GetGame(notification.GameID)
	if err != nil {
		return nil, fmt.Errorf("The game id '%s' is invalid.", notification.GameID)
	}

	board, err := API.GetBoard(g.ID)
	if err != nil {
		return nil, fmt.Errorf("Error getting board: %s", err)
	}

	turn, ok := g.FindTurn(notification.TurnID)
	if !ok {
		return nil, fmt.Errorf("The turn id '%s' is invalid.", notification.TurnID)
	}

	var message string
	if g.InProgress {
		switch len(notification.PreviousTurnOrders) > 0 {
		case true:
			message = fmt.Sprintf("# %s\n## Previous Turn\n%s\n## New Turn: %s",
				g.Map.Name,
				notification.PreviousTurnOrders.PrettyPrint(),
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
