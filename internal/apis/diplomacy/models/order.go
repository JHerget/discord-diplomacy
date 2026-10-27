package models

type Order struct {
	ID          string `json:"id"`
	PlayerName  string `json:"playerName"`
	CreatedDate int    `json:"createdDate"`
	Value       string `json:"value"`
}

type OrderCommand struct {
	PlayerName string `json:"playerName"`
	Value      string `json:"value"`
	Valid      bool   `json:"valid"`
}

type CreateOrderRequest struct {
	PlayerName string `json:"playerName"`
	Value      string `json:"value"`
}
