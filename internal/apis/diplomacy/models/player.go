package models

type Player struct {
	ID        string  `json:"id"`
	UserID    *string `json:"userID"`
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	IsPlaying bool    `json:"isPlaying"`
}

type CreatePlayerRequest struct {
	UserID *string `json:"userID"`
}
