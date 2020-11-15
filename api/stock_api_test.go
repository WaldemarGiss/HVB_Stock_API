package api

import (
	"HVB_Stock_API/api/dto"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (mockService *MockService) Calculate(baseURL string, tenant string, username string, password string) (dto.OutputDTO, error) {
	args := mockService.Called(baseURL, tenant, username, password)
	return args.Get(0).(dto.OutputDTO), args.Error(1)
}
