package api

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (mockService *MockService) CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error) {
	args := mockService.Called(key, responseEntity)
	return args.Get(0).(dto.OutputDTO), args.Error(1)
}

func TestStockController_CalculateEarning(t *testing.T) {

	CalculateEarning := "CalculateEarning"

	t.Run("CalculateEarningWithStatusBadRequest", func(t *testing.T) {

		errorNew := &customError.ErrorStock{Code: 400, Text: "Bad Request"}
		serviceMock := MockService{}
		serviceMock.On(CalculateEarning, "", dto.OutputDTO{}).Return(dto.OutputDTO{}, errorNew)

		req := httptest.NewRequest("GET", "/HVB_Stock_API/calcGain", bytes.NewBuffer(nil))

		recorder := httptest.NewRecorder()
		router := gin.Default()
		CreateRouter(router, &serviceMock)

		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		expected := `"bad request"`
		assert.Equal(t, expected, recorder.Body.String())

	})

	t.Run("CalculateEarningWithStatusOK", func(t *testing.T) {

		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b"

		outputDTOIN := dto.OutputDTO{
			Share:   "AAPL",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}
		outputDTOOUT := dto.OutputDTO{
			Share:   "AAPL",
			Numbers: 1,
			Price:   100,
			Xrate:   1.3,
			Value:   130,
		}

		errorNew := &customError.ErrorStock{Code: 200, Text: "Status OK"}
		serviceMock := MockService{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(outputDTOOUT, errorNew)
		r := gin.Default()
		CreateRouter(r, &serviceMock)

		req := httptest.NewRequest("GET", "/HVB_Stock_API/calcGain?share=AAPL&number=1&singlePrice=100", bytes.NewBuffer(nil))
		req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b")
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)

		expected := `{"share":"AAPL","numbers":1,"price":100,"xrate":1.3,"value":130}`
		assert.Equal(t, expected, recorder.Body.String())
	})

	t.Run("CalculateEarningWithStatusUnauthorized", func(t *testing.T) {

		key := ""

		outputDTOIN := dto.OutputDTO{
			Share:   "AAPL",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}

		errorNew := &customError.ErrorStock{Code: 401, Text: "401 Unauthorized"}
		serviceMock := MockService{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(dto.OutputDTO{}, errorNew)

		req := httptest.NewRequest("GET", "/HVB_Stock_API/calcGain?share=AAPL&number=1&singlePrice=100", bytes.NewBuffer(nil))

		recorder := httptest.NewRecorder()
		r := gin.Default()
		CreateRouter(r, &serviceMock)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		expected := `"401 Unauthorized"`
		assert.Equal(t, expected, recorder.Body.String())
	})
	t.Run("CalculateEarningWithStatusNoContent", func(t *testing.T) {

		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b"

		outputDTOIN := dto.OutputDTO{
			Share:   "APPL",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}

		errorNew := &customError.ErrorStock{Code: 204, Text: "no content"}
		serviceMock := MockService{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(dto.OutputDTO{}, errorNew)

		req := httptest.NewRequest("GET", "/HVB_Stock_API/calcGain?share=APPL&number=1&singlePrice=100", bytes.NewBuffer(nil))
		req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b")
		recorder := httptest.NewRecorder()
		r := gin.Default()
		CreateRouter(r, &serviceMock)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusNoContent, recorder.Code)

		expected := ``
		assert.Equal(t, expected, recorder.Body.String())
	})
	t.Run("CalculateEarningWithStatusForbidden", func(t *testing.T) {

		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5c"

		outputDTOIN := dto.OutputDTO{
			Share:   "APPL",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}

		errorNew := &customError.ErrorStock{Code: 403, Text: "403 Forbidden"}
		serviceMock := MockService{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(dto.OutputDTO{}, errorNew)

		req := httptest.NewRequest("GET", "/HVB_Stock_API/calcGain?share=APPL&number=1&singlePrice=100", bytes.NewBuffer(nil))
		req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5c")
		recorder := httptest.NewRecorder()
		r := gin.Default()
		CreateRouter(r, &serviceMock)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusForbidden, recorder.Code)

		expected := `"403 Forbidden"`
		assert.Equal(t, expected, recorder.Body.String())
	})
	t.Run("CalculateEarningWithStatusForbidden", func(t *testing.T) {

		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5c"

		outputDTOIN := dto.OutputDTO{
			Share:   "APPL",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}

		errorNew := &customError.ErrorStock{Code: 500, Text: "internal server error"}
		serviceMock := MockService{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(dto.OutputDTO{}, errorNew)

		req := httptest.NewRequest("GET", "/HVB_Stock_API/calcGain?share=APPL&number=1&singlePrice=100", bytes.NewBuffer(nil))
		req.Header.Add("x-rapidapi-key", "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5c")
		recorder := httptest.NewRecorder()
		r := gin.Default()
		CreateRouter(r, &serviceMock)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		expected := `"internal server error"`
		assert.Equal(t, expected, recorder.Body.String())
	})
}
