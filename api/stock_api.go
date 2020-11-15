package api

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type stockService interface {
	CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error)
}

type StockAPI struct {
	stockService stockService
}

func ProvideStockAPI(service stockService) *StockAPI {
	return &StockAPI{stockService: service}
}

func (stockAPI *StockAPI) CalculateEarning(req *gin.Context) {

	share := req.Query("share")
	number := req.Query("number")
	singlePrice := req.Query("singlePrice")

	if share == "" || number == "" || singlePrice == "" {
		req.JSON(http.StatusBadRequest, "bad request")
	} else {
		var responseEntity dto.OutputDTO
		var stockError *customError.ErrorStock

		responseEntity.Share = share
		responseEntity.Numbers, _ = strconv.ParseFloat(strings.Replace(number, ",", ".", 1), 64)
		responseEntity.Price, _ = strconv.ParseFloat(strings.Replace(singlePrice, ",", ".", 1), 64)

		key := req.GetHeader("x-rapidapi-key")

		ret, err := stockAPI.stockService.CalculateEarning(key, responseEntity)

		if errors.As(err, &stockError) {
			if stockError.Code == 200 {
				req.JSON(http.StatusOK, ret)
			} else if stockError.Code == 204 {
				req.JSON(http.StatusNoContent, stockError.Text)
			} else if stockError.Code == 401 {
				req.JSON(http.StatusUnauthorized, stockError.Text)
			} else if stockError.Code == 403 {
				req.JSON(http.StatusForbidden, stockError.Text)
			} else if stockError.Code == 500 {
				req.JSON(http.StatusInternalServerError, stockError.Text)
			}
		}
	}
}

func CreateRouter(router *gin.Engine, stockService stockService) {
	stockController := ProvideStockAPI(stockService)
	router.GET("HVB_Stock_API/calcGain", stockController.CalculateEarning)
}
