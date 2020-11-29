package repository

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/internal/entities"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockFinanceAPI struct {
	mock.Mock
}

func (mockFinanceAPI *MockFinanceAPI) GetCurrency(key string) (entities.RegularMarketPrice, error) {
	args := mockFinanceAPI.Called(key)
	return args.Get(0).(entities.RegularMarketPrice), args.Error(1).(error)
}

func TestStockRepository_GetCurrency(t *testing.T) {
	key := "tokenKey"
	getCurrency := "GetCurrency"

	t.Run("GetCurrencyWithoutHost", func(t *testing.T) {

		outputDTOIN := entities.RegularMarketPrice{}

		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(entities.RegularMarketPrice{}, errorNew)
		output, err := ProvideStockRepository().GetCurrency(key)
		expected := "stock-api returned Status : 500 and text : internal server error"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())
	})

	t.Run("GetCurrencyWithoutBase_URL", func(t *testing.T) {

		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		outputDTOIN := entities.RegularMarketPrice{}

		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(entities.RegularMarketPrice{}, errorNew)
		output, err := ProvideStockRepository().GetCurrency(key)
		expected := "stock-api returned Status : 500 and text : internal server error"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())

	})

	//TODO: HIER BENÖTIGE ICH HILFE!!!!!!!!!!!!!!!!!!

	t.Run("GetCurrencyWithFullResponse", func(t *testing.T) {

		//
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "{1.3}")
		}))
		defer server.Close()

		req := httptest.NewRequest("GET", server.URL, nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		fmt.Println(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()

		/*	os.Setenv("BASE_URL", server.URL+"/")
			os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")*/

		/*outputEntity := entities.RegularMarketPrice{Raw: 1.3}

		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(outputEntity, nil)

		output, _ := ProvideStockRepository().GetCurrency(key)
		assert.Equal(t, string(bodyBytes), output.Raw)*/

	})
}
