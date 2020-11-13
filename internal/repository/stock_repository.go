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

func (stockRepository StockRepository) Calculate(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock) {

	host, baseUrl := getEnv()

	url := baseUrl + "get-statistics?symbol=" + responseEntity.Share + "&region=DE"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.OutputDTO{}, customError.ErrorStock{Code: 500, Text: "internal server error"}
	}

	req.Header.Add("x-rapidapi-key", key)
	req.Header.Add("x-rapidapi-host", host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.OutputDTO{}, customError.ErrorStock{Code: 401, Text: "unauthorized"}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("cant read body")
	}

	var response entities.Result

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("cant unmarshal into entity")
	}

	responseEntity.Value = response.Price.RegularMarketPrice.Raw

	fak, err := stockRepository.GetCurrency(key)

	if responseEntity.Value == 0 && fak.Raw == 0 {
		fmt.Println("401")
		return dto.OutputDTO{}, customError.ErrorStock{Code: 401, Text: "unauthorized"}
	} else if responseEntity.Value == 0 {
		fmt.Println("204")
		return dto.OutputDTO{}, customError.ErrorStock{Code: 204, Text: "input wrong"}
	}

	if err != nil {
		return dto.OutputDTO{}, customError.ErrorStock{Code: 500, Text: "could not get exchangerate"}
	}
	responseEntity.Xrate = fak.Raw
	return responseEntity, customError.ErrorStock{}

}

func (stockRepository StockRepository) GetCurrency(key string) (entities.RegularMarketPrice, error) {

	host, baseUrl := getEnv()

	url := baseUrl + "get-statistics?symbol=EURUSD=X&region=DE"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return entities.RegularMarketPrice{}, err
	}

	req.Header.Add("x-rapidapi-key", key)
	req.Header.Add("x-rapidapi-host", host)

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
