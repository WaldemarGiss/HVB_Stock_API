package repository

import (
	"HVB_Stock_API/internal/entities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StockRepository struct {
}

func ProvideStockRepository() *StockRepository {
	return &StockRepository{}
}

func (stockRepository StockRepository) Calculate(tenant string) (entities.RegularMarketPrice, float64) {
	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-statistics?symbol=" + tenant + "&region=DE"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b")
	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response entities.Result

	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Print(err)
	}
	fak := stockRepository.GetCurrency()

	//fmt.Printf("%+v\n", response.Price.RegularMarketPrice.Raw)
	return response.Price.RegularMarketPrice, fak.Raw

}

func (stockRepository StockRepository) GetCurrency() entities.RegularMarketPrice {
	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-statistics?symbol=EURUSD=X&region=DE"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b")
	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response entities.Result

	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Print(err)
	}
	return response.Price.RegularMarketPrice
}

//TODO: Env auslesen?
