package service

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"strings"
)

type stockRepository interface {
	Calculate(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock)
}

type StockService struct {
	stockRepository stockRepository
}

func ProvideStockService(repository stockRepository) *StockService {
	return &StockService{stockRepository: repository}
}

func (stockService *StockService) Calculate(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock) {

	response, err := stockService.stockRepository.Calculate(key, responseEntity)

	if err.Code != 0 {
		return dto.OutputDTO{}, err
	}

	if strings.Contains(strings.ToLower(response.Share), "eur") {
		response.Value = (response.Value - response.Price) * response.Numbers
		return response, customError.ErrorStock{}
	}

	response.Value = ((response.Value / response.Xrate) - response.Price) * response.Numbers
	return response, customError.ErrorStock{}
}
