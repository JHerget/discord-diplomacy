package models

type Turn struct {
	ID         string  `json:"id"`
	PhaseID    string  `json:"phaseID"`
	Orders     []Order `json:"orders"`
	TurnNumber int     `json:"turnNumber"`
	StartDate  int     `json:"startDate"`
	EndDate    int     `json:"endDate"`
}
