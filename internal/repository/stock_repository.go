package repository

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"HVB_Stock_API/internal/entities"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type StockRepository struct {
}

func ProvideStockRepository() *StockRepository {
	return &StockRepository{}
}

func (stockRepository StockRepository) CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error) {

	//fetch password from .env.file.
	//return empty entity with filled errEntity if not OK
	host, ok := os.LookupEnv("HOST")
	if ok != true {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//fetch baseURL from .env.file.
	//return empty entity with filled errEntity if not OK
	baseUrl, ok := os.LookupEnv("BASE_URL")
	if ok != true {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//assemble URL for request
	url := baseUrl + "get-statistics?symbol=" + responseEntity.Share + "&region=DE"
	//url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-holders"

	//initialize the request.
	//ignoring error. it will catched at http.DefaultClient.Do()
	req, _ := http.NewRequest("GET", url, nil)

	//fetch necessary information from header
	req.Header.Add("x-rapidapi-key", key)
	req.Header.Add("x-rapidapi-host", host)

	//client do http.NewRequest.
	//return empty entity and filled errEntity if Client Do fail
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}
	defer res.Body.Close()

	//return reponse directly from extern finance api
	//empty entity and filled errEntity
	if res.StatusCode == 403 {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: res.StatusCode, Text: res.Status}
	} else if res.StatusCode == 401 {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: res.StatusCode, Text: res.Status}
	}

	//read response from Client.
	//ignoring error. will be catched in json.unmarshal
	body, _ := ioutil.ReadAll(res.Body)

	//declare responseEntity
	var response entities.Result

	//unmarshal responseClient into responseEntity
	//throw empty entity and filled errEntity
	err = json.Unmarshal(body, &response)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//fill responseEntity with value
	responseEntity.Value = response.Price.RegularMarketPrice.Raw

	//return empty entity and filled errEntity
	if responseEntity.Value == 0 {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 204, Text: "no content"}
	}

	//get exchangerate EUR/USD
	//return empty entity and filled errEntity if err != nil
	fak, err := stockRepository.GetCurrency(key)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//fill responseEntity with exchangeRate
	responseEntity.Xrate = fak.Raw

	//return filled entity and filled errEntity
	return responseEntity, &customError.ErrorStock{Code: 200, Text: "Status OK"}
}

//
//Method to get the current exchange-rate EUR/USD
//
func (stockRepository StockRepository) GetCurrency(key string) (entities.RegularMarketPrice, error) {

	//fetch password from .env.file.
	//return empty entity with filled errEntity if not OK
	host, ok := os.LookupEnv("HOST")
	if ok != true {
		return entities.RegularMarketPrice{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//fetch baseURL from .env.file.
	//return empty entity with filled errEntity if not OK
	baseUrl, ok := os.LookupEnv("BASE_URL")
	if ok != true {
		return entities.RegularMarketPrice{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//assemble URL for request
	url := baseUrl + "get-statistics?symbol=EURUSD=X&region=DE"

	//initialize the request.
	//ignoring error. it will catched at http.DefaultClient.Do()
	request, _ := http.NewRequest("GET", url, nil)

	//fetch necessary information from header
	request.Header.Add("x-rapidapi-key", key)
	request.Header.Add("x-rapidapi-host", host)

	//client do http.NewRequest.
	//return empty entity and filled errEntity if Client Do fail
	responseClient, err := http.DefaultClient.Do(request)
	if err != nil {
		return entities.RegularMarketPrice{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}
	defer responseClient.Body.Close()

	//read response from Client.
	//ignoring error. will be catched in json.unmarshal
	body, _ := ioutil.ReadAll(responseClient.Body)

	//declare responseEntity
	var responseEntity entities.Result

	//unmarshal responseClient into responseEntity
	//throw empty entity and filled errEntity
	err = json.Unmarshal(body, &responseEntity)
	if err != nil {
		return entities.RegularMarketPrice{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	//after success return responseEntity.Price.RegularMarketPrice
	//and error as nil
	return responseEntity.Price.RegularMarketPrice, nil
}
