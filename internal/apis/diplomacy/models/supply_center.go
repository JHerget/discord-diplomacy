package models

type SupplyCenter struct {
	ControlledBy *string     `json:"controlledBy"`
	Coordinates  Coordinates `json:"coordinates"`
}
