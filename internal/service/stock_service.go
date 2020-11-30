package service

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"errors"
	"strings"
)

type stockRepository interface {
	CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error)
}

type StockService struct {
	stockRepository stockRepository
}

func ProvideStockService(repository stockRepository) *StockService {
	return &StockService{stockRepository: repository}
}

func (stockService *StockService) CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error) {

	//encapsulate service-layer from repository-layer
	response, err := stockService.stockRepository.CalculateEarning(key, responseEntity)

	//declare entity with pointer
	var stockError *customError.ErrorStock

	//compare err.code. if not 200, return empty entity and filled error
	if errors.As(err, &stockError) {
		if stockError.Code != 200 {
			return dto.OutputDTO{}, err
		}
	}

	//example "btc-eur" is already the correct value
	//and dont need to be recalculated with the exchange rate EUR/USD
	//filter between stock with and without EURo in string
	if strings.Contains(strings.ToLower(response.Share), "eur") {
		response.Value = (response.Value - response.Price) * response.Numbers
		return response, &customError.ErrorStock{Code: 200, Text: "Status OK"}
	}

	//calculate with exchange rate EUR/USD
	response.Value = ((response.Value / response.Xrate) - response.Price) * response.Numbers

	//defer fully reponse with errEntity status 200 back to controller - layer
	return response, &customError.ErrorStock{Code: 200, Text: "Status OK"}
}
