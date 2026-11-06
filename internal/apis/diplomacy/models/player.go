package models

type Player struct {
	ID        string  `json:"id"`
	UserID    *string `json:"userId"`
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	IsPlaying bool    `json:"isPlaying"`
}

type CreatePlayersRequest struct {
	UserID      string   `json:"userId"`
	PlayerNames []string `json:"playerNames"`
}
