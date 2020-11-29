package repository

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"HVB_Stock_API/internal/entities"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockFinanceAPI struct {
	mock.Mock
}

func (mockFinanceAPI *MockFinanceAPI) CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error) {
	args := mockFinanceAPI.Called(key, responseEntity)
	return args.Get(0).(dto.OutputDTO), args.Error(1).(error)
}

func TestStockRepository_CalculateEarning(t *testing.T) {
	key := "tokenKey"
	calculateEarning := "CalculateEarning"
	getCurrency := "GetCurrency"

	responseEntity := dto.OutputDTO{
		Share:   "AAPL",
		Numbers: 1,
		Price:   100,
		Xrate:   0,
		Value:   0,
	}

	t.Run("GetCalculateEarningWithoutHost", func(t *testing.T) {
		outputDTOIN := dto.OutputDTO{}

		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(calculateEarning, key, responseEntity).Return(dto.OutputDTO{}, errorNew)
		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)
		expected := "stock-api returned Status : 500 and text : internal server error"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())
	})
	t.Run("GetCalculateEarningWithoutBaseUrl", func(t *testing.T) {
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")
		outputDTOIN := dto.OutputDTO{}

		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(calculateEarning, key, responseEntity).Return(dto.OutputDTO{}, errorNew)
		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)
		expected := "stock-api returned Status : 500 and text : internal server error"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())
	})
	t.Run("CalculateEarningWithFullResponse", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"price": {"regularMarketPrice":{"raw":1000}}}`)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		outputEntity := entities.RegularMarketPrice{Raw: 1.3}

		//outputEntity := entities.RegularMarketPrice{}
		errorResponse := &customError.ErrorStock{Code: 200, Text: "Status OK"}

		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(outputEntity, errorResponse)

		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)
		assert.Equal(t, errorResponse, err)

		exception := 1000.00
		assert.Equal(t, exception, output.Value)
	})
	t.Run("CalculateEarningCantDoClient", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"price": {"regularMarketPrice":{"raw":1000}}}`)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", "||||||")

		errorResponse := &customError.ErrorStock{Code: 500, Text: "internal server error"}

		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)
		assert.Equal(t, errorResponse, err)
		assert.Equal(t, dto.OutputDTO{}, output)
	})
	t.Run("CalculateEarningStatusUnauthorized ", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		errorResponse := &customError.ErrorStock{Code: 401, Text: "401 Unauthorized"}

		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)

		assert.Equal(t, errorResponse, err)
		assert.Equal(t, dto.OutputDTO{}, output)
	})
	t.Run("CalculateEarningStatusForbidden", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(403)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		errorResponse := &customError.ErrorStock{Code: 403, Text: "403 Forbidden"}

		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)

		assert.Equal(t, errorResponse, err)
		assert.Equal(t, dto.OutputDTO{}, output)
	})
	t.Run("CalculateEarningCantUnmarshal", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, "text")
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		errorResponse := &customError.ErrorStock{Code: 500, Text: "internal server error"}

		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)

		assert.Equal(t, errorResponse, err)
		assert.Equal(t, dto.OutputDTO{}, output)
	})
	t.Run("CalculateEarningStatusNoContent", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"price": {"regularMarketPrice":{"raw":0}}}`)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		errorResponse := &customError.ErrorStock{Code: 204, Text: "no content"}

		output, err := ProvideStockRepository().CalculateEarning(key, responseEntity)

		assert.Equal(t, errorResponse, err)
		assert.Equal(t, dto.OutputDTO{}, output)
	})
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

		output, err := ProvideStockRepository().GetCurrency(key)
		expected := "stock-api returned Status : 500 and text : internal server error"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())
	})

	t.Run("GetCurrencyWithoutBase_URL", func(t *testing.T) {

		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		outputDTOIN := entities.RegularMarketPrice{}

		output, err := ProvideStockRepository().GetCurrency(key)
		expected := "stock-api returned Status : 500 and text : internal server error"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())

	})

	t.Run("GetCurrencyWithFullResponse", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"price": {"regularMarketPrice":{"raw":1.3}}}`)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		outputEntity := entities.RegularMarketPrice{}

		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(outputEntity, nil)

		output, _ := ProvideStockRepository().GetCurrency(key)
		exception := 1.3
		assert.Equal(t, exception, output.Raw)
	})
	t.Run("GetCurrencyCantDoClient", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `some text`)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", "||||||||||||||||||||||")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		outputEntity := entities.RegularMarketPrice{}
		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(outputEntity, errorNew)

		_, err := ProvideStockRepository().GetCurrency(key)
		exception := errorNew
		assert.Equal(t, exception, err)
	})
	t.Run("GetCurrencyCantUnmarshalBody", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `some text`)
		}))
		defer server.Close()

		os.Setenv("BASE_URL", server.URL+"/")
		os.Setenv("host", "apidojo-yahoo-finance-v1.p.rapidapi.com")

		outputEntity := entities.RegularMarketPrice{}
		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		mockFinanceApi := MockFinanceAPI{}
		mockFinanceApi.On(getCurrency, key).Return(outputEntity, errorNew)

		_, err := ProvideStockRepository().GetCurrency(key)
		exception := errorNew
		assert.Equal(t, exception, err)
	})
}
