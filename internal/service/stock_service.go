package service

import (
	"HVB_Stock_API/internal/entities"
	"strings"
)

type stockRepository interface {
	Calculate(tenant string) (entities.RegularMarketPrice, float64)
}

type StockService struct {
	stockRepository stockRepository
}

func ProvideStockService(repository stockRepository) *StockService {
	return &StockService{stockRepository: repository}
}

func (stockService *StockService) Calculate(tenant string) float64 {

	stock, xRate := stockService.stockRepository.Calculate(tenant)
	if strings.Contains(strings.ToLower(tenant), "eur") {
		return stock.Raw
	}
	return stock.Raw / xRate
}
