package types

type OBUData struct {
	ObuID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"long"`
}

type Distance struct {
	Value float64 `json:"value"`
	ObuID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type Invoice struct {
	ObuID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}
