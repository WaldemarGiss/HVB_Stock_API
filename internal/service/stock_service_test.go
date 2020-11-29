package service

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (mockRepository *MockRepository) CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error) {
	args := mockRepository.Called(key, responseEntity)
	return args.Get(0).(dto.OutputDTO), args.Error(1)
}

func TestStockService_CalculateEarning(t *testing.T) {
	CalculateEarning := "CalculateEarning"

	t.Run("CalculateEarningWithStatusNot200", func(t *testing.T) {
		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b"

		outputDTOIN := dto.OutputDTO{
			Share:   "",
			Numbers: 0,
			Price:   0,
			Xrate:   0,
			Value:   0,
		}

		errorNew := &customError.ErrorStock{Code: 400, Text: "Bad Request"}
		serviceMock := MockRepository{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(dto.OutputDTO{}, errorNew)
		output, err := ProvideStockService(&serviceMock).CalculateEarning(key, outputDTOIN)
		expected := "stock-api returned Status : 400 and text : Bad Request"
		assert.Equal(t, outputDTOIN, output)
		assert.Equal(t, expected, err.Error())

	})

	t.Run("CalculateEarningWithStatus200AndWithEUR", func(t *testing.T) {
		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b"

		outputDTOIN := dto.OutputDTO{
			Share:   "BTC-EUR",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}
		outputDTOOUT := dto.OutputDTO{
			Share:   "BTC-EUR",
			Numbers: 1,
			Price:   100,
			Xrate:   1.3,
			Value:   130,
		}

		errorNew := &customError.ErrorStock{Code: 200, Text: "Status OK"}
		serviceMock := MockRepository{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(outputDTOOUT, errorNew)
		output, err := ProvideStockService(&serviceMock).CalculateEarning(key, outputDTOIN)

		expectedEntity := dto.OutputDTO{
			Share:   "BTC-EUR",
			Numbers: 1,
			Price:   100,
			Xrate:   1.3,
			Value:   30,
		}
		assert.Equal(t, expectedEntity, output)

		expectedErr := &customError.ErrorStock{Code: 200, Text: "Status OK"}
		assert.Equal(t, expectedErr, err)

	})
	t.Run("CalculateEarningWithStatus200AndWithoutEUR", func(t *testing.T) {
		key := "4dfc8434famsh4013439d5a8c726p1a0e1djsn1e63e0ec5b5b"

		outputDTOIN := dto.OutputDTO{
			Share:   "BTC",
			Numbers: 1,
			Price:   100,
			Xrate:   0,
			Value:   0,
		}
		outputDTOOUT := dto.OutputDTO{
			Share:   "BTC",
			Numbers: 1,
			Price:   100,
			Xrate:   1.3,
			Value:   130,
		}

		errorNew := &customError.ErrorStock{Code: 200, Text: "Status OK"}
		serviceMock := MockRepository{}
		serviceMock.On(CalculateEarning, key, outputDTOIN).Return(outputDTOOUT, errorNew)
		output, err := ProvideStockService(&serviceMock).CalculateEarning(key, outputDTOIN)

		expectedEntity := dto.OutputDTO{
			Share:   "BTC",
			Numbers: 1,
			Price:   100,
			Xrate:   1.3,
			Value:   0,
		}
		assert.Equal(t, expectedEntity, output)

		expectedErr := &customError.ErrorStock{Code: 200, Text: "Status OK"}
		assert.Equal(t, expectedErr, err)

	})
}
