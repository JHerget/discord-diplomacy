package models

type Order struct {
	ID          string `json:"id"`
	PhaseID     string `json:"phaseID"`
	PlayerName  string `json:"playerName"`
	CreatedDate int    `json:"createdDate"`
	Value       string `json:"value"`
}

type CreateOrderRequest struct {
	PlayerName string `json:"playerName"`
	PhaseID    string `json:"phaseID"`
	Value      string `json:"value"`
}
