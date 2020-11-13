package dto

type OutputDTO struct {
	Share   string  `json:"share"`
	Numbers float64 `json:"numbers"`
	Price   float64 `json:"price"`
	Xrate   float64 `json:"xrate"`
	Value   float64 `json:"value"`
}
