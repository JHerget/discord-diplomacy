package models

type Phase struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PhaseOrder  int8   `json:"phaseOrder"`
}
