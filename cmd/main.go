package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://rapidapi.p.rapidapi.com/stock/v2/get-chart?interval=5m&symbol=AMRN&range=1d&region=US"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b")
	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)

	fmt.Println(string(body))

}
