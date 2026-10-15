package models

import "errors"

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

func (g *Game) FindTurnByNumber(turnNumber int) (*Turn, bool) {
	for _, turn := range g.Turns {
		if turn.TurnNumber == turnNumber {
			return &turn, true
		}
	}

	return nil, false
}

type CreateGameRequest struct {
	ExternalID    *string `json:"externalId"`
	MapID         string  `json:"mapId"`
	DaysPerTurn   int     `json:"daysPerTurn"`
	TurnStartHour int     `json:"turnStartHour"`
	Timezone      int     `json:"timezone"`
	StartDate     int     `json:"startDate"`
}
