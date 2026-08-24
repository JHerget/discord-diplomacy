package models

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
