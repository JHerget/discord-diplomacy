package models

import (
	"fmt"
)

type Turn struct {
	ID         string  `json:"id"`
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

func (t *Turn) PrettyPrint() string {
	if len(t.Orders) == 0 {
		return ""
	}

	var value string
	for _, order := range t.Orders {
		value = fmt.Sprintf(`%s\n\n### %s\n%s`, value, order.PlayerName, order.PrettyPrint())
	}

	return value
}
