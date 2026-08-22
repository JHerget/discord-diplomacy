package models

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type LocationReference struct {
	ID    string `json:"id"`
	Coast *Coast `json:"coast"`
}
