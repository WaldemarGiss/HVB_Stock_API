package api

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type stockService interface {
	CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, customError.ErrorStock)
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

		responseEntity.Share = share
		responseEntity.Numbers, _ = strconv.ParseFloat(strings.Replace(number, ",", ".", 1), 64)
		responseEntity.Price, _ = strconv.ParseFloat(strings.Replace(singlePrice, ",", ".", 1), 64)

		key := req.GetHeader("x-rapidapi-key")

		ret, err := stockAPI.stockService.CalculateEarning(key, responseEntity)
		fmt.Println(err.Code)
		if err.Code == 200 {
			req.JSON(http.StatusOK, ret)
		} else if err.Code == 204 {
			req.JSON(http.StatusNoContent, err.Text)
		} else if err.Code == 401 {
			req.JSON(http.StatusUnauthorized, err.Text)
		} else if err.Code == 403 {
			req.JSON(http.StatusForbidden, err.Text)
		} else if err.Code == 500 {
			req.JSON(http.StatusInternalServerError, err.Text)
		}
	}
}

func CreateRouter(router *gin.Engine, stockService stockService) {
	stockController := ProvideStockAPI(stockService)
	router.GET("HVB_Stock_API/calcGain", stockController.CalculateEarning)
}
