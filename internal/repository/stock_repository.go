package repository

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"HVB_Stock_API/internal/entities"
	"encoding/json"
	"fmt"
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

	host, baseUrl := getEnv()

	url := baseUrl + "get-statistics?symbol=" + responseEntity.Share + "&region=DE"
	//url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-holders"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}
	if key != "" {
		req.Header.Add("x-rapidapi-key", key)
		req.Header.Add("x-rapidapi-host", host)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	defer res.Body.Close()

	if res.StatusCode == 403 {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: res.StatusCode, Text: res.Status}
	} else if res.StatusCode == 401 {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: res.StatusCode, Text: res.Status}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	var response entities.Result

	err = json.Unmarshal(body, &response)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	responseEntity.Value = response.Price.RegularMarketPrice.Raw

	if responseEntity.Value == 0 {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 204, Text: "no content"}
	}

	fak, err := stockRepository.GetCurrency(key)
	if err != nil {
		return dto.OutputDTO{}, &customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	responseEntity.Xrate = fak.Raw

	return responseEntity, &customError.ErrorStock{Code: 200, Text: "Status OK"}
}

func (stockRepository StockRepository) GetCurrency(key string) (entities.RegularMarketPrice, error) {

	host, baseUrl := getEnv()

	url := baseUrl + "get-statistics?symbol=EURUSD=X&region=DE"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return entities.RegularMarketPrice{}, err
	}

	if key != "" {
		req.Header.Add("x-rapidapi-key", key)
		req.Header.Add("x-rapidapi-host", host)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return entities.RegularMarketPrice{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return entities.RegularMarketPrice{}, err
	}
	var response entities.Result

	err = json.Unmarshal(body, &response)
	if err != nil {
		return entities.RegularMarketPrice{}, err
	}

	return response.Price.RegularMarketPrice, nil
}

func getEnv() (string, string) {

	//fetch password
	host, ok := os.LookupEnv("HOST")
	if ok != true {
		fmt.Println("can't find a HOST")
	}

	//fetch baseURL
	baseURL, ok := os.LookupEnv("BASE_URL")
	if ok != true {
		fmt.Println("can't find a base_url")
	}

	return host, baseURL
}
