package models

import (
	"fmt"
)

type Turn struct {
	ID         string  `json:"id"`
	PhaseID    string  `json:"phaseID"`
	Orders     []Order `json:"orders"`
	TurnNumber int     `json:"turnNumber"`
	StartDate  int     `json:"startDate"`
	EndDate    int     `json:"endDate"`
}

func (t *Turn) Name() string {
	turn := (t.TurnNumber + 1) / 2
	season := "Spring"

	if t.TurnNumber%2 == 0 {
		season = "Fall"
	}

	return fmt.Sprintf("%s %d", season, 1900+turn)
}
