package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Result struct {
	Price Price `json:"price"`
}

type Price struct {
	RegularMarketPrice RegularMarketPrice `json:"regularMarketPrice"`
}

type RegularMarketPrice struct {
	Raw float64 `json:"raw"`
}

func main() {

	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-statistics?symbol=TSLA&region=DE"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b")
	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var ausgabe Result

	err := json.Unmarshal(body, &ausgabe)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%+v\n", ausgabe.Price.RegularMarketPrice.Raw)

	//fmt.Println(res)

	//fmt.Println(string(body))

}
