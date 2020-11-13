package api

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type stockService interface {
	Calculate(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock)
}

type StockAPI struct {
	stockService stockService
}

func ProvideStockAPI(service stockService) *StockAPI {
	return &StockAPI{stockService: service}
}

func (stockAPI *StockAPI) Calculate(req *gin.Context) {

	if req.GetHeader("x-rapidapi-key") == "" {
		req.JSON(http.StatusUnauthorized, "key missing")
	} else {
		share := req.Query("share")
		number := req.Query("number")
		singlePrice := req.Query("singlePrice")

		if share == "" || number == "" || singlePrice == "" {
			req.JSON(http.StatusBadRequest, "bad request")
		} else {
			var responseEntity dto.OutputDTO

			responseEntity.Share = share
			responseEntity.Numbers, _ = strconv.ParseFloat(strings.Replace(number, ",", ".", 1), 64)
			responseEntity.Price, _ = strconv.ParseFloat(strings.Replace(singlePrice, ",", ".", 1), 64)

			key := req.GetHeader("x-rapidapi-key")

			ret, err := stockAPI.stockService.Calculate(key, responseEntity)

			if err.Code == 0 {
				req.JSON(http.StatusOK, ret)
			} else if err.Code == 204 {
				req.JSON(http.StatusNoContent, "no content")
			} else if err.Code == 500 {
				req.JSON(http.StatusInternalServerError, "internal server error")
			}

		}
	}

}

func CreateRouter(router *gin.Engine, stockService stockService) {
	stockController := ProvideStockAPI(stockService)
	router.GET("HVB_Stock_API/calcGain", stockController.Calculate)
}
