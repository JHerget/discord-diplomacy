package models

import "strings"

type Order struct {
	ID          string `json:"id"`
	PlayerName  string `json:"playerName"`
	CreatedDate int    `json:"createdDate"`
	Value       string `json:"value"`
}

type CreateOrderRequest struct {
	PlayerName string `json:"playerName"`
	Value      string `json:"value"`
}

func (o *Order) PrettyPrint() string {
	return strings.ReplaceAll(o.Value, ",", "\n")
}
