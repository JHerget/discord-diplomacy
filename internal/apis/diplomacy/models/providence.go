package models

type Providence struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	SupplyCenter  *SupplyCenter      `json:"supplyCenter"`
	Unit          *Unit              `json:"unit"`
	Coordinates   Coordinates        `json:"coordinates"`
	Type          ProvidenceType     `json:"type"`
	Routes        []string           `json:"routes"`
	CoastalRoutes map[Coast][]string `json:"coastalRoutes"`
}
