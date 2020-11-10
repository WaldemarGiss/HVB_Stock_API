package entities

type Result struct {
	Price Price `json:"price"`
}

type Price struct {
	RegularMarketPrice RegularMarketPrice `json:"regularMarketPrice"`
}

type RegularMarketPrice struct {
	Raw float64 `json:"raw"`
}
