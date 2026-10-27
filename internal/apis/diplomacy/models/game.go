package models

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Game struct {
	ID            string       `json:"id"`
	ExternalID    *string      `json:"externalId"`
	OwnerID       string       `json:"ownerID"`
	Map           MapSummary   `json:"map"`
	Board         []Providence `json:"board"`
	Players       []Player     `json:"players"`
	Turns         []Turn       `json:"turns"`
	DaysPerTurn   int          `json:"daysPerTurn"`
	TurnStartHour int          `json:"turnStartHour"`
	Timezone      int          `json:"timezone"`
	StartDate     int          `json:"startDate"`
	EndDate       int          `json:"endDate"`
	InProgress    bool         `json:"inProgress"`
	IsDeleted     bool         `json:"isDeleted"`
}

type CreateGameRequest struct {
	ExternalID    *string `json:"externalId"`
	MapID         string  `json:"mapId"`
	DaysPerTurn   int     `json:"daysPerTurn"`
	TurnStartHour int     `json:"turnStartHour"`
	Timezone      int     `json:"timezone"`
	StartDate     int     `json:"startDate"`
}

func (g *Game) CurrentTurn() *Turn {
	target := len(g.Turns)
	if target <= 0 {
		return nil
	}

	for _, turn := range g.Turns {
		if turn.TurnNumber == target {
			return &turn
		}
	}

	return nil
}

func (g *Game) FindPlayerByUserID(userID string) (*Player, error) {
	for i := range g.Players {
		if g.Players[i].UserID != nil && *g.Players[i].UserID == userID {
			return &g.Players[i], nil
		}
	}

	return nil, errors.New("player not found")
}

func (g *Game) FindTurn(turnID string) (*Turn, bool) {
	for _, turn := range g.Turns {
		if turn.ID == turnID {
			return &turn, true
		}
	}

	return nil, false
}

func (g *Game) GameEndSummary() string {
	playerScores := make(map[string]int)

	for _, p := range g.Board {
		if p.SupplyCenter == nil {
			continue
		}
		if p.SupplyCenter.ControlledBy == nil {
			continue
		}

		playerScores[*p.SupplyCenter.ControlledBy]++
	}

	highestScore := 0
	for _, score := range playerScores {
		if score > highestScore {
			highestScore = score
		}
	}

	var leaders []string
	for player, score := range playerScores {
		if score == highestScore {
			leaders = append(leaders, player)
		}
	}
	sort.Strings(leaders)

	switch len(leaders) {
	case 0:
		return "The game has ended!"
	case 1:
		return fmt.Sprintf("The game has ended! %s wins!", leaders[0])
	default:
		return fmt.Sprintf("The game has ended in a draw between %s!", strings.Join(leaders, " and "))
	}
}
