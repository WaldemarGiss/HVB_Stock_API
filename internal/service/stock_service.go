package service

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"strings"
)

type stockRepository interface {
	CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock)
}

type StockService struct {
	stockRepository stockRepository
}

func ProvideStockService(repository stockRepository) *StockService {
	return &StockService{stockRepository: repository}
}

func (stockService *StockService) CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock) {

	response, err := stockService.stockRepository.CalculateEarning(key, responseEntity)

	if err.Code != 200 {
		return dto.OutputDTO{}, err
	}

	if strings.Contains(strings.ToLower(response.Share), "eur") {
		response.Value = (response.Value - response.Price) * response.Numbers
		return response, customError.ErrorStock{Code: 200, Text: "status OK"}
	}

	response.Value = ((response.Value / response.Xrate) - response.Price) * response.Numbers
	return response, customError.ErrorStock{Code: 200, Text: "status OK"}
}
